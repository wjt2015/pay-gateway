package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	etcdfileutil "github.com/coreos/etcd/pkg/fileutil"
	"github.com/fsnotify/fsnotify"
	"github.com/pjoc-team/tracing/logger"
	"io/ioutil"
	"os"
	"sync"
)

// Store storage for Services
type Store interface {
	// Put put service
	Put(serviceName string, service *Service) error
	// Get get service name
	Get(serviceName string) (*Service, error)
}

// fileStore use file storage to implements the store interface
type fileStore struct {
	locker       sync.RWMutex
	filePath     string
	lockFilePath string
	file         *os.File
	lockedFile   *etcdfileutil.LockedFile
	services     map[string]*Service
}

// NewFileStore create file store
func NewFileStore(ctx context.Context, filePath string) (Store, error) {
	log := logger.Log()

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, etcdfileutil.PrivateFileMode)
	// file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		err2 := file.Close()
		if err2 != nil {
			log.Error(err2.Error())
		}
	}()
	lockFilePath := fmt.Sprintf("%s.lock", filePath)

	fs := &fileStore{
		filePath:     filePath,
		lockFilePath: lockFilePath,
		file:         file,
	}
	services, err := fs.readAll()
	if err != nil {
		log.Errorf("failed to read file, error: %v", err.Error())
		return nil, err
	}
	fs.services = services

	fs.watch(ctx)

	return fs, nil
}

func (f *fileStore) watch(ctx context.Context) {
	log := logger.Log()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Warnf("interrupt.")
				return
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Infof("event:", event)
				f.processFileEvent(ctx, event)

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Infof("error:", err)
			}
		}
	}()

	err = watcher.Add(f.filePath)
	if err != nil {
		log.Fatal(err)
	}
}

func (f *fileStore) processFileEvent(ctx context.Context, event fsnotify.Event) {
	log := logger.ContextLog(ctx)
	if event.Op&fsnotify.Write != fsnotify.Write {
		return
	}
	f.locker.Lock()
	defer f.locker.Unlock()

	log.Infof("modified file: %v", event.Name)
	all, err2 := f.readAll()
	if err2 != nil {
		log.Errorf("failed to read, error: %v", err2.Error())
		return
	}
	f.services = all
}

func (f *fileStore) lock() error {
	file, err := etcdfileutil.TryLockFile(
		f.lockFilePath, os.O_WRONLY|os.O_CREATE, 0777,
	)
	if err != nil {
		return err
	}
	f.lockedFile = file
	return nil
}
func (f *fileStore) unlock() error {
	if f.lockedFile == nil {
		return nil
	}
	err := os.Remove(f.lockFilePath)
	if err != nil {
		return err
	}
	return nil
}

func (f *fileStore) Put(serviceName string, service *Service) error {
	log := logger.Log()
	for {
		err := f.lock()
		if err != nil {
			log.Errorf(
				"failed to put service, "+
					"because lock failed serviceName: %v service: %#v", serviceName, service,
			)
			continue
		}
		break
	}
	defer func() {
		err := f.unlock()
		if err != nil {
			log.Errorf(
				"failed to put service, "+
					"because unlock failed serviceName: %v service: %#v", serviceName, service,
			)
		}
	}()
	file, err := os.OpenFile(f.filePath, os.O_RDWR|os.O_CREATE, etcdfileutil.PrivateFileMode)
	// file, err := os.Open(filePath)
	if err != nil {
		log.Errorf(
			"failed to put service, " +
				"because open file failed serviceName: %v service",
		)
		return err
	}
	defer func() {
		err2 := file.Close()
		if err2 != nil {
			log.Error(err2.Error())
		}
	}()

	services, err := f.readAll()
	if err != nil {
		return err
	}
	services[serviceName] = service
	marshal, err := json.Marshal(services)
	if err != nil {
		return err
	}
	_, err = file.Write(marshal)
	if err != nil {
		return err
	}
	return nil
}

func (f *fileStore) Get(serviceName string) (*Service, error) {
	f.locker.RLock()
	defer f.locker.RUnlock()
	return f.services[serviceName], nil
}

func (f *fileStore) readAll() (map[string]*Service, error) {
	log := logger.Log()
	raw, err := ioutil.ReadFile(f.filePath)
	if err != nil {
		log.Errorf("failed to read file: %v error: %v", f.filePath, err.Error())
		return nil, err
	}
	services := make(map[string]*Service)
	if len(raw) == 0 {
		return services, nil
	}
	err = json.Unmarshal(raw, &services)
	if err != nil {
		log.Errorf("failed to read file: %v error: %v", f.filePath, err.Error())
		return nil, err
	}
	return services, nil
}

package notify

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	s := &Scheduler{}
	s.QueueConfig = &QueueConfig{QueueType: "mysql", ConfigValue: &MySQLConfig{}}
	bytes, _ := yaml.Marshal(s)
	fmt.Println(string(bytes))
	fmt.Println(s.QueueConfig)
	queue, e := InstanceQueue(*s.QueueConfig, &Service{})
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	fmt.Println(queue)

}

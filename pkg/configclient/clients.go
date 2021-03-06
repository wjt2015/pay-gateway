package configclient

import (
	"context"
	"fmt"
	"github.com/pjoc-team/pay-gateway/pkg/config"
	"github.com/pjoc-team/tracing/logger"
	"github.com/spf13/pflag"
)

// ConfigClients 所有配置
type ConfigClients interface {
	// GetAppChannelConfig 获取渠道配置
	GetAppChannelConfig(ctx context.Context, appID string, method string) (
		[]*AppIDChannelConfig, error,
	)

	// GetAppConfig 获取应用配置
	GetAppConfig(ctx context.Context, appID string) (*MerchantConfig, error)

	// GetPayConfig 获取支付配置
	GetPayConfig(ctx context.Context) (*PayConfig, error)

	// GetNotifyConfig notify config
	GetNotifyConfig(ctx context.Context) (*NotifyConfig, error)
}

// configClients 所有配置
type configClients struct {
	PayConfigServer            *configClient
	NotifyConfigServer         *configClient
	ServiceConfigServer        *configClient
	ChannelServiceConfigServer *configClient
	MerchantConfigServer       *configClient
	PersonalMerchantServer     *configClient
	AppIDChannelConfigServer   *configClient
	FlagSet                    *pflag.FlagSet
}

// configClient 包装配置，简化获取配置函数
type configClient struct {
	url       string
	s         config.Server
	configURL ConfigURL
}

func (c *configClient) UnmarshalGetConfig(
	ctx context.Context, ptr interface{}, keys ...string,
) error {
	log := logger.ContextLog(ctx)
	if c == nil {
		err := fmt.Errorf("config is not initialized")
		return err
	} else if c.s == nil {
		err := fmt.Errorf("config is not initialized, please add flag: %v", c.configURL.Flag())
		log.Errorf(err.Error())
		return err
	}
	return c.s.UnmarshalGetConfig(ctx, ptr, keys...)
}

// NewConfigClients 创建配置
func NewConfigClients(opts ...Option) (ConfigClients, *pflag.FlagSet, error) {
	o, err := newOpts()
	if err != nil {
		return nil, nil, err
	}
	o.apply(opts...)

	c := &configClients{}

	err = c.initConfigs(o)
	if err != nil {
		return nil, nil, err
	}

	return c, c.FlagSet, nil
}

// newConfigClient 使用url创建配置客户端
func newConfigClient(url ConfigURL) (*configClient, error) {
	log := logger.Log()
	c := &configClient{
		configURL: url,
		url:       url.URL(),
	}
	if !c.configURL.Required() {
		return c, nil
	}
	if url.URL() == "" {
		err := fmt.Errorf("config is not initialized, please add flag: %v", c.configURL.Flag())
		log.Fatal(err.Error())
		return nil, err
	}
	server, err := config.InitConfigServer(url.URL())
	if err != nil {
		log.Fatalf(
			"failed to init config client: %v url: %v error: %v", url.Flag(), url.URL(), err.Error(),
		)
		return nil, err
	}
	c.s = server

	return c, nil
}

func (c *configClients) initConfigs(o *options) error {
	client, err := newConfigClient(o.PayConfigServerURL)
	if err != nil {
		return err
	}
	c.PayConfigServer = client

	client, err = newConfigClient(o.NotifyConfigServerURL)
	if err != nil {
		return err
	}
	c.NotifyConfigServer = client

	client, err = newConfigClient(o.ServiceConfigServerURL)
	if err != nil {
		return err
	}
	c.ServiceConfigServer = client

	client, err = newConfigClient(o.ChannelServiceConfigServerURL)
	if err != nil {
		return err
	}
	c.ChannelServiceConfigServer = client

	client, err = newConfigClient(o.MerchantConfigServerURL)
	if err != nil {
		return err
	}
	c.MerchantConfigServer = client

	client, err = newConfigClient(o.PersonalMerchantServerURL)
	if err != nil {
		return err
	}
	c.PersonalMerchantServer = client

	client, err = newConfigClient(o.AppIDChannelConfigServerURL)
	if err != nil {
		return err
	}
	c.AppIDChannelConfigServer = client

	c.FlagSet = o.ps
	return nil
}

func (c *configClients) GetAppChannelConfig(
	ctx context.Context, appID string, method string,
) ([]*AppIDChannelConfig, error) {
	log := logger.ContextLog(ctx)
	appConfig := make([]*AppIDChannelConfig, 0)
	err := c.AppIDChannelConfigServer.UnmarshalGetConfig(ctx, &appConfig, appID, method)
	if err != nil {
		log.Errorf(
			"failed to get channel config of appID: %v method: %v error: %v", appID, method,
			err.Error(),
		)
		return nil, err
	}
	return appConfig, nil
}

func (c *configClients) GetAppConfig(ctx context.Context, appID string) (*MerchantConfig, error) {
	log := logger.ContextLog(ctx)
	merchantConfig := &MerchantConfig{}
	err := c.MerchantConfigServer.UnmarshalGetConfig(ctx, merchantConfig, appID)
	if err != nil {
		log.Errorf(
			"failed to get merchant config of appID: %v method: %v error: %v", appID, err.Error(),
		)
		return nil, err
	}
	return merchantConfig, nil
}

func (c *configClients) GetPayConfig(ctx context.Context) (*PayConfig, error) {
	log := logger.ContextLog(ctx)
	payConfig := &PayConfig{}
	err := c.PayConfigServer.UnmarshalGetConfig(ctx, payConfig)
	if err != nil {
		log.Errorf(
			"failed to get pay config of error: %v", payConfig,
			err.Error(),
		)
		return nil, err
	}
	return payConfig, nil
}

func (c *configClients) GetNotifyConfig(ctx context.Context) (*NotifyConfig, error) {
	log := logger.ContextLog(ctx)
	notifyConfig := &NotifyConfig{}
	err := c.NotifyConfigServer.UnmarshalGetConfig(ctx, notifyConfig)
	if err != nil {
		log.Errorf(
			"failed to get notify config of error: %v", notifyConfig,
			err.Error(),
		)
		return nil, err
	}
	return notifyConfig, nil
}

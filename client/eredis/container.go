package eredis

import (
	"github.com/go-redis/redis"
	"github.com/gotomicro/ego/core/conf"
	"github.com/gotomicro/ego/core/elog"
)

type Option func(c *Container)

type Container struct {
	config *Config
	name   string
	logger *elog.Component
}

func DefaultContainer() *Container {
	return &Container{
		logger: elog.EgoLogger.With(elog.FieldMod("client.egorm")),
	}
}

func Load(key string) *Container {
	c := DefaultContainer()
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, &config); err != nil {
		c.logger.Panic("parse config error", elog.FieldErr(err), elog.FieldKey(key))
		return c
	}

	c.config = config
	c.name = key

	return c
}

func WithStub() Option {
	return func(c *Container) {
		if c.config.Addr == "" && len(c.config.Addrs) == 0 {
			c.logger.Panic("no address in redis config", elog.FieldName(c.name))
		}
		if c.config.Addr != "" {
			c.config.Addrs = []string{c.config.Addr}
		}
		c.config.Mode = StubMode
	}
}

func WithCluster() Option {
	return func(c *Container) {
		c.config.Mode = ClusterMode
	}
}

// Build ...
func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}
	count := len(c.config.Addrs)
	if count < 1 {
		c.logger.Panic("no address in redis config", elog.FieldName(c.name))
	}
	if len(c.config.Mode) == 0 {
		c.config.Mode = StubMode
		if count > 1 {
			c.config.Mode = ClusterMode
		}
	}
	var client redis.Cmdable
	switch c.config.Mode {
	case ClusterMode:
		if count == 1 {
			c.logger.Warn("redis config has only 1 address but with cluster mode")
		}
		client = c.buildCluster()
	case StubMode:
		if count > 1 {
			c.logger.Warn("redis config has more than 1 address but with stub mode")
		}
		client = c.buildStub()
	default:
		c.logger.Panic("redis mode must be one of (stub, cluster)")
	}
	return &Component{
		Config: c.config,
		Client: client,
	}
}

func (c *Container) buildCluster() *redis.ClusterClient {
	clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        c.config.Addrs,
		MaxRedirects: c.config.MaxRetries,
		ReadOnly:     c.config.ReadOnly,
		Password:     c.config.Password,
		MaxRetries:   c.config.MaxRetries,
		DialTimeout:  c.config.DialTimeout,
		ReadTimeout:  c.config.ReadTimeout,
		WriteTimeout: c.config.WriteTimeout,
		PoolSize:     c.config.PoolSize,
		MinIdleConns: c.config.MinIdleConns,
		IdleTimeout:  c.config.IdleTimeout,
	})
	if err := clusterClient.Ping().Err(); err != nil {
		switch c.config.OnDialError {
		case "panic":
			c.logger.Panic("start cluster redis", elog.FieldErr(err))
		default:
			c.logger.Error("start cluster redis", elog.FieldErr(err))
		}
	}
	return clusterClient
}

func (c *Container) buildStub() *redis.Client {
	stubClient := redis.NewClient(&redis.Options{
		Addr:         c.config.Addrs[0],
		Password:     c.config.Password,
		DB:           c.config.DB,
		MaxRetries:   c.config.MaxRetries,
		DialTimeout:  c.config.DialTimeout,
		ReadTimeout:  c.config.ReadTimeout,
		WriteTimeout: c.config.WriteTimeout,
		PoolSize:     c.config.PoolSize,
		MinIdleConns: c.config.MinIdleConns,
		IdleTimeout:  c.config.IdleTimeout,
	})

	if err := stubClient.Ping().Err(); err != nil {
		switch c.config.OnDialError {
		case "panic":
			c.logger.Panic("start stub redis", elog.FieldErr(err))
		default:
			c.logger.Error("start stub redis", elog.FieldErr(err))
		}
	}
	return stubClient
}

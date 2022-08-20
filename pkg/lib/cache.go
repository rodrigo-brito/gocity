package lib

import (
	"time"

	"github.com/karlseguin/ccache"
)

type Cache interface {
	Get(key string) (bool, []byte)
	Set(key string, value []byte, TTL time.Duration)
	GetSet(key string, set func() ([]byte, error), TTL time.Duration) ([]byte, error)
}

func NewCache() Cache {
	return &cache{
		client: ccache.New(ccache.Configure()),
	}
}

type cache struct {
	client *ccache.Cache
}

func (c *cache) Get(key string) (bool, []byte) {
	item := c.client.Get(key)
	if item != nil {
		return true, item.Value().([]byte)
	}

	return false, nil
}

func (c *cache) Set(key string, value []byte, TTL time.Duration) {
	c.client.Set(key, value, TTL)
}

func (c *cache) GetSet(key string, getValue func() ([]byte, error), TTL time.Duration) ([]byte, error) {
	hit, value := c.Get(key)
	if hit {
		return value, nil
	}

	value, err := getValue()
	if err != nil {
		return nil, err
	}

	c.Set(key, value, TTL)

	return value, nil
}

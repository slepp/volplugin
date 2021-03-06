package config

import (
	"errors"
	"path"

	"github.com/alecthomas/units"
	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

const (
	rootVolume = "volumes"
	rootUse    = "users"
	rootTenant = "tenants"
)

var (
	quux = units.Mebibyte
)

var defaultPaths = []string{rootVolume, rootUse, rootTenant}

// ErrExist indicates when a key in etcd exits already. Used for create logic.
var ErrExist = errors.New("Already exists")

// Request provides a request structure for communicating with the
// volmaster.
type Request struct {
	Volume string `json:"volume"`
	Tenant string `json:"tenant"`
}

// RequestCreate provides a request structure for creating new volumes.
type RequestCreate struct {
	Tenant string            `json:"tenant"`
	Volume string            `json:"volume"`
	Opts   map[string]string `json:"opts"`
}

// TopLevelConfig is the top-level struct for communicating with the intent store.
type TopLevelConfig struct {
	etcdClient client.KeysAPI
	prefix     string
}

// NewTopLevelConfig creates a TopLevelConfig struct which can drive communication
// with the configuration store.
func NewTopLevelConfig(prefix string, etcdHosts []string) (*TopLevelConfig, error) {
	etcdCfg := client.Config{
		Endpoints: etcdHosts,
	}

	etcdClient, err := client.New(etcdCfg)
	if err != nil {
		return nil, err
	}

	config := &TopLevelConfig{
		prefix:     prefix,
		etcdClient: client.NewKeysAPI(etcdClient),
	}

	config.etcdClient.Set(context.Background(), config.prefix, "", &client.SetOptions{Dir: true})
	for _, path := range defaultPaths {
		config.etcdClient.Set(context.Background(), config.prefixed(path), "", &client.SetOptions{Dir: true})
	}

	return config, nil
}

func (c *TopLevelConfig) prefixed(strs ...string) string {
	str := c.prefix
	for _, s := range strs {
		str = path.Join(str, s)
	}

	return str
}

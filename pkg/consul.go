package pkg

import (
	"errors"
	"fmt"

	"os"
	"path"
	"path/filepath"
	"unsafe"

	"github.com/hashicorp/consul/api"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/viper"

	"consul-webhook/config"
)

var ErrNoConsulKey = errors.New("consul: no consul key found")

const consulPath = "/hxyljk/devops/repo"

var consulConf *ConsulConf

type ConsulConf struct {
	Repo       *config.Repo `toml:"repo"`
	Address    string       `toml:"address"`
	Scheme     string       `toml:"scheme"`
	Datacenter string       `toml:"datacenter"`
	Token      string       `toml:"token"`
}

type Factory struct {
	client *Client
}

type Consul interface {
	Get(key string) ([]byte, error)
	Put(key, data string) error
	Keys(prefix string, separators ...string) ([]string, error)
}

type Client struct {
	config *ConsulConf
	client *api.Client
}

func (c *Client) Get(key string) (value []byte, err error) {
	kv := c.client.KV()
	kvPair, _, err := kv.Get(key, nil)
	if err != nil {
		return
	}
	if kvPair == nil {
		err = ErrNoConsulKey
		return
	}

	return kvPair.Value, nil
}

func (c *Client) Put(key string) error {
	kv := c.client.KV()
	repoPath := GetTmpPath()
	file := path.Join(repoPath, key)
	f, _ := os.ReadFile(file)
	pair := api.KVPair{
		Key:   key,
		Value: f,
	}
	_, err := kv.Put(&pair, nil)
	return err
}

func (c *Client) Delete(key string) error {
	kv := c.client.KV()
	_, err := kv.Delete(key, nil)
	return err
}

// Keys  获取consul目录所有节点
// 这个底层方法不会去扣除目录节点，注意
func (c *Client) Keys(prefix string, separators ...string) ([]string, error) {
	var result []string
	var err error

	var separator string = ""
	if len(separators) > 0 {
		separator = separators[0]
	}

	result, _, err = c.client.KV().Keys(prefix, separator, nil)
	return result, err
}

// makeString 高效转为string
func makeString(data []byte) string {
	return *(*string)(unsafe.Pointer(&data))
}

func NewFactory() (*Factory, error) {
	apiConfig := api.DefaultConfig()
	localCfg := config.GetConfig().Consul
	apiConfig.Address = localCfg.Address
	apiConfig.Scheme = localCfg.Scheme
	apiConfig.Token = localCfg.Token
	consulConf := &ConsulConf{
		Address: localCfg.Address,
		Token:   localCfg.Token,
		Scheme:  localCfg.Scheme,
	}
	apiClient, err := api.NewClient(apiConfig)
	if err != nil {
		return nil, err
	}
	return &Factory{
		client: &Client{
			config: consulConf,
			client: apiClient,
		},
	}, nil
}

func (f *Factory) NewClient() *Client {
	return f.client
}

func GetConsulConfig() (*ConsulConf, error) {
	consulPath := consulPath
	if consulConf != nil {
		return consulConf, nil
	}
	consulFac, err := NewFactory()
	if err != nil {
		return nil, err
	}

	consulClient := consulFac.NewClient()
	value, err := consulClient.Get(filepath.ToSlash(filepath.Join(viper.GetString("app.env"), consulPath)))

	if err != nil {
		return nil, fmt.Errorf("consulClient.Get consul key[%s],error:%v", consulPath, err)
	}
	consulConf = new(ConsulConf)
	err = toml.Unmarshal(value, consulConf)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal mysql consul config error:%v,key[%s]", err, consulPath)
	}

	return consulConf, nil
}

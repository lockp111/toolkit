package consul

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/reader"
	"github.com/micro/go-micro/config/source"
	"github.com/micro/go-plugins/config/source/consul"
	log "github.com/sirupsen/logrus"
)

type conf struct {
	source  source.Source
	prefix  string
	address string
}

var consulConf *conf

func init() {
	addr := os.Getenv("CONSUL_ADDR")
	if len(addr) == 0 {
		addr = "127.0.0.1:8500"
	}

	consulConf = &conf{
		address: addr,
	}
}

func getPrefixedPath(path ...string) []string {
	prefixPaths := strings.Split(consulConf.prefix, "/")
	path = append(prefixPaths, path...)
	return path
}

// GetAddress ..
func GetAddress() string {
	return consulConf.address
}

// InitSource Directly init source. Use it without micro service
func InitSource(addr string, prefix ...string) {
	consulConf.address = addr

	var opts = []source.Option{consul.WithAddress(consulConf.address)}
	if len(prefix) != 0 {
		consulConf.prefix = strings.Trim(prefix[0], "/")
		opts = append(opts,
			consul.WithPrefix(consulConf.prefix),
			consul.StripPrefix(false),
		)
	}

	consulConf.source = consul.NewSource(opts...)
}

// ConfigGet ...
func ConfigGet(x interface{}, path ...string) error {
	conf := config.NewConfig()
	if err := conf.Load(consulConf.source); err != nil {
		return err
	}
	defer conf.Close()

	path = getPrefixedPath(path...)
	if err := conf.Get(path...).Scan(x); err != nil {
		return err
	}

	return nil
}

// ConfigWatch ...
func ConfigWatch(scanFunc func(reader.Value), path ...string) error {
	conf := config.NewConfig()
	if err := conf.Load(consulConf.source); err != nil {
		return err
	}

	path = getPrefixedPath(path...)
	w, err := conf.Watch(path...)
	if err != nil {
		return err
	}

	go func() {
		val := conf.Get(path...)
		scanFunc(val)

		for {
			v, err := w.Next()
			if err != nil {
				log.WithFields(log.Fields{
					"path": path,
				}).WithError(err).
					Error("Config watch next value")
				continue
			}

			scanFunc(v)
		}
	}()

	return nil
}

// Sync ...
func Sync(service string, conf interface{}) error {
	data, _ := json.MarshalIndent(conf, "", "\t")

	apiConf := api.DefaultConfig()
	apiConf.Address = consulConf.address

	// Get a new client
	client, err := api.NewClient(apiConf)
	if err != nil {
		return err
	}

	// Get a handle to the KV API
	kv := client.KV()

	// PUT a new KV pair
	p := &api.KVPair{Key: service, Value: data}
	_, err = kv.Put(p, nil)
	return err
}

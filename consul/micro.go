package consul

import (
	"os"
	"strings"

	"github.com/micro/go-config"
	"github.com/micro/go-config/reader"
	"github.com/micro/go-config/source"
	"github.com/micro/go-config/source/consul"
	log "github.com/sirupsen/logrus"
)

type conf struct {
	source  source.Source
	prefix  string
	address string
}

var consulConf *conf

func init() {
	addr := os.Getenv("CONSULADDR")
	if len(addr) == 0 {
		addr = "127.0.0.1:8500"
	}

	consulConf = &conf{
		address: addr,
	}
}

func getPrefixedPath(path ...string) []string {
	prefixPaths := strings.Split(consulConf.prefix, "/")
	path = append(prefixPaths[1:], path...)

	return path
}

// GetConsulAddress ..
func GetConsulAddress() string {
	return consulConf.address
}

// InitSource Directly init source. Use it without micro service
func InitSource(addr string, prefix ...string) {
	consulConf.address = addr
	var opts = []source.Option{consul.WithAddress(consulConf.address)}
	if len(prefix) != 0 {
		consulConf.prefix = prefix[0]
		opts = append(opts,
			consul.WithPrefix(consulConf.prefix),
			consul.StripPrefix(true),
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

	path = getPrefixedPath(path...)
	defer conf.Close()

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

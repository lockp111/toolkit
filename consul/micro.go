package consul

import (
	"strings"

	"github.com/micro/go-config"
	"github.com/micro/go-config/reader"
	"github.com/micro/go-config/source"
	"github.com/micro/go-config/source/consul"
	log "github.com/sirupsen/logrus"
)

var (
	consulSrouce    source.Source
	configPrefix    string
	registryAddress string
)

func getPrefixedPath(path ...string) []string {
	prefixPaths := strings.Split(configPrefix, "/")
	path = append(prefixPaths[1:], path...)

	return path
}

// GetRegistryAddress ..
func GetRegistryAddress() string {
	return registryAddress
}

// InitSource Directly init source. Use it without micro service
func InitSource(addr, prefix string) {
	configPrefix = prefix
	registryAddress = addr
	consulSrouce = consul.NewSource(
		consul.WithAddress(registryAddress),
		consul.WithPrefix(configPrefix),
		consul.StripPrefix(true),
	)
}

// ConfigGet ...
func ConfigGet(x interface{}, path ...string) error {
	conf := config.NewConfig()
	if err := conf.Load(consulSrouce); err != nil {
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
	if err := conf.Load(consulSrouce); err != nil {
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

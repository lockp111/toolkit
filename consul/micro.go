package consul

import (
	"strings"

	"github.com/micro/cli"
	"github.com/micro/go-config"
	"github.com/micro/go-config/reader"
	"github.com/micro/go-config/source"
	"github.com/micro/go-config/source/consul"
	"github.com/micro/go-micro"
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

// InitService ...
func InitService(service micro.Service, prefix string, flags ...micro.Option) {
	var options = []micro.Option{
		micro.Flags(cli.StringFlag{
			Name:        "config_prefix",
			Value:       prefix,
			Usage:       "consul config prefix",
			EnvVar:      "CONFIG_PREFIX",
			Destination: &configPrefix,
		}),

		micro.Action(func(c *cli.Context) {
			registryAddress = c.GlobalString("registry_address")
		}),
	}

	options = append(options, flags...)
	service.Init(options...)

	var opts = []source.Option{consul.WithAddress(registryAddress)}
	if configPrefix != "" {
		opts = append(opts,
			consul.WithPrefix(configPrefix),
			consul.StripPrefix(true),
		)
	}
	consulSrouce = consul.NewSource(opts...)
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

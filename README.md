# Utils

## env
保存一些程序常用到的环境变量

## consul
目前的go-micro把registry和config进行了分离，在多台服务器上运行服务并指定consul时，在默认条件下会遇到麻烦。
因此用此包用来初始化服务，使config和registry都使用同一个选项 --registry_address

```go

// Run ...
func Run() error {
	service := micro.NewService(
		micro.Name(ServiceName),
		micro.RegisterTTL(ttl),
		micro.RegisterInterval(interval),
		micro.Version(ver),
	)

	consul.InitService(service, "/config/your_service_prefix")
    consul.ConfigGet(...) //
}
```

```bash
bin-to-exec --registry_address "http://consul_host"
```

## log
使用 "github.com/sirupsen/logrus"

### 日志级别
日志级别放在环境变量里面
LOGLEVEL=debug bin-to-exec

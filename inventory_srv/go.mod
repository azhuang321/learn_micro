module inventory_srv

go 1.16

require (
	github.com/gookit/ini/v2 v2.0.11
	github.com/natefinch/lumberjack v2.0.0+incompatible
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.19.0
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gorm.io/driver/mysql v1.1.2
	gorm.io/gorm v1.21.14
)

require (
	github.com/go-redis/redis/v8 v8.11.3
	github.com/hashicorp/consul/api v1.10.1
	github.com/nacos-group/nacos-sdk-go v1.0.8
	github.com/satori/go.uuid v1.2.0
)

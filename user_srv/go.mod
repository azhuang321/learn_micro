module user_srv

go 1.16

require (
	github.com/apache/rocketmq-client-go/v2 v2.1.0
	github.com/fsnotify/fsnotify v1.5.1
	github.com/hashicorp/consul/api v1.10.1
	github.com/mbobakov/grpc-consul-resolver v1.4.4
	github.com/nacos-group/nacos-sdk-go v1.0.9
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/viper v1.9.0
	go.uber.org/zap v1.19.1
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gorm.io/driver/mysql v1.1.2
	gorm.io/gorm v1.21.15
)

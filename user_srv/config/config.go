package config

// Config 全局配置
type Config struct {
	RunMod string `mapstructure:"run_mod" json:"run_mod"` //运行模式
	Port   int32  `mapstructure:"port" json:"port"`       //端启动监听口
	Logger Logger `mapstructure:"logger" json:"logger"`   //日志配置
}

// Logger 日志配置
type Logger struct {
	FileName   string `mapstructure:"file_name" json:"file_name"`
	MaxSize    int    `mapstructure:"max_size" json:"max_size"`
	MaxBackups int    `mapstructure:"max_backups" json:"max_backups"`
	MaxAge     int    `mapstructure:"max_age" json:"max_age"`
	Compress   bool   `mapstructure:"compress" json:"compress"`
}

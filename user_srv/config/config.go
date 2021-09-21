package config

// Config 全局配置
type Config struct {
	RunMod  string `mapstructure:"run_mod" json:"run_mod"`   //运行模式
	Port    int32  `mapstructure:"port" json:"port"`         //端启动监听口
	MD5Salt string `mapstructure:"md5_salt" json:"md5_salt"` //加密盐
	Logger  Logger `mapstructure:"logger" json:"logger"`     //日志配置
	Mysql   Mysql  `mapstructure:"mysql" json:"mysql"`       //数据库配置
}

// Logger 日志配置
type Logger struct {
	FileName   string `mapstructure:"file_name" json:"file_name"`
	MaxSize    int    `mapstructure:"max_size" json:"max_size"`
	MaxBackups int    `mapstructure:"max_backups" json:"max_backups"`
	MaxAge     int    `mapstructure:"max_age" json:"max_age"`
	Compress   bool   `mapstructure:"compress" json:"compress"`
}

// Mysql mysql配置
type Mysql struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Database string `mapstructure:"database" json:"database"`
	Username string `mapstructure:"username" json:"username"`
	Password string `mapstructure:"password" json:"password"`
}

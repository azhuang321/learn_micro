package config

type ServerConfig struct {
	Name   string   `json:"name"`
	Tags   []string `json:"tags"`
	Mysql  Mysql    `json:"mysql"`
	Consul Consul   `json:"consul"`
}
type Mysql struct {
	Db       string `json:"db"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}
type Consul struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

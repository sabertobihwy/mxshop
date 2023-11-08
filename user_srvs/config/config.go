package config

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Username string `mapstructure:"username" json:"username"`
	Password string `mapstructure:"password" json:"password"`
	Db       string `mapstructure:"db" json:"db"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type ServiceConfig struct {
	Name         string       `mapstructure:"name" json:"name"`
	MysqlConfig  MysqlConfig  `mapstructure:"mysql" json:"mysql"`
	ConsulConfig ConsulConfig `mapstructure:"consul" json:"consul"`
	Host         string       `mapstructure:"host" json:"host"`
}

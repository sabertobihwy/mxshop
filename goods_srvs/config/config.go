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
	Tags         []string     `mapstructure:"tags" json:"tags"`
	MysqlConfig  MysqlConfig  `mapstructure:"mysql" json:"mysql"`
	ConsulConfig ConsulConfig `mapstructure:"consul" json:"consul"`
	Host         string       `mapstructure:"host" json:"host"`
}
type NacosConfig struct {
	Host      string `mapstructure:"host" json:"host"`
	Port      int    `mapstructure:"port" json:"port"`
	Namespace string `mapstructure:"namespace" json:"namespace"`
	Dataid    string `mapstructure:"dataid" json:"dataid"`
	Group     string `mapstructure:"group" json:"group"`
	User      string `mapstructure:"user" json:"user"`
	Password  string `mapstructure:"password" json:"password"`
}

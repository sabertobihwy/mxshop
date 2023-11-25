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

type RedisConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}
type GoodService struct {
	Name string `mapstructure:"name" json:"name"`
}
type InventoryService struct {
	Name string `mapstructure:"name" json:"name"`
}
type ServiceConfig struct {
	Name         string           `mapstructure:"name" json:"name"`
	MysqlConfig  MysqlConfig      `mapstructure:"mysql" json:"mysql"`
	ConsulConfig ConsulConfig     `mapstructure:"consul" json:"consul"`
	RedisConfig  RedisConfig      `mapstructure:"redis" json:"redis"`
	Host         string           `mapstructure:"host" json:"host"`
	Tags         []string         `mapstructure:"tags" json:"tags"`
	GoodSrv      GoodService      `mapstructure:"goods_srvs" json:"goods_srvs"`
	InventorySrv InventoryService `mapstructure:"stocks_srvs" json:"stocks_srvs"`
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

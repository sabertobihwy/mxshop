package initialize

import (
	"encoding/json"
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"mxshop_srvs/stocks_srvs/global"
)

func GetSystemEnv(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env) // can get sys env variable
}

func InitConfig() {
	flg := GetSystemEnv("MXSHOP_CONFIG_FLAG")
	filePrefix := "config"
	configFilePath := fmt.Sprintf("stocks_srvs/%s_pro.yaml", filePrefix)
	if flg {
		configFilePath = fmt.Sprintf("stocks_srvs/%s_local.yaml", filePrefix)
	}
	zap.S().Info(configFilePath)

	v := viper.New() // return a Viper
	v.SetConfigFile(configFilePath)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	nc := global.NacosConfig
	if err := v.Unmarshal(nc); err != nil { // unserilaize to Object
		panic(err)
	}

	//create ServerConfig
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(global.NacosConfig.Host, uint64(global.NacosConfig.Port),
			constant.WithContextPath("/nacos")),
	}

	//create ClientConfig
	cc := *constant.NewClientConfig(
		constant.WithNamespaceId(global.NacosConfig.Namespace),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("tmp/nacos/log"),
		constant.WithCacheDir("tmp/nacos/cache"),
		constant.WithLogLevel("debug"),
	)

	// create config client
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err.Error())
	}
	//fmt.Printf("GetConfig,info : %s,%s,%s \n",
	//	global.NacosConfig.Dataid, global.NacosConfig.Group, global.NacosConfig.Namespace)
	content, err := client.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.Dataid,
		Group:  global.NacosConfig.Group,
	})
	//fmt.Printf("GetConfig,config : %s", content)
	err = json.Unmarshal([]byte(content), &global.ServiceConfig)
	if err != nil {
		panic(err.Error())
	}

}

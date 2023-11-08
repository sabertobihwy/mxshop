package initialize

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mxshop_srvs/user_srvs/global"
)

func GetSystemEnv(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env) // can get sys env variable
}

func InitConfig() {
	flg := GetSystemEnv("MXSHOP_CONFIG_FLAG")
	filePrefix := "config"
	configFilePath := fmt.Sprintf("user_srvs/%s_pro.yaml", filePrefix)
	if flg {
		configFilePath = fmt.Sprintf("user_srvs/%s_local.yaml", filePrefix)
	}
	zap.S().Info(configFilePath)

	v := viper.New() // return a Viper
	v.SetConfigFile(configFilePath)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	sc := global.ServiceConfig
	if err := v.Unmarshal(sc); err != nil { // unserilaize to Object
		panic(err)
	}
}

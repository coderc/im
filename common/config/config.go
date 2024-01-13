package config

import (
	"time"

	"os"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(GetEnv("config", "./im.yaml", "配置文件"))
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

// GetEndPointsForDiscovery 获取服务发现的地址
func GetEndPointsForDiscovery() []string {
	return viper.GetStringSlice("discovery.endpoints")
}

// GetTimeoutForDiscovery 获取链接服务发现集群的超时时间(ms)
func GetTimeoutForDiscovery() time.Duration {
	return viper.GetDuration("discovery.timeout") * time.Second
}

func GetServicePathForIPconfig() string {
	return viper.GetString("ip_config.service_path")
}

func IsDebug() bool {
	env := viper.GetString("global.env")
	return env == "debug"
}

// GetEnv 获取环境变量中的值
func GetEnv(name, value, usage string) string {
	envVar := os.Getenv(name)
	if envVar == "" {
		logger.Warnf("please check env var, (name: \"%s\", default: \"%s\", usage: \"%s\")", name, value, usage)
		return value
	}
	return envVar
}

package config

import (
	"strconv"
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
	str := os.Getenv(name)
	if str == "" {
		logger.Warnf("please check env var, (name: \"%s\", default: \"%s\", usage: \"%s\")", name, value, usage)
		return value
	}
	return str
}

func GetEnvInt32(name string, value int32, usage string) int32 {
	vInt64 := GetEnvInt64(name, int64(value), usage)
	return int32(vInt64)
}

func GetEnvInt64(name string, value int64, usage string) int64 {
	str := GetEnv(name, "", usage)
	if str == "" {
		return value
	}

	v, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		logger.Warnf("ParseInt failed, name: %s, usage: %s, err: %v", name, usage, err)
		return value
	}
	return v
}

package util

import (
	"time"

	"github.com/spf13/viper"
)

// 用于保存所有从环境变量/配置文件中读到的配置
type Config struct {
	//`mapstructure` tag(标签)
	//用于将viper从env\json\ini等文件或者环境变量读取来的
	//这个标签的作用是告诉 mapstructure 在解码过程中将 Viper 配置中的“db_driver”键的值赋给 DBDriver 字段。
	DBdriver             string        `mapstructure:"db_driver"`
	DBsource             string        `mapstructure:"db_source"`
	ServerAddress        string        `mapstructure:"Server_address"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	RedisAddress         string        `mapstructure:"REDIS_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)  //参数为配置文件的路径
	viper.SetConfigName("app") //参数为配置文件的名称
	viper.SetConfigType("env") //参数为配置文件的类型
	viper.AutomaticEnv()       //检查环境变量
	err = viper.ReadInConfig() //读取配置
	if err != nil {
		return
	}
	viper.Unmarshal(&config)
	return
	//采用传参式返回，即给返回值声明了变量，
	//这种情况下，不需要在return中写出返回值，因为其会根据变量内容自动返回
}

//此程序基本只需要读取app.env即可，因为config里的变量信息字段都在app.env里写好

package utils

import (
	"bytes"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
)


func GetFromConsulKV(endpoint, key string) (value []byte, err error) {
	// init config connection to consul
	config := api.DefaultConfig()
	if endpoint != "" {
		config.Address = endpoint
	}

	// init consul client
	client, err := api.NewClient(config)
	if err != nil {
		return
	}
	kv := client.KV()

	// get key
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return
	}
	if pair == nil {
		err = fmt.Errorf("remote conf key is not existed: %v", key)
		return
	}

	value = pair.Value
	return
}


// LoadConfig -- read conf from byte
func LoadConfig(configType string, value []byte, isMerge bool) (err error) {
	viper.SetConfigType(configType)
	if !isMerge {
		err = viper.ReadConfig(bytes.NewBuffer(value))
	} else {
		err = viper.MergeConfig(bytes.NewBuffer(value))
	}
	return
}

func ReadConfig(fileName string, configPaths ...string) bool {
	viper.SetConfigName(fileName)
	if len(configPaths) < 1 {
		// look for current dir
		viper.AddConfigPath(".")
	} else {
		for _, configPath := range configPaths {
			viper.AddConfigPath(configPath)
		}
	}
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Cannot read config file. %s", err)
		return false
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

	return true
}

// ReadConfigByFile read config file by file path
func ReadConfigByFile(file string) bool {
	viper.SetConfigFile(file)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Cannot read config file. %s", err)
		return false
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

	return true
}
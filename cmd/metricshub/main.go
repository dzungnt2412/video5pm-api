package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
	"video5pm-api/core/utils"
	"video5pm-api/pkg/database"
	"video5pm-api/pkg/logger"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"

	"video5pm-api/cmd/metricshub/api"
)

var (
	httpHost            string
	httpPort            int64
	configType          string
	configFile          string
	configRemoteAddress string
	configRemoteKeys    string
	mysqlConn           *gorm.DB
	err                 error
)

func init() {
	InitDefaultFlags()
	flag.Parse()

	if configType == "file" {
		if utils.IsStringEmpty(configFile) {
			// read by default
			utils.ReadConfig("conf", "./conf", ".")
		} else {
			// read by input file
			utils.ReadConfigByFile(configFile)
		}

	} else {
		configRemoteKeys := utils.StringSlice(utils.StringTrimSpace(configRemoteKeys), ",")

		if len(configRemoteKeys) < 1 {
			logger.Log.Warn("This app has no conf defined")
		} else {
			for index := 0; index < len(configRemoteKeys); index++ {
				isMerge := true
				if index == 0 {
					isMerge = false
				}
				remoteKey := utils.StringTrimSpace(configRemoteKeys[index])
				if utils.IsStringEmpty(remoteKey) {
					continue
				}
				if remoteKey[0:1] != "/" {
					logger.Log.Errorw(fmt.Sprintf("Invalid key: %v. Remote key must start with /", remoteKey))
					continue
				}
				valueBytes, err := utils.GetFromConsulKV(configRemoteAddress, remoteKey)
				if err != nil {
					logger.Log.Warnf("Could not get key %v from consul. Details: %v", remoteKey, err)
					continue
				}
				err = utils.LoadConfig("toml", valueBytes, isMerge)
				if err != nil {
					logger.Log.Errorw(fmt.Sprintf("Could not load conf from remote key %v. Details: %v", remoteKey, err))
					continue
				}
				logger.Log.Infof("Loaded conf remote key %v", remoteKey)
			}
		}
	}
}

func InitDefaultFlags() {
	// flags for config
	flag.StringVar(&configType, "config-type", "file", "Configuration type: file or remote")
	flag.StringVar(&configFile, "config-file", "", "Configuration file")
	flag.StringVar(&configRemoteAddress, "config-remote-address", "", "Configuration remote address. ip:port")
	flag.StringVar(&configRemoteKeys, "config-remote-keys", "", "Configuration remote keys. Separate by ,")

	// flags for http
	flag.StringVar(&httpHost, "http-host", "", "HTTP listen host")
	flag.Int64Var(&httpPort, "http-port", 8888, "HTTP listen port")
}

func main() {
	// init db
	mysqlConn, err = database.InitDB()
	if err != nil {
		logger.Log.Infof("cant connect mysql", err)
		return
	}

	readTimeout := time.Duration(viper.GetInt("server.read_timeout")) * time.Second
	writeTimeout := time.Duration(viper.GetInt("server.write_timeout")) * time.Second
	endPoint := fmt.Sprintf("%v:%v", httpHost, httpPort)
	maxHeaderBytes := 1 << 20

	logger.Log.Infof("[info] start http server listening %s", endPoint)
	routersInit := api.InitRouter(mysqlConn)
	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		logger.Log.Info("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		_ = mysqlConn.Close()
		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Log.Panicf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	logger.Log.Infof("Server running: %s", endPoint)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Log.Fatalf("Could not listen on %s: %v\n", endPoint, err)
	}

	<-done
	logger.Log.Info("Server stopped")
}

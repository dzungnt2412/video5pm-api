package database

import (
	"fmt"
	"github.com/spf13/viper"
	"lionnix-metrics-api/pkg/logger"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//connectionString - build connection string to MySQL from config
func connectionString() string {
	param := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		viper.GetString("mysql.username"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.database"),
	)

	// Additional params
	additionParams := make(map[string]string)

	if len(viper.GetString("mysql.charset")) > 0 {
		additionParams["charset"] = viper.GetString("mysql.charset")
	}
	switch viper.GetBool("mysql.parse_time") {
	case true:
		additionParams["parseTime"] = "True"
	case false:
		additionParams["parseTime"] = "False"
	}
	if len(viper.GetString("mysql.loc")) > 0 {
		additionParams["loc"] = viper.GetString("mysql.loc")
	}
	additionParamsStr := make([]string, 0)
	for k, v := range additionParams {
		additionParamsStr = append(additionParamsStr, fmt.Sprintf("%s=%s", k, v))
	}

	return fmt.Sprintf("%s?%s", param, strings.Join(additionParamsStr, "&"))
}

func InitDB() (*gorm.DB, error) {
	connectionString := connectionString()
	logger.Log.Info("Connecting to MysSQL: ", connectionString)

	db, err := gorm.Open("mysql", connectionString)
	if err != nil {
		logger.Log.Errorf("Cannot connect to MySQL, %v", err)
	}

	db.LogMode(viper.GetBool("mysql.log_mode"))
	return db, err
}

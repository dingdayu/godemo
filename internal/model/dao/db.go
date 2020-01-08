// dao 所有定义的模型操作方法，均在此包里。
package dao

import (
	"fmt"
	"os"
	"sync"
	"time"

	"demo/pkg/config"
	"demo/pkg/log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var once sync.Once
var masterDB *gorm.DB

func Init() {
	var err error
	once.Do(func() {
		masterDB, err = gorm.Open("mysql", config.GetString("db.master.mysql"))
		if err == nil {
			fmt.Println("\033[1;30;42m[info]\033[0m db [master] connect success")
		} else {
			fmt.Printf("\033[1;30;41m[error]\033[0m db [master] connect error: %s", err.Error())
			os.Exit(1)
		}
		masterDB.SetLogger(log.New())
		masterDB.LogMode(config.GetBool("db.master.log_model"))
		masterDB.DB().SetConnMaxLifetime(time.Minute * time.Duration(config.GetInt("db.master.conn_max_lifetime")))
		masterDB.DB().SetMaxIdleConns(config.GetInt("db.master.max_idle_conn"))
		masterDB.DB().SetMaxOpenConns(config.GetInt("db.master.max_open_conn"))

	})
}

func GetMasterDB() *gorm.DB {
	return masterDB
}

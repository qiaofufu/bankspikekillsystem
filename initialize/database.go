package initialize

import (
	"bank-product-spike-system/global"
	"bank-product-spike-system/models"

	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDatabaseTables() {
	err := global.DB.AutoMigrate(
		&models.Admin{},
		&models.User{},
		&models.Product{},
		&models.SpikeActivity{},
		&models.FilterTable{},
		&models.FilterRoles{},
		&models.ActivityUser{},
		&models.ActivityGuidanceInformation{},
		&models.SlideShow{},
		&models.Article{},
		&models.Order{},
		&models.SpikeResult{},
		&models.ActivityAttention{},
		&models.Node{},
		&models.Edge{},
	)
	if err != nil {
		global.Logger.Fatal("数据库自动迁移失败 " + err.Error())
	}
}

func InitDatabase() {
	host := viper.GetString("database.host")
	port := viper.GetString("database.port")
	username := viper.GetString("database.username")
	password := viper.GetString("database.password")
	dbname := viper.GetString("database.dbname")
	charset := viper.GetString("database.charset")
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=Local", username, password, host, port, dbname, charset)
	db, err := gorm.Open(mysql.Open(args), &gorm.Config{})
	if err != nil {
		panic("Database init error: " + err.Error())
	}
	global.DB = db
}

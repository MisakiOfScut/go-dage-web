package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

var GSqlDB *sqlx.DB

const (
	CHARSET = "utf8"
)

type Config struct {
	MysqlHost string `json:"mysql_host"`
	MysqlPort string `json:"mysql_port"`
	MysqlUser string `json:"mysql_user"`
	MysqlPwd  string `json:"mysql_pwd"`
	DataBase  string `json:"database"`
}

func InitMysql(c *Config){
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", c.MysqlUser, c.MysqlPwd,
		c.MysqlHost,
		c.MysqlPort,
		c.DataBase, CHARSET)

	var err error
	// 打开连接失败
	db, err := sqlx.Open("mysql", dbDSN)
	// defer MysqlDb.Close();
	if err != nil {
		zap.S().Panicf("dbDSN: %s, 数据源配置不正确: %v", dbDSN, err)
	}

	// 最大连接数
	db.SetMaxOpenConns(100)
	// 闲置连接数
	db.SetMaxIdleConns(20)
	// 最大连接周期
	db.SetConnMaxLifetime(100 * time.Second)

	if err = db.Ping(); nil != err {
		zap.S().Panicf("数据库链接失败:%s", err.Error())
	}

	GSqlDB = db
}


func Close(){
	defer GSqlDB.Close()
}
package main

import (
	"encoding/json"
	"flag"
	"github.com/natefinch/lumberjack"
	"go-dage-web/app"
	"go-dage-web/app/utils/mysql"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
)

func main() {
	configFile := flag.String("path", "config.json", "config path")
	flag.Parse()
	config, err := ReadConfigFromFile(*configFile)
	if err != nil {
		zap.S().Fatal(err)
	}

	// 获取编码器,NewJSONEncoder()输出json格式，NewConsoleEncoder()输出普通文本格式
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 指定时间格式
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	// 文件writeSyncer
	fileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./app.log", // 日志文件存放目录
		MaxSize:    1024,        // 文件大小限制,单位MB
		MaxBackups: 10,          // 最大保留日志文件数量
		MaxAge:     30,          // 日志文件保留天数
		Compress:   false,       // 是否压缩处理
	})
	fileCore := zapcore.NewCore(encoder, fileWriteSyncer, zapcore.DebugLevel)

	logger := zap.New(fileCore, zap.AddCaller()) // AddCaller()为显示文件名和行号
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	mysql.InitMysql(&config.Config)
	defer mysql.Close()

	app.StartHttpServer()
}

type Config struct {
	mysql.Config
}

func ReadConfigFromFile(filePath string) (Config, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}
	config := Config{}
	if err := json.Unmarshal(bytes, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}

package middlewares

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"workshop01/db/mongo"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var (
	collection string
	err        error
	fLogger    *lumberjack.Logger
	//errLog     *log.Logger
)

type (
	//Logger struct logger for echo
	Logger struct {
		Time       time.Time `bson:"time" json:"time"`
		Lv         string    `bson:"level" json:"level"`
		Prefix     string    `bson:"prefix" json:"prefix"`
		Message    string    `bson:"-" json:"message"`
		Data       ctxLogger `bson:"data" json:"data"`
		Collection string    `bson:"-"`
	}

	ctxLogger struct {
		ID  string      `json:"id" bson:"id"`
		Req interface{} `json:"req" bson:"request"`
		Res interface{} `json:"res" bson:"response"`
	}
)

//Init log
func init() {

	/*
		e, err := os.OpenFile("./log/foo.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

		if err != nil {
			fmt.Printf("error opening file: %v", err)
			os.Exit(1)
		}

		log.New(e, "", log.Ldate|log.Ltime)
	*/

	/*
		    errLog = log.New(e, "", log.Ldate|log.Ltime)
		    errLog.SetOutput(&lumberjack.Logger{
		      Filename:   "./foo.log",
		      MaxSize:    1, // megabytes after which new file is created
		      MaxBackups: 3, // number of backups
		      MaxAge:     28, //days
			})
	*/

	//some time shutdown database, you will need this.
	fLogger = &lumberjack.Logger{
		Filename: "./log/foo.log", //viper.GetString("logger.filename"), //config.DataConfig.Logger
		MaxSize:  650,             // megabytes
		MaxAge:   15,              //days
		Compress: true,            // disabled by default
	}

	// Set output to use your custom package logger and using middleware
	//log.SetPrefix("prefix")
	//log.SetOutput(&Logger{Collection: "logger"}) //middleware log to mongodb

	fmt.Println("init logs...")
}

func (lg *Logger) Write(logByte []byte) (n int, err error) {

	//fmt.Println("write logs >>>")
	fmt.Printf("\n logBytes: %s\n", logByte)
	fmt.Printf("\n lg: %v\n", &lg)

	err = json.Unmarshal(logByte, &lg)
	if err != nil {
		fmt.Println("\n err, json >>>", err)
		return
	}
	go func() {
		fLogger.Write(logByte)
	}()
	err = json.NewDecoder(strings.NewReader(lg.Message)).Decode(&lg.Data)
	go func() {
		conn := mongo.MgoManager().Copy()
		defer conn.Close()

		if err2 := conn.DB("document").C(lg.Collection).Insert(&lg); err2 != nil {
			fmt.Printf("\ntime:%s,%s\n", time.Now(), lg.Message)
		} else {
			fmt.Printf("\n not err, time:%s,%s\n", time.Now(), lg.Message)
		}
	}()
	return len(logByte), nil
}

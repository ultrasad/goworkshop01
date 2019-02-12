package middlewares

import (
	"encoding/json"
	"fmt"
	"path/filepath"
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
	//Logger struct logger from go log
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

	// Logs struct log from echo
	Logs struct {
		ID           string    `json:"id" bson:"id"`
		Time         time.Time `json:"time" json:"time"`
		RemoteIP     string    `json:"remote_ip" bson:"remote_ip"`
		Host         string    `json:"host" bson:"host"`
		Method       string    `json:"method" bson:"method"`
		URI          string    `json:"uri" bson:"uri"`
		Status       int       `json:"status" bson:"status"`
		Latency      int       `json:"latency" bson:"latency"`
		LatencyHuman string    `json:"latency_human" bson:"latency_human"`
		BytesIn      int       `json:"bytes_in" bson:"bytes_in"`
		BytesOut     int       `json:"bytes_out" bson:"bytes_out"`
		Collection   string    `bson:"-"`
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
		//Filename: "./log/foo.log", //viper.GetString("logger.filename"), //config.DataConfig.Logger
		Filename: filepath.Join("./logs", "logfile.log"),
		MaxSize:  650,  // megabytes
		MaxAge:   15,   //days
		Compress: true, // disabled by default
	}

	// Set output to use your custom package logger and using middleware
	//log.SetPrefix("prefix")
	//log.SetOutput(&Logger{Collection: "logger"}) //middleware log to mongodb

	fmt.Println("init logs...")
}

func (lg *Logger) Write(logByte []byte) (n int, err error) {

	//fmt.Println("write logs >>>")
	//fmt.Printf("log input: %s\n", logByte)
	//fmt.Printf("logger: %v\n", &lg)

	//fmt.Printf("log Logger: %s\n", logByte)

	err = json.Unmarshal(logByte, &lg)
	if err != nil {
		fmt.Println("\n err Logger, json Unmarshal >>>", err)
		return
	}

	//fmt.Printf("\n log Message: %s\n", lg.Message)
	//fmt.Printf("\n &lg.Data: %v\n", &lg.Data)
	//fmt.Printf("\n &lg.Message: %v\n", strings.NewReader(lg.Message))

	// Not Write Log File
	/*
		go func() {
			fLogger.Write(logByte)
		}()
	*/

	//enc := json.NewEncoder(os.Stdout)
	//d := map[string]int{"apple": 5, "lettuce": 7}
	//enc.Encode(d)

	err = json.NewDecoder(strings.NewReader(lg.Message)).Decode(&lg.Data)

	if err != nil {
		//fmt.Println("\n err json decode >>>", err)
		return
	}

	go func() {
		conn := mongo.MgoManager().Copy()
		defer conn.Close()

		if err := conn.DB("document").C(lg.Collection).Insert(&lg); err != nil {
			fmt.Printf("\n err time:%s,%s\n", time.Now(), lg.Message)
		} else {
			//fmt.Printf("\n not err, time:%s\n", time.Now())
		}
	}()
	return len(logByte), nil
}

func (lg *Logs) Write(logEcho []byte) (n int, err error) {

	//fmt.Printf("log Logs : %s\n", logEcho)

	err = json.Unmarshal(logEcho, &lg)
	if err != nil {
		fmt.Println("\n err Logs, json Unmarshal >>>", err)
		return
	}

	go func() {
		fLogger.Write(logEcho)
	}()

	//fmt.Printf("\n &lg Logs: %#v\n", &lg)

	go func() {
		conn := mongo.MgoManager().Copy()
		defer conn.Close()

		if err := conn.DB("document").C(lg.Collection).Insert(&lg); err != nil {
			fmt.Printf("\n err Logs time:%s, %s\n", time.Now(), err)
		} else {
			//fmt.Printf("\n not err, time:%s\n", time.Now())
		}
	}()

	return len(logEcho), nil
}

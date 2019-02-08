package main

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"workshop01/db/mongo"
	"workshop01/routers"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//logrus
var logger = logrus.New()

// Info top-level functions to wrap Logrus
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Debug top-level functions to wrap Logrus
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

//WithConn is with connect
func WithConn(conn net.Conn) *logrus.Entry {
	addr := "unknown"
	if conn != nil {
		addr = conn.RemoteAddr().String()
	}
	return logger.WithField("addr", addr)
}

//RequestFields is request
func RequestFields(req *http.Request) logrus.Fields {
	return logrus.Fields{"userAgent": req.UserAgent()}
}

//WithRequest is with request
func WithRequest(req *http.Request) *logrus.Entry {
	return logger.WithFields(RequestFields(req))
}

var (
	con net.Conn
)

var (
	req *http.Request
)

func proxy(w http.ResponseWriter, r *http.Request) {
	u, _ := url.Parse("http://myip.ninja")
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}

func main() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.ReadInConfig()
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetDefault("port", "8082")

	/*
		var (
			buf    bytes.Buffer
			logger = log.New(&buf, "logger: ", log.Lshortfile)
		)

		logger.Print("Hello, log file!")
		fmt.Print(&buf)
	*/

	//http.HandleFunc("/", proxy)
	//log.Fatal(http.ListenAndServe(":8082", nil))

	e := echo.New()

	//Start Mongo Connect
	mongo.ConnectMgo()

	// Start Router
	routers.Init(e)

}

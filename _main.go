package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
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

func main() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.ReadInConfig()
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("port", "8082")

	mongoHost := viper.GetString("mongo.host")
	mongoUser := viper.GetString("mongo.user")
	mongoPass := viper.GetString("mongo.pass")
	port := fmt.Sprintf(":%v", viper.GetString("port"))

	connString := fmt.Sprintf("%v:%v@%v", mongoUser, mongoPass, mongoHost)
	conn, err := mgo.Dial(connString)
	if err != nil {
		log.Printf("dial mongodb server with connection string %q: %v", connString, err)
		return
	}

	h := &handler{
		mongo: conn,
		db:    "document",
		col:   "todo",
	}

	e := echo.New()

	//Start Mongo Connect
	//mongo.ConnectMgo()
	//fmt.Println("conn main => ", conn)
	/*h := &handler{
		mongo: conn,
		db:    "document",
		col:   "todo",
	}*/

	// Start Router
	//routers.Init(e)

	//logrus
	//logger.Info("Some info. Earth is not flat")

	WithConn(con).Info("Connected")
	logger.WithFields(logrus.Fields{
		"animal": "walrus",
	}).Info("A walrus appears")

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	fmt.Printf("Starting...")

	e.GET("/todos", h.list)
	e.POST("/todos", h.create)
	e.GET("/todos/:id", h.view)
	e.PUT("/todos/:id", h.done)
	e.DELETE("/todos/:id", h.remove)

	e.Logger.Fatal(e.Start(port))

}

type todo struct {
	ID    bson.ObjectId `json:"id" bson:"_id"`
	Topic string        `json:"topic" bson:"topic"`
	Done  bool          `json:"done" bson:"done"`
}

type handler struct {
	mongo *mgo.Session
	db    string
	col   string
}

func (h *handler) list(c echo.Context) error {
	conn := h.mongo.Copy()
	defer conn.Close()
	var ts []todo
	if err := conn.DB(h.db).C(h.col).Find(nil).All(&ts); err != nil {
		return err
	}
	c.JSON(http.StatusOK, ts)
	return nil
}

func (h *handler) view(c echo.Context) error {
	conn := h.mongo.Copy()
	defer conn.Close()
	var t todo
	id := bson.ObjectIdHex(c.Param("id"))

	if err := conn.DB(h.db).C(h.col).FindId(id).One(&t); err != nil {
		return err
	}
	c.JSON(http.StatusOK, t)
	return nil
}

func (h *handler) create(c echo.Context) error {
	id := bson.NewObjectId()
	var t todo
	if err := c.Bind(&t); err != nil {
		return err
	}
	t.ID = id
	t.Done = false

	conn := h.mongo.Copy()
	defer conn.Close()
	if err := conn.DB(h.db).C(h.col).Insert(t); err != nil {
		return err
	}

	c.JSON(http.StatusOK, t)
	return nil
}

func (h *handler) done(c echo.Context) error {
	conn := h.mongo.Copy()
	defer conn.Close()
	var t todo
	id := bson.ObjectIdHex(c.Param("id"))

	if err := conn.DB(h.db).C(h.col).FindId(id).One(&t); err != nil {
		return err
	}

	//update more done
	if err := c.Bind(&t); err != nil {
		return err
	}

	t.Done = true
	if err := conn.DB(h.db).C(h.col).UpdateId(id, t); err != nil {
		return err
	}
	c.JSON(http.StatusOK, t)
	return nil
}

func (h *handler) remove(c echo.Context) error {
	conn := h.mongo.Copy()
	defer conn.Close()
	id := bson.ObjectIdHex(c.Param("id"))

	if err := conn.DB(h.db).C(h.col).RemoveId(id); err != nil {
		return err
	}
	c.JSON(http.StatusOK, echo.Map{
		"result": "success",
	})
	return nil
}

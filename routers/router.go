package routers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"workshop01/controllers"
	"workshop01/middlewares"

	//"github.com/google/logger"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

/*
//Counter Handle
type Counter struct {
	n int
}

func (ctr *Counter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctr.n++
	fmt.Fprintf(w, "counter = %d\n", ctr.n)
}
*/

/*
func middlewareOne(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("from middleware one")
		return next(c)
	}
}

func middlewareTwo(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("from middleware two")
		return next(c)
	}
}

func middlewareSomething(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("from middleware something", w, r)
		next.ServeHTTP(w, r)
	})
}

func customGenerator() string {
	fmt.Printf("Custom UUID")
	return "UUID"
}
*/

//Init func
func Init(e *echo.Echo) {

	port := fmt.Sprintf(":%v", viper.GetString("port"))

	// Middleware
	//e.Use(middleware.Logger())
	//e.Use(middleware.Recover())

	//Set output to use your custom package logger and using middleware
	//log.SetPrefix("prefix")
	//log.SetOutput(&middlewares.Logger{Collection: "logger"}) //middleware log to mongodb

	/*
		f, err := os.OpenFile("./log/foo.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()
		log.SetOutput(f)
	*/

	//middlewares.Init()
	//log.SetPrefix("prefix")
	logger := log.New(os.Stderr, "", 0)
	logger.SetOutput(&middlewares.Logger{Collection: "logger"}) //middleware log to mongodb

	// Logger request, response
	e.Use(middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{
		Handler: func(c echo.Context, reqBody, resBody []byte) {

			reqB := ""
			if len(reqBody) > 0 {
				reqB = string(reqBody)
			}

			logger.Printf(`{"id":"%s","req":%s,"res":%s}`, c.Response().Header().Get(echo.HeaderXRequestID), reqB, resBody)
			// logger.Init must be called first to setup logger
			//logger.Init("./log")
			//logger.Info("Failed to find player! uid=%d plid=%d cmd=%s xxx=%d", 1234, 678942, "getplayer", 102020101)
			//logger.Warn("Failed to parse protocol! uid=%d plid=%d cmd=%s", 1234, 678942, "getplayer")
			//logger.Infof(`{"id":"%s","req":%s,"res":%s}`, c.Response().Header().Get(echo.HeaderXRequestID), reqB, resBody)
			//fmt.Printf(`{"id":"%s","req":%s,"res":%s}`, c.Response().Header().Get(echo.HeaderXRequestID), reqB, resBody)
		},
	}))

	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "hanajung" && password == "secret" {
			return true, nil
		}
		return false, nil
	}))

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	fmt.Printf("Starting...")

	// CustomRequestID
	/*
		e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
			Generator: func() string {
				return customGenerator()
			},
		}))
	*/

	//ctr := new(Counter)
	//http.Handle("/counter", ctr)

	//e.Use(middlewareOne)
	//e.Use(middlewareTwo)

	//e.Use(echo.WrapMiddleware(middlewareSomething))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	e.GET("/todos", controllers.List)
	e.POST("/todos", controllers.Create)
	e.GET("/todos/:id", controllers.View)
	e.PUT("/todos/:id", controllers.Done)
	e.DELETE("/todos/:id", controllers.Delete)

	//GoRoutine
	e.GET("/hello", controllers.CallHelloRoutine)

	e.Logger.Fatal(e.Start(port))
}

package routers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
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

	//OK
	//middlewares.Init()
	//log.SetPrefix("prefix")
	logger := log.New(os.Stderr, "", 0)
	logger.SetOutput(&middlewares.Logger{Collection: "logger"}) //middleware log to mongodb

	// Logger request, response
	e.Use(middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{
		Handler: func(c echo.Context, reqBody, resBody []byte) {

			reqB := "\"\""
			if len(reqBody) > 0 {
				reqB = string(reqBody)
			}

			logger.Printf(`{"time": "%s", "message": "{}", "level": "info","data": {"id":"%s","req":%s,"res":%s}}`, time.Now().Format("2006-01-02T15:04:05Z"), c.Response().Header().Get(echo.HeaderXRequestID), reqB, resBody)
		},
	}))

	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "hanajung" && password == "secret" {
			return true, nil
		}
		return false, nil
	}))

	/*e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))*/

	//RFC3339local := "2006-01-02T15:04:05Z"
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_custom}","id":"${id}","remote_ip":"${remote_ip}","host":"${host}",` +
			`"method":"${method}","uri":"${uri}","status":${status}, "latency":${latency},` +
			`"latency_human":"${latency_human}","bytes_in":${bytes_in},` +
			`"bytes_out":${bytes_out}}` + "\r\n",
		CustomTimeFormat: "2006-01-02T15:04:05Z",
		Output:           &middlewares.Logs{Collection: "logs"},
	}))

	fmt.Printf("Starting...")

	/*
		url1, err := url.Parse("www.yahoo.com")
		if err != nil {
			e.Logger.Fatal(err)
		}
		targets := []*middleware.ProxyTarget{&middleware.ProxyTarget{URL: url1}}
		e.Group("/myblog", middleware.ProxyWithConfig(middleware.ProxyConfig{
			Balancer: &middleware.RobinBalancer{
				Targets: targets,
			},
		}))
	*/

	/*
		url1, err := url.Parse("http://localhost:1323")
		if err != nil {
			e.Logger.Fatal(err)
		}
		url2, err := url.Parse("http://localhost:1323")
		if err != nil {
			e.Logger.Fatal(err)
		}
		targets := []*middleware.ProxyTarget{
			{
				URL: url1,
			},
			{
				URL: url2,
			},
		}

		e.Group("/users", middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))
	*/

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

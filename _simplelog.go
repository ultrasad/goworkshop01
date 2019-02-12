package main

type logger struct {
    ch chan []byte
}

func (l logger) Write(p []byte) (n int, err error) {
    l.ch <- p
    return len(p), nil
}

func main() {
    ch, wg := make(chan []byte), new(sync.WaitGroup)
    wg.Add(1)

    go runLogger(ch, wg)
    runWebServer(ch)

    close(ch)
    wg.Wait()
}

func runLogger(ch chan []byte) {
    defer wg.Done()

    for entry := range ch {
        // handle entry - send to SQL, or group them and send them in batches,
        // handle errors, perform retries, etc.
    }

    // handle any cleanup, flush/close db connection, etc.
}

// Note you will need to implement graceful shutdown: https://echo.labstack.com/cookbook/graceful-shutdown
// so that runLogger() can also complete gracefully.
func runWebServer(ch chan []byte) {
    // ...
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: `{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}","host":"${host}",` +
				`"method":"${method}","uri":"${uri}","status":${status}, "latency":${latency},` +
				`"latency_human":"${latency_human}","bytes_in":${bytes_in},` +
				`"bytes_out":${bytes_out}}` + "\r\n",
			Output: logger{ch},
		}))
    // ...
}
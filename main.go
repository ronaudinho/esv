package main

import (
	"log"
	"math/rand"
	"runtime"
	"strconv"
	"time"

	"github.com/arl/statsviz"
	"github.com/arl/statsviz/websocket"
	"github.com/labstack/echo/v4"
)

// statsviz routes
const (
	statsvizRoot = "/admin/debug/stats/"
	statsvizWs   = "/admin/debug/stats/ws"
)

func main() {
	// force the GC to work to make the plots "move".
	// reference: examples in https://github.com/arl/statsviz/blob/master/_example
	go work()

	e := echo.New()

	// "*" to allow getting files necessary for UI stylings
	e.GET(statsvizRoot+"*", echo.WrapHandler(statsviz.IndexAtRoot(statsvizRoot)))
	e.GET(statsvizWs, ws)

	e.Logger.Fatal(e.Start(":1323"))
}

// echo requires explicitly upgrading to websocket,
// hence copying statsviz.Ws to make it work here.
// reference: https://echo.labstack.com/cookbook/websocket
func ws(c echo.Context) error {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("ws: upgrade error:", err)
		return err
	}
	defer ws.Close()

	tick := time.NewTicker(time.Second)
	defer tick.Stop()

	// unexported struct in statsviz
	// reference: https://github.com/arl/statsviz/blob/master/statsviz.go
	var stats struct {
		Mem          runtime.MemStats
		NumGoroutine int
	}
	for {
		select {
		case <-tick.C:
			runtime.ReadMemStats(&stats.Mem)
			stats.NumGoroutine = runtime.NumGoroutine()
			if err := ws.WriteJSON(stats); err != nil {
				return err
			}
		}
	}
}

// reference: https://github.com/arl/statsviz/blob/master/_example/work.go
func work() {
	m := map[string][]byte{}

	for {
		b := make([]byte, 512+rand.Intn(16*1024))
		m[strconv.Itoa(len(m)%(10*100))] = b

		if len(m)%(10*100) == 0 {
			m = make(map[string][]byte)
		}

		time.Sleep(10 * time.Millisecond)
	}
}

package main

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/arl/statsviz"
	"github.com/labstack/echo/v4"
)

// statsviz routes
const (
	statsvizRoot = "/debug/statsviz/"
)

func main() {
	// force the GC to work to make the plots "move".
	// reference: examples in https://github.com/arl/statsviz/blob/master/_example
	go work()

	e := echo.New()

	mux := http.NewServeMux()
	statsviz.Register(mux)

	// "*" to allow getting files necessary for UI stylings
	e.GET(statsvizRoot+"*", echo.WrapHandler(mux))

	e.Logger.Fatal(e.Start(":1323"))
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

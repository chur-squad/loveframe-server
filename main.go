package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"runtime"
	"syscall"
	"github.com/labstack/echo/v4"
	"github.com/chur-squad/loveframe-server/handler"
	_error "github.com/chur-squad/loveframe-server/error"
)

var (
	echoAddr = flag.String("echo_addr", ":8080", "echo address")
	tls      = flag.Bool("tls", false, "tls")
)

// main is a point for start.
func main() {
	// parse args
	flag.Parse()
	fmt.Println(flag.Args())
	// check go processor
	fmt.Printf("logical cpu = %d, go processor = %d, \n", runtime.NumCPU(), runtime.GOMAXPROCS(0))

	// initialize handler config
	cfg, err := createConfigForHandler()
	if err != nil {
		panic(err)
	}

	// initialize handler
	h, err := handler.NewHandler(handler.WithConfig(cfg))
	if err != nil {
		panic(err)
	}

	// initialize echo server
	e, err := initEchoServer(h)
	if err != nil {
		panic(err)
	}
	// create signal handler
	signalCtx, signalStop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	// serve echo
	fmt.Println("[echo] start server")
	go func() {
		if *tls {
			e.StartAutoTLS(*echoAddr)
		} else {
			e.Start(*echoAddr)
		}
	}()
	//server loop
	// wait for signal
	<-signalCtx.Done()
	signalStop()

	// stop echo gracefully
	fmt.Println("[echo] stop server gracefully")
	e.Shutdown(context.Background())
}

func initEchoServer(h *handler.Handler) (*echo.Echo, error) {
	// create echo server
	e := echo.New()

	// add route
	if err := addRoute(e, h); err != nil {
		return nil, _error.WrapError(err)
	}

	return e, nil
}


func createConfigForHandler() (*handler.Config, error) {
	// generate handler config
	cfg := &handler.Config{
		/*
		CdnEndpoint: 				env.GetCdnEndpoint(),
		GroupSalt:                  env.GetGroupCodeSalt(),
		UserSalt:                   env.GetUserCodeSalt(),
		*/
	}
	return cfg, nil
}
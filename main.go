package main

import (
	"context"
	"flag"
	"fmt"
	_error "github.com/chur-squad/loveframe-server/error"
	"github.com/chur-squad/loveframe-server/handler"
	"github.com/chur-squad/loveframe-server/env"
	"github.com/labstack/echo/v4"
	"os/signal"
	"runtime"
	"syscall"
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

	h.Mysql.AddUser(1, "Kim jaehyun")

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
		CdnEndpoint: env.GetCdnEndpoint(),
		GroupSalt:   env.GetGroupCodeSalt(),
		UserSalt:    env.GetUserCodeSalt(),
	}

	var (
		dbUsername, dbPassword, dbHost, dbPort, dbName string
	)

	dbUsername = env.GetDatabaseUsername()
	dbPassword = env.GetDatabasePassword()
	dbHost = env.GetDatabaseHost()
	dbPort = env.GetDatabasePort()
	dbName = env.GetDatabaseName()

	// add database DSN config
	cfg.MysqlDSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4,utf8&timeout=1s&autocommit=true&parseTime=true",
		dbUsername, dbPassword, dbHost, dbPort, dbName)

	return cfg, nil
}

/*
func connectDatabase() (database *sql.DB, err error) {
	user := env.GetDatabaseUsername()
	password := env.GetDatabasePassword()
	dbname := env.GetDatabaseName()
	dbconfig := fmt.Sprintf("%s:%s@/%s", user, password, dbname)
	
	db, err := sql.Open("mysql", dbconfig)
	if err != nil {
		return nil, _error.WrapError(err)
	}
	db.SetConnMaxLifetime(time.Minute * 1)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	fmt.Print()
	return db, nil
}
*/
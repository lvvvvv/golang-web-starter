package main

import (
	"flag"
	"fmt"
	"os"

	"github.acsdev.net/starters/golang-web-starter/server"

	"github.com/kataras/iris"
)

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func main() {

	var configFile string

	cfg := server.NewDefaultConfig()

	// Parse flags
	flag.BoolVar(&cfg.Debug, "D", true, "Enable Debug logging.")
	flag.StringVar(&cfg.Addr, "a", ":3000", "Network host to listen on.")
	flag.StringVar(&cfg.Database.Driver, "db-driver", "postgres", "database driver")
	flag.IntVar(&cfg.Database.Pool, "db-pool", 16, "database pool")
	flag.StringVar(&cfg.Database.ConnString, "db-conn", "dbname=whiteraven_dev user=postgres password=123123123 host=10.20.30.11 sslmode=disable", "database connection string")

	flag.StringVar(&configFile, "c", "config/config.toml", "Configuration file.")

	flag.Parse()

	// Parse config if given
	if configFile != "" {
		fmt.Printf("loading config file %s\n", configFile)
		if err := cfg.ParseToml(configFile); err != nil {
			die("Failure to parse config, %v", err)
		}
	}

	err := server.SetupDatabase(cfg)
	if err != nil {
		die("failure to setup database %v\n", err)
	}

	server.SetupIris(cfg)

	// register the routes & the public API
	server.RegisterRoutes()
	server.RegisterAPI()

	// start the server
	iris.Listen(cfg.Addr)
}

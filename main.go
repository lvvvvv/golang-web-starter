package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kataras/iris"

	"github.acsdev.net/starters/golang-web-starter/config"
	"github.acsdev.net/starters/golang-web-starter/web"
)

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func main() {

	cfg := config.Default()

	// Parse flags
	flag.BoolVar(&cfg.Debug, "D", true, "Enable Debug logging.")

	flag.BoolVar(&cfg.Database.Enabled, "db-enabled", true, "database enabled")
	flag.StringVar(&cfg.Database.Driver, "db-driver", "postgres", "database driver")
	flag.IntVar(&cfg.Database.Pool, "db-pool", 16, "database pool")
	flag.StringVar(&cfg.Database.Addr, "db-addr", "localhost:5432", "database addr")
	flag.StringVar(&cfg.Database.Name, "db-name", "cozy_dev", "database name")
	flag.StringVar(&cfg.Database.User, "db-user", "sykipper", "database user")
	flag.StringVar(&cfg.Database.Pass, "db-pass", "123123123", "database pass")
	flag.BoolVar(&cfg.Database.SSL, "db-ssl", false, "database tls")

	flag.BoolVar(&cfg.Redis.Enabled, "redis-enabled", true, "redis enabled")
	flag.StringVar(&cfg.Redis.Addr, "redis-addr", "192.168.0.133:6379", "redis addr")
	flag.StringVar(&cfg.Redis.Pass, "redis-pass", "123123123", "redis pass")

	flag.StringVar(&cfg.Web.Addr, "web-addr", ":3000", "web server addr listen on.")
	flag.StringVar(&cfg.Web.RootDir, "web-rootdir", ".", "web server root dir")
	flag.StringVar(&cfg.Web.PublicDir, "web-publicdir", "public", "web server public dir")
	flag.StringVar(&cfg.Web.JwtSecret, "web-jwtsecret", "ad4wgsfserqfg2", "web jwt secret")

	flag.Parse()

	config.ConfigureIris(&cfg)

	//connect db, return nil if not enabled
	_ = config.ConfigureDatabase(&cfg)

	//connect redis, return nil if not enabled
	_ = config.ConfigureRedis(&cfg)

	// register the routes & the public API
	web.RegisterRoutes(&cfg)

	//start the web server
	iris.Listen(cfg.Web.Addr)

}

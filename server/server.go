package server

import (
	"database/sql"
	"os"
	"time"

	"github.acsdev.net/starters/golang-web-starter/server/api"
	"github.acsdev.net/starters/golang-web-starter/server/routes"
	"github.com/iris-contrib/logger"
	"github.com/iris-contrib/middleware/i18n"
	mLogger "github.com/iris-contrib/middleware/logger"
	"github.com/iris-contrib/middleware/recovery"
	"github.com/kataras/iris"
	"github.com/rs/cors"
	"gopkg.in/mgutz/dat.v1"
	"gopkg.in/mgutz/dat.v1/kvs"
	"gopkg.in/mgutz/dat.v1/sqlx-runner"
)

//SetupDatabase setup db connection and put into iris context
func SetupDatabase(cfg *Config) error {
	// create a normal database connection through database/sql
	db, err := sql.Open(cfg.Database.Driver, cfg.Database.ConnString)
	if err != nil {
		return err
	}
	// ensures the database can be pinged with an exponential backoff (15 min)
	runner.MustPing(db)

	// set to reasonable values for production
	// db.SetMaxIdleConns(4)
	db.SetMaxOpenConns(cfg.Database.Pool)

	if cfg.Database.CacheEnabled {
		// Redis: namespace is the prefix for keys and should be unique
		store, err := kvs.NewRedisStore(cfg.Database.CacheNamespace, cfg.Database.CacheRedisAddr, cfg.Database.CacheRedisPwd)
		if err != nil {
			return err
		}
		runner.SetCache(store)
	}

	// set this to enable interpolation
	dat.EnableInterpolation = true

	// set to check things like sessions closing.
	// Should be disabled in production/release builds.
	dat.Strict = false

	// Log any query over 10ms as warnings. (optional)
	runner.LogQueriesThreshold = 10 * time.Millisecond

	ddb := runner.NewDB(db, cfg.Database.Driver)

	iris.UseFunc(func(c *iris.Context) {
		c.Set("db", ddb)
		c.Next()
	})

	return nil
}

func SetupIris(cfg *Config) error {

	// set the template engine
	iris.Config.Render.Template.Engine = iris.DefaultEngine    // or iris.DefaultEngine
	iris.Config.Render.Template.Layout = "layouts/layout.html" // = ./templates/layouts/layout.html

	iris.Use(recovery.New(os.Stderr)) // optional

	iris.Use(i18n.New(i18n.Config{Default: "en-US",
		Languages: map[string]string{
			"en-US": "./locales/locale_en-US.ini",
			"zh-CN": "./locales/locale_zh-CN.ini"}}))

	// register cors middleware
	iris.Use(cors.New(cors.Options{}))

	// iris.UseTemplate(html.New(html.Config{Layout: "layout.html"})).Directory("./templates", ".html")
	// set the favicon
	iris.Favicon("./webapp/public/images/favicon.ico")

	// set static folder(s)
	iris.Static("/public", "./webapp/public", 1)

	// logCfg := logger.Config{
	// 	EnableColors: false, //enable it to enable colors for all, disable colors by iris.Logger.ResetColors(), defaults to false
	// 	// Status displays status code
	// 	Status: true,
	// 	// IP displays request's remote address
	// 	IP: true,
	// 	// Method displays the http method
	// 	Method: true,
	// 	// Path displays the request path
	// 	Path: true,
	// }

	// set the global middlewares
	iris.Use(mLogger.New(logger.New(logger.DefaultConfig())))

	// set the custom errors
	iris.OnError(iris.StatusNotFound, func(ctx *iris.Context) {
		ctx.Render("errors/404.html", iris.Map{"Title": iris.StatusText(iris.StatusNotFound)})
	})

	iris.OnError(iris.StatusInternalServerError, func(ctx *iris.Context) {

		ctx.Render("errors/500.html", nil)
	})

	iris.UseFunc(func(c *iris.Context) {
		c.Set("cfg", cfg)
		c.Next()
	})

	return nil
}

func SetupWorker(cfg *Config) error {
	return nil
}

func RegisterRoutes() {
	// register index using a 'Handler'
	iris.Handle("GET", "/", routes.Index())

	// this is other way to declare a route
	// using a 'HandlerFunc'
	iris.Get("/about", routes.About)

	// Dynamic route

	iris.Get("/profile/:username", routes.Profile)("user-profile")
	// user-profile is the custom,optional, route's Name: with this we can use the {{ url "user-profile" $username}} inside userlist.html

	iris.Get("/all", routes.UserList)
}

func RegisterAPI() {
	// this is other way to declare routes using the 'API'
	iris.API("/users", api.UserAPI{})
}

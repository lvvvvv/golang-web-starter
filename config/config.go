package config

import (
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/iris-contrib/logger"
	"github.com/iris-contrib/middleware/cors"
	"github.com/iris-contrib/middleware/i18n"
	"github.com/iris-contrib/middleware/recovery"
	"github.com/kataras/iris"
	"gopkg.in/pg.v4"

	mLogger "github.com/iris-contrib/middleware/logger"
)

type Database struct {
	Enabled    bool
	Driver     string
	ConnString string
	Addr       string
	Name       string
	User       string
	Pass       string
	SSL        bool
	Pool       int
}

type Redis struct {
	Enabled bool
	Addr    string
	Pass    string
}

type Web struct {
	Addr      string
	RootDir   string
	PublicDir string
	JwtSecret string
}

// Config Server Configuration
type Config struct {
	Debug bool

	Database Database
	Redis    Redis
	Web      Web
}

func Default() Config {
	return Config{
		Debug:    true,
		Database: DefaultDatabase(),
		Redis:    DefaultRedis(),
		Web:      DefaultWeb(),
	}
}

func DefaultDatabase() Database {
	return Database{
		Enabled: true,
		Driver:  "postgres",
		Addr:    "127.0.0.1",
		Name:    "starter",
		User:    "postgres",
		Pass:    "123123123",
		SSL:     false,
		Pool:    16,
	}
}

func DefaultRedis() Redis {
	return Redis{
		Enabled: true,
		Addr:    "127.0.0.1:6397",
		Pass:    "123123123",
	}
}

func DefaultWeb() Web {
	return Web{
		Addr:      "0.0.0.0:8822",
		RootDir:   ".",
		PublicDir: "./public",
		JwtSecret: "aasdgwsgevwsfsetw",
	}
}

//ConfigureDatabase setup db connection
func ConfigureDatabase(cfg *Config) *pg.DB {
	if !cfg.Database.Enabled {
		return nil
	}
	db := pg.Connect(&pg.Options{
		Addr:     cfg.Database.Addr,
		Database: cfg.Database.Name,
		User:     cfg.Database.User,
		Password: cfg.Database.Pass,
		PoolSize: cfg.Database.Pool,
	})

	iris.UseFunc(func(c *iris.Context) {
		c.Set("db", db)
		c.Next()
	})

	return db
}

//ConfigureRedis setup redis connection
func ConfigureRedis(cfg *Config) *redis.Pool {
	if !cfg.Redis.Enabled {
		return nil
	}
	rds := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", cfg.Redis.Addr)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", cfg.Redis.Pass); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	iris.UseFunc(func(c *iris.Context) {
		c.Set("rds", rds)
		c.Next()
	})

	return rds

}

//ConfigureIris setup iris opts
func ConfigureIris(cfg *Config) {

	// set the template engine
	iris.Config.Render.Template.Engine = iris.DefaultEngine                       // or iris.DefaultEngine
	iris.Config.Render.Template.Layout = cfg.Web.RootDir + "/layouts/layout.html" // = ./templates/layouts/layout.html

	iris.Use(recovery.New(os.Stderr)) // optional

	iris.Use(i18n.New(i18n.Config{Default: "en-US",
		Languages: map[string]string{
			"en-US": cfg.Web.RootDir + "/locales/locale_en-US.ini",
			"zh-CN": cfg.Web.RootDir + "/locales/locale_zh-CN.ini"}}))

	// register cors middleware
	iris.Use(cors.New(cors.Options{}))

	// iris.UseTemplate(html.New(html.Config{Layout: "layout.html"})).Directory("./templates", ".html")
	// set the favicon
	iris.Favicon(cfg.Web.PublicDir + "/images/favicon.ico")

	// set static folder(s)
	iris.Static("/public", cfg.Web.PublicDir, 1)

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

}

package web

import (
	"github.acsdev.net/starters/golang-web-starter/config"
	"github.acsdev.net/starters/golang-web-starter/web/api"
	"github.acsdev.net/starters/golang-web-starter/web/routes"
	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
)

func RegisterRoutes(cfg *config.Config) {

	secret := []byte(cfg.Web.JwtSecret)

	//Jwt middleware
	myJwtMiddleware := jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	//put jwtSecret into iris context
	iris.UseFunc(func(c *iris.Context) {
		c.Set("jwtSecret", secret)
		c.Next()
	})

	// register index using a 'Handler'
	iris.Handle("GET", "/", routes.Index())

	// this is other way to declare a route
	// using a 'HandlerFunc'
	iris.Get("/about", routes.About)

	// Dynamic route

	iris.Get("/profile/:username", routes.Profile)("user-profile")
	// user-profile is the custom,optional, route's Name: with this we can use the {{ url "user-profile" $username}} inside userlist.html

	iris.Get("/all", routes.UserList)

	//api
	iris.Post("/account/login", api.UserLogin)
	iris.Post("/account/register", api.UserRegister)
	iris.Get("/account/logout", myJwtMiddleware.Serve, api.UserLogout)

}

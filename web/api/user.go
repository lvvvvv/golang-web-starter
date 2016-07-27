package api

import (
	"fmt"
	"time"

	"github.acsdev.net/starters/golang-web-starter/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"github.com/syklevin/cozy/server/errors"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/pg.v4"
)

type inputParams struct {
	Login string `json:"login"`
	Pwd   string `json:"pwd"`
	Role  string `json:"role"`
	Mode  string `json:"mode"`
}

func UserLogin(ctx *iris.Context) {

	db, _ := ctx.Get("db").(*pg.DB)
	jwtSecret, _ := ctx.Get("jwtSecret").([]byte)

	var input inputParams
	var err error

	err = ctx.ReadJSON(&input)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, errors.NewApiError(errors.InvalidParams, "failure to parse json"))
		return
	}

	if input.Login == "" || input.Pwd == "" {
		ctx.JSON(iris.StatusBadRequest, errors.NewApiError(errors.InvalidParams, "login or password should not empty"))
		return
	}

	var user models.User

	err = db.Model(&user).
		Where("email = ?", input.Login).
		Limit(1).
		Select()

	if err != nil {
		ctx.JSON(iris.StatusBadRequest, errors.NewApiError(errors.UserNotFound, "user not found"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Pwd))

	if err != nil {
		ctx.JSON(iris.StatusBadRequest, errors.NewApiError(errors.InvalidParams, "password not match"))
		return
	}

	sess, err := createUserSession(ctx, db, &user)

	cuser := &models.UserJwt{
		SID:  sess.ID,
		ID:   user.ID,
		Name: user.Name,
		Role: user.Role,
	}

	token, err := cuser.GenerateJwt(jwtSecret)

	if err != nil {
		ctx.JSON(iris.StatusBadRequest, errors.NewApiError(errors.InvalidParams, "password not match"))
		return
	}

	ctx.JSON(iris.StatusOK, map[string]string{"jwt": token})

}

func UserRegister(ctx *iris.Context) {

	db, _ := ctx.Get("db").(*pg.DB)
	jwtSecret, _ := ctx.Get("jwtSecret").([]byte)

	// pg.SetQueryLogger(log.New(os.Stdout, "", log.LstdFlags))

	var input inputParams
	var err error

	err = ctx.ReadJSON(&input)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, errors.NewApiError(errors.InvalidParams, "failure to parse json"))
		return
	}

	if input.Login == "" || input.Pwd == "" {
		ctx.JSON(iris.StatusBadRequest, errors.NewApiError(errors.InvalidParams, "login or password should not empty"))
		return
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(input.Pwd), bcrypt.DefaultCost)

	if err != nil {
		ctx.JSON(iris.StatusBadRequest, errors.NewApiError(errors.InvalidParams, "password not match"))
		return
	}

	user := models.User{
		Name:      "guest",
		Email:     input.Login,
		Password:  string(hashedPwd),
		Role:      input.Role,
		Status:    models.USER_CREATE,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}

	err = db.Create(&user)

	if err != nil {
		fmt.Printf("step 1 err %+v\n", err)
		ctx.Text(iris.StatusBadRequest, err.Error())
		return
	}

	sess, err := createUserSession(ctx, db, &user)

	if err != nil {
		ctx.Text(iris.StatusBadRequest, err.Error())
		return
	}

	cuser := &models.UserJwt{
		SID:  sess.ID,
		ID:   user.ID,
		Name: user.Name,
		Role: user.Role,
	}

	token, err := cuser.GenerateJwt(jwtSecret)

	if err != nil {
		ctx.JSON(iris.StatusBadRequest, errors.NewApiError(errors.InvalidParams, "password not match"))
		return
	}

	ctx.JSON(iris.StatusOK, map[string]string{"jwt": token})

}

func UserLogout(ctx *iris.Context) {
	jwtToken, _ := ctx.Get("jwt").(*jwt.Token)
	db, _ := ctx.Get("db").(*pg.DB)

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
		cuser, _ := claims["user"].(*models.UserJwt)

		sess := models.UserSession{
			ID:        cuser.SID,
			Status:    models.SESSION_DEACTIVE,
			UpdatedAt: time.Now(),
		}

		err := db.Update(&sess)

		if err != nil {
			iris.Logger.Warningf("failure to logout user %v\n", err)
		}

	} else {
		iris.Logger.Warningf("invalid jwt\n")
	}

	ctx.Text(iris.StatusOK, "")

}

func createUserSession(ctx *iris.Context, db *pg.DB, user *models.User) (*models.UserSession, error) {
	userAgent := ctx.RequestHeader("User-Agent")
	referer := ctx.RequestHeader("Referer")
	clientIp := ctx.RemoteAddr()

	sess := models.UserSession{
		UserID:    user.ID,
		Referer:   referer,
		UserAgent: userAgent,
		ClientIp:  clientIp,
		Status:    models.SESSION_ACTIVE,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}

	err := db.Create(&sess)

	if err != nil {
		return nil, err
	}
	return &sess, nil

}

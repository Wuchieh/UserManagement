package server

import (
	"context"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"time"
)

func Run(ip string, s string, expiredTime int, mode string) (err error) {
	switch strings.ToLower(mode) {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "test":
		gin.SetMode(gin.TestMode)
	}
	r = gin.Default()
	store := cookie.NewStore([]byte(s))
	store.Options(
		sessions.Options{
			MaxAge: expiredTime, // 單位秒
			Path:   "/",
		})
	// session 的名稱會在 browser 變成 cookie 的 key
	r.Use(sessions.Sessions("us", store))

	routeInit(r)
	go cornTab()
	srv = &http.Server{
		Addr:    ip,
		Handler: r,
	}
	log.Println("Server Running!")
	log.Println("RunMode: ", gin.Mode())
	err = srv.ListenAndServe()
	return
}

func Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	return err
}

func reboot() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Println(err)
		return
	}
}

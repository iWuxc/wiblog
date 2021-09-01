package main

import (
	"fmt"
	"path/filepath"
	"wiblog/pkg/conf"
	"wiblog/pkg/core/wiblog"
	"wiblog/pkg/core/wiblog/admin"
	"wiblog/pkg/core/wiblog/file"
	"wiblog/pkg/core/wiblog/page"
	"wiblog/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Hi, it's, App", conf.Conf.WiBlogApp.Name)
	endRun := make(chan error, 1)
	runHttpServer(endRun)
	fmt.Println(<-endRun)
}

func runHttpServer(endRun chan error) {
	if !conf.Conf.WiBlogApp.EnableHTTP {
		return
	}

	//set mode
	if conf.Conf.RunMode == conf.ModeProd {
		gin.SetMode(gin.ReleaseMode)
	}

	//create a http
	e := gin.Default()

	e.Use(middleware.UserMiddleware())
	e.Use(middleware.SessionMiddleware(middleware.SessionOpts{
		Name:   "su",
		Secure: conf.Conf.RunMode == conf.ModeProd,
		Secret: []byte("ZGlzvcmUoMTAsICI="),
	}))

	//static files, page
	staticFile := filepath.Join(conf.WorkDir, "assets")
	e.Static("static", staticFile)

	//// static files
	file.RegisterRouters(e)

	// frontend pages
	page.RegisterRoutes(e)

	// unauthz api
	admin.RegisterRoutes(e)

	group := e.Group("/admin", wiblog.AuthFilter)
	{
		page.RegisterRoutesAuthz(group)
		//admin.RegisterRoutesAuthz(group)
	}

	address := fmt.Sprintf(":%d", conf.Conf.WiBlogApp.HTTPPort)
	go func() {
		_ = e.Run(address)
	}()
	fmt.Println("HTTP server running on: " + address)
}

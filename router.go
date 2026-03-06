package main

import (
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(cors.Default())
	router.Use(func(c *gin.Context) {
		c.Set("APP_SITE", os.Getenv("APP_SITE"))
		c.Set("CF_TOKEN", os.Getenv("CF_TOKEN"))
		c.Next()
	})
	router.LoadHTMLGlob("templates/*")
	router.Static("/assets", "./assets")
	router.GET("/script.min.js", func(c *gin.Context) {
		c.File("./assets/script.min.js")
	})
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.File("./assets/favicon.ico")
	})
	router.GET("/style.less", func(c *gin.Context) {
		c.File("./assets/style.less")
	})
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"site":      c.MustGet("APP_SITE"),
			"cf_token":  c.MustGet("CF_TOKEN"),
			"themeList": themeList,
		})
	})
	router.GET("/@:name", handleCounter)
	router.GET("/heart-beat", func(c *gin.Context) {
		c.String(http.StatusOK, "alive")
	})
	return router
}

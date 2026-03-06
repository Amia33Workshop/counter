package main

import (
	"math/rand"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
)

type getCountImageOptions struct {
	Count     int
	Theme     string
	Padding   int
	Offset    int
	Align     string
	Scale     float64
	Pixelated bool
	DarkMode  string
	Prefix    int
}

var validName = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func handleCounter(c *gin.Context) {
	name := c.Param("name")
	if !validName.MatchString(name) || len(name) > 32 {
		c.String(http.StatusBadRequest, "Invalid name")
		return
	}
	theme := c.DefaultQuery("theme", "moebooru")
	if theme == "random" {
		keys := make([]string, 0, len(themeList))
		for k := range themeList {
			keys = append(keys, k)
		}
		theme = keys[rand.Intn(len(keys))]
	} else if _, ok := themeList[theme]; !ok {
		theme = "moebooru"
	}
	padding, err := strconv.Atoi(c.DefaultQuery("padding", "7"))
	if err != nil || padding < 0 || padding > 16 {
		padding = 7
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < -500 || offset > 500 {
		offset = 0
	}
	align := c.DefaultQuery("align", "top")
	if align != "top" && align != "center" && align != "bottom" {
		align = "top"
	}
	scale, err := strconv.ParseFloat(c.DefaultQuery("scale", "1"), 64)
	if err != nil || scale < 0.1 || scale > 2 {
		scale = 1
	}
	pixelated := c.DefaultQuery("pixelated", "1")
	if pixelated != "0" && pixelated != "1" {
		pixelated = "1"
	}
	darkmode := c.DefaultQuery("darkmode", "auto")
	if darkmode != "0" && darkmode != "1" && darkmode != "auto" {
		darkmode = "auto"
	}
	num, err := strconv.Atoi(c.DefaultQuery("num", "0"))
	if err != nil || num < 0 || num > 1e15 {
		num = 0
	}
	prefix, err := strconv.Atoi(c.DefaultQuery("prefix", "-1"))
	if err != nil || prefix < -1 || prefix > 999999 {
		prefix = -1
	}
	counter, err := getCountByName(name, num)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting count")
		return
	}
	if name == "demo" {
		c.Header("cache-control", "max-age=31536000")
	} else {
		c.Header("cache-control", "max-age=0, no-cache, no-store, must-revalidate")
	}
	c.Header("content-type", "image/svg+xml")
	svg := getCountImage(getCountImageOptions{
		Count:     counter.Num,
		Theme:     theme,
		Padding:   padding,
		Offset:    offset,
		Align:     align,
		Scale:     scale,
		Pixelated: pixelated == "1",
		DarkMode:  darkmode,
		Prefix:    prefix,
	})
	c.String(http.StatusOK, svg)
	LogDebug("ip: " + strconv.Quote(c.ClientIP()) + ", name: " + strconv.Quote(name) + ", theme: " + strconv.Quote(theme) + ", user-agent: " + strconv.Quote(c.Request.UserAgent()) + ", referrer: " + strconv.Quote(c.Request.Referer()))
}

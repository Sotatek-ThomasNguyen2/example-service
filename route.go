package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Router struct {
	route *gin.Engine
	u     *User
}

func authenBasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
		if len(auth) != 2 || auth[0] != "Basic" {
			c.AbortWithStatusJSON(401, "Unauthorized")
			return
		}
		c.Set("now", time.Now())
		c.Set("$apikey", auth[1])
		c.Next()
	}
}

func (r *Router) router() {
	r.route.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, ready to serve!",
		})
	})
	r.route.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.route.GET("user/user-partner", authenBasicAuth(), r.handleListUserPartner)
}

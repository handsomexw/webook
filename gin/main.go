package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	server := gin.Default()
	server.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello",
		})
	})

	server.POST("/post", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello post方法")
	})

	server.GET("/users/:name", func(ctx *gin.Context) {
		name := ctx.Param("name")
		ctx.JSON(http.StatusOK, gin.H{
			"name": name,
		})
	})

	server.GET("/views/*.html", func(ctx *gin.Context) {
		val := ctx.Param(".html")
		ctx.JSON(http.StatusOK, gin.H{
			"value": val,
		})
	})

	server.GET("/keys", func(ctx *gin.Context) {
		key := ctx.Query("id")
		ctx.JSON(http.StatusOK, gin.H{
			"id": key,
		})
	})

	server.Run(":9090")
}

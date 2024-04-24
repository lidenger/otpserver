package admin_ui

import (
	"embed"
	"github.com/gin-gonic/gin"
	"strings"
)

//go:embed static/*
var staticFS embed.FS

func Initialize(g *gin.Engine) {

	g.GET("/assets/:filepath", func(c *gin.Context) {
		c.Writer.WriteHeader(200)
		filePath := c.Param("filepath")
		b, _ := staticFS.ReadFile("static/assets/" + filePath)
		contentType := ""
		if strings.HasSuffix(filePath, ".js") {
			contentType = "text/javascript"
		}
		if strings.HasSuffix(filePath, ".css") {
			contentType = "text/css"
		}
		c.Header("Content-Type", contentType)
		_, _ = c.Writer.Write(b)
		c.Writer.Flush()
	})

	g.GET("/", func(c *gin.Context) {
		c.Writer.WriteHeader(200)
		c.Header("Content-Type", "text/html; charset=utf-8")
		b, _ := staticFS.ReadFile("static/index.html")
		_, _ = c.Writer.Write(b)
		c.Writer.Flush()
	})

	g.GET("/favicon.ico", func(c *gin.Context) {
		c.Writer.WriteHeader(200)
		b, _ := staticFS.ReadFile("static/favicon.ico")
		_, _ = c.Writer.Write(b)
		c.Writer.Flush()
	})
}

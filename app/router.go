package app

import (
	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go-dage-web/app/service/node_service"
	"go-dage-web/app/service/script_service"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func StartHttpServer() {
	router := gin.New()
	router.Use(ginzap.Ginzap(zap.L(), time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(zap.L(), true))

	// html
	router.LoadHTMLGlob("web/static/template/*")
	router.GET("/index", func(context *gin.Context) {
		context.HTML(http.StatusOK, "home.html", nil)
	})
	router.GET("/edit.html", func(context *gin.Context) {
		context.HTML(http.StatusOK, "edit.html", nil)
	})

	// file
	router.StaticFS("/static",http.Dir("./web/static/"))
	router.StaticFS("/temp",http.Dir("./web/static/temp/"))

	// data api
	api := router.Group("/api")
	{
		node := api.Group("/node")
		{
			node.GET("/getAll", node_service.GetList)
			node.GET("/add", node_service.AddNode)
			node.GET("/del", node_service.DelNode)
		}

		script := api.Group("/script")
		{
			script.GET("/getAll", script_service.GetAll)
			script.POST("/add", script_service.AddVersion)
			script.POST("/check", script_service.Check)
			script.POST("/publish", script_service.Publish)
		}
	}

	err := router.Run(":8000")
	if err != nil {
		panic(err)
	}
}

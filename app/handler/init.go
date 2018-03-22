package handler

import (
	"api/app/entity"
	"api/config"
	"master/middleware"
	"master/utilities/ulog"

	"github.com/gin-gonic/gin"
)

func InitEngine(conf *config.Config) *gin.Engine {
	if conf.App.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware(conf.App.Whitelist))
	r.Use(gin.LoggerWithWriter(ulog.Logger().Request))
	r.LoadHTMLGlob("public/*.html")
	r.Static("static", "./public/static")

	if conf.App.Debug {
		r.Use(gin.Logger())
	}

	indexHandler := indexHandler{
		Category: entity.NewCategory(),
	}
	groupIndex := r.Group("")
	{
		GET(groupIndex, "", indexHandler.Index)
	}
	// Product
	productHandler := productHandler{
		productEntity: entity.NewProduct(),
	}
	groupProduct := r.Group("/products")
	{
		GET(groupProduct, "/:slug", productHandler.GetDetail)
		GET(groupProduct, "", productHandler.GetList)
		POST(groupProduct, "", productHandler.Create)
	}

	return r
}

func GET(group *gin.RouterGroup, relativePath string, f func(*gin.Context)) {
	route(group, "GET", relativePath, f)
}

func POST(group *gin.RouterGroup, relativePath string, f func(*gin.Context)) {
	route(group, "POST", relativePath, f)
}

func PUT(group *gin.RouterGroup, relativePath string, f func(*gin.Context)) {
	route(group, "PUT", relativePath, f)
}

func DELETE(group *gin.RouterGroup, relativePath string, f func(*gin.Context)) {
	route(group, "DELETE", relativePath, f)
}

func route(group *gin.RouterGroup, method string, relativePath string, f func(*gin.Context)) {
	hanld := middleware.ErrorHandler(group.BasePath() + relativePath)
	switch method {
	case "POST":
		group.POST(relativePath, hanld, f)
	case "GET":
		group.GET(relativePath, hanld, f)
	case "PUT":
		group.PUT(relativePath, hanld, f)
	case "DELETE":
		group.DELETE(relativePath, hanld, f)
	}
}
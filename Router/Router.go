package Router

import (
	"../Handlers"
	"github.com/julienschmidt/httprouter"
)

func Routes() (routes *httprouter.Router) {
	routes = httprouter.New()
	routes.GET("/", Handlers.HomeHandler)
	routes.GET("/upload", Handlers.UploadHandler)
	routes.GET("/uploads/:name", Handlers.UploadNameHandler)
	routes.POST("/upload", Handlers.UploadHandler)
	routes.GET("/list", Handlers.ListHandler)
	return
}

package routes

import (
	"database/sql"
	controller "gogeek/controllers"
	"log"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(db *sql.DB) *gin.Engine {

	// Get a user resource
	router := gin.Default()
	log.Println("Server started on: http://localhost:8000")
	router.LoadHTMLGlob("templates/*/*.html")
	router.Static("/assets", "./assets")
	router.GET("/", controller.Loginview)
	router.GET("/newregister", controller.Register)
	router.POST("/savedata", controller.Regsave)
	router.POST("/login", controller.Login)
	router.GET("/index", controller.Index)
	router.GET("/productdetail", controller.Productdetail)
	router.GET("/reviewrating", controller.Reviewrating)
	router.POST("/ratingsave", controller.Reviewratingsave)
	router.POST("/addorders", controller.Orders)
	router.GET("/myorders", controller.Ordersview)
	router.GET("/logout", controller.DeleteSession)

	return router
}

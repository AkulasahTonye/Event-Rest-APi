package routes

import (
	"example.com/api-testing/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.GET("/events", getEvents)    // GET PUT, POST, PATCH, DELETE
	server.GET("/events/:id", getEvent) // events/1, events/5 etc....

	Authenticate := server.Group("/")
	Authenticate.Use(middlewares.Authenticate)
	Authenticate.POST("/events", createEvent)
	Authenticate.PUT("/events/:id", updateEvent)
	Authenticate.DELETE("/events/:id", deleteEvent)

	Authenticate.POST("/events/:id/register", registerForEvent)
	Authenticate.DELETE("/events/:id/register", cancelRegistration)

	server.POST("/signup", signup)
	server.POST("/login", login)

}

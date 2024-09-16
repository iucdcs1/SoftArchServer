package routes

import (
	"github.com/gin-gonic/gin"
	"softarch/controllers"
)

func SetupRouter(route *gin.Engine) {
	public := route.Group("/messages")
	{
		public.GET("", controllers.GetMessages)
		public.POST("", controllers.SendMessage)
		public.GET("/count", controllers.GetCount)
	}
}

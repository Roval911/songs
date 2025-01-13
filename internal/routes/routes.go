package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "songs/docs"
	"songs/internal/hanlers"
)

func SetupRouter(songHandler *hanlers.Handler) *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	// Swagger документация
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Публичные маршруты
	public := router.Group("/api/v1")
	{
		public.GET("/songs", songHandler.GetSongs)
		public.GET("/song/:id/lyrics", songHandler.GetLyrics)
		public.POST("/song", songHandler.AddSong)
		public.DELETE("/song/:id", songHandler.DeleteSong)
		public.PUT("/song/:id", songHandler.UpdateSong)
	}

	return router
}

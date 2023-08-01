package routes

import (
	"arhive-manager-go/src/config"
	"database/sql"
	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func Run(db *sql.DB) {
	router.Use(JSONMiddleware())
	router.SetTrustedProxies(nil)
	getRoutes(db)
	err := router.Run(":" + config.ApplicationConfig.RestPort)
	if err != nil {
		panic("Failed to run")
	}
}

func JSONMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	}
}

func getRoutes(db *sql.DB) {
	v1 := router.Group("/v1")
	AddRoutes(v1, db)
}

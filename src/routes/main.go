package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func Run(db *sql.DB) {
	router.Use(JSONMiddleware())
	_ = router.SetTrustedProxies(nil)
	getRoutes(db)

	go func() {
		if err := router.Run(":8080"); err != nil {
			panic(err)
		}
	}()
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

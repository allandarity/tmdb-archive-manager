package routes

import (
	. "arhive-manager-go/src/database"
	"arhive-manager-go/src/service"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type ErrorMessage struct {
	Type    int    `json:"type"`
	Message string `json:"message"`
	Error   error  `json:"error"`
}

func AddRoutes(rg *gin.RouterGroup, db *sql.DB) {
	users := rg.Group("/tv")

	users.GET("/:tmdbId", func(c *gin.Context) {
		id := c.Param("tmdbId")

		tmdbId, err := GetContentEntryById(db, "tmdb_id", id)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, ErrorMessage{
				Type:    500,
				Message: "Failed to retrieve id",
				Error:   err,
			})
			return
		}

		externalId, err := service.ConvertID(tmdbId)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, ErrorMessage{
				Type:    500,
				Message: "Failed to call out to convert id",
				Error:   err,
			})
			return
		}

		posterId, err := GetPosterByContentId(db, tmdbId.Id)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, ErrorMessage{
				Type:    500,
				Message: "Failed to manage poster",
				Error:   err,
			})
			return
		}
		posterPath, err := service.GetPoster(posterId, tmdbId)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, ErrorMessage{
				Type:    500,
				Message: "Failed to get poster from tmdb",
				Error:   err,
			})
			return
		}

		posterImageData, err := service.FetchImage(posterPath.FilePath)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, ErrorMessage{
				Type:    500,
				Message: "Failed to get poster byte data",
				Error:   err,
			})
			return
		}

		UpdateImdbIdForGivenRow(db, externalId)

		//TODO: Make this not 3 sep updates
		UpdatePosterDatabaseEntry(db, tmdbId.Id, "poster_data", posterImageData)
		UpdatePosterDatabaseEntry(db, tmdbId.Id, "poster_url", posterPath.FilePath)
		UpdatePosterDatabaseEntry(db, tmdbId.Id, "content_id", tmdbId.Id)

		content, err := GetContentEntryById(db, "tmdb_id", id) //getting the id again for the updated value

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, ErrorMessage{
				Type:    500,
				Message: "Failed to find returning id in database",
				Error:   err,
			})
			return
		}

		c.JSON(http.StatusOK, content)
	})
}

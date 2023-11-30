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
	tvEndpoint := rg.Group("/tv")

	tvEndpoint.GET("/:tmdbId", func(c *gin.Context) {

		//TODO: add this queue stuff as the below will populate 5 entries but i think ratelimited on pulling the poster
		//https://webdevstation.com/posts/simple-queue-implementation-in-golang/

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

		if err != nil || externalId.ImdbId == "" {
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

	tvEndpoint.GET("/random/:count", func(c *gin.Context) {
		count := c.Param("count")
		contents, err := GetSpecifiedAmountOfContent(db, count)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, ErrorMessage{
				Type:    500,
				Message: "Failed to retrieve content",
				Error:   err,
			})
			return
		}
		c.JSON(http.StatusOK, contents)
	})

}

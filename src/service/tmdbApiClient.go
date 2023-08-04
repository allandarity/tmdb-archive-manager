package service

import (
	"arhive-manager-go/src/config"
	"arhive-manager-go/src/database"
	. "arhive-manager-go/src/model"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func ConvertID(content database.ContentEntry) (TmdbTvItem, error) {

	if content.ImdbId.Valid {
		fmt.Printf("IMDB id for input %d already exists, using saved information \n", content.TmdbId)
		return TmdbTvItem{
			ID:     content.TmdbId,
			ImdbId: content.ImdbId.String,
		}, nil
	}

	authKey := fmt.Sprintf("Bearer %s", config.ApplicationConfig.TmdbKey)
	endpoint := fmt.Sprintf("https://api.themoviedb.org/3/tv/%d/external_ids", content.TmdbId)
	req, err := http.NewRequest("GET", endpoint, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", authKey)

	if err != nil {
		fmt.Println(err)
		return TmdbTvItem{}, err
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return TmdbTvItem{}, err
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return TmdbTvItem{}, err
	}

	fmt.Println(data)
	var entry TmdbTvItem
	err = json.Unmarshal(data, &entry)
	if err != nil {
		fmt.Println(err)
		return TmdbTvItem{}, err
	}

	return entry, nil
}

func GetPoster(poster database.PosterEntry, content database.ContentEntry) (TmdbTvPoster, error) {
	if poster.PosterUrl != "" {
		fmt.Printf("Poster info for input %d already exists, using saved information \n", poster.Id)
		return TmdbTvPoster{
			FilePath: poster.PosterUrl,
		}, nil
	}
	log.Println("Starting request for POSTERS")
	authKey := fmt.Sprintf("Bearer %s", config.ApplicationConfig.TmdbKey)
	log.Println(authKey)
	log.Println(content.TmdbId)
	endpoint := fmt.Sprintf("https://api.themoviedb.org/3/tv/%d/images?include_image_language=en", content.TmdbId)
	req, err := http.NewRequest("GET", endpoint, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", authKey)

	if err != nil {
		fmt.Println(err)
		return TmdbTvPoster{}, err
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return TmdbTvPoster{}, err
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	responseText := string(data)
	fmt.Println("Response Body:", responseText)
	if err != nil {
		fmt.Println(err)
		return TmdbTvPoster{}, err
	}
	log.Println("Starting unmarshall for POSTER")
	var entry TmdbTvPosterContainer
	err = json.Unmarshal(data, &entry)
	if err != nil {
		fmt.Println(err)
		return TmdbTvPoster{}, err
	}

	log.Println(entry)
	return entry.Posters[0], nil
}

func FetchImage(url string) ([]byte, error) {

	var posterUrl = fmt.Sprintf("%s%s", "https://www.themoviedb.org/t/p/original", url)
	log.Println(fmt.Sprintf("URL provided for poster is: %s", posterUrl))
	resp, err := http.Get(posterUrl)
	if err != nil {
		fmt.Println("Error fetching image:", err)
		return nil, err
	}
	defer resp.Body.Close()

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading image data:", err)
		return nil, err
	}
	return imageData, nil
}

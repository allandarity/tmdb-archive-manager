package service

import (
	"arhive-manager-go/src/config"
	"arhive-manager-go/src/database"
	. "arhive-manager-go/src/model"
	"encoding/json"
	"fmt"
	"io"
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

	var entry TmdbTvItem
	err = json.Unmarshal(data, &entry)
	if err != nil {
		fmt.Println(err)
		return TmdbTvItem{}, err
	}

	return entry, nil
}

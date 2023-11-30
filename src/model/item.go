package model

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
)

type TvShow struct {
	ID         int     `json:"id"`
	Title      string  `json:"original_name"`
	Popularity float64 `json:"popularity"`
}

func ReadDownloadedFileForTvShow(path string) (map[int]TvShow, int, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
		return nil, 0, err
	}
	defer file.Close()
	byteVal, err := os.ReadFile(file.Name())
	if err != nil {
		log.Fatal(err)
		return nil, 0, err
	}

	showMap := make(map[int]TvShow)
	var tvShow TvShow
	for _, line := range bytes.Split(byteVal, []byte{'\n'}) {
		if err := json.Unmarshal(line, &tvShow); err != nil {
			log.Printf("Reached the end of parsing: error is %s\n", err)
		}
		showMap[tvShow.ID] = tvShow
	}

	return showMap, tvShow.ID, nil
}

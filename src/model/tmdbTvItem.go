package model

type TmdbTvItem struct {
	ID     int    `json:"id"`
	ImdbId string `json:"imdb_id"`
	TvdbId int    `json:"tvdb_id"`
}

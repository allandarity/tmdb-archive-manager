package model

type TmdbTvItem struct {
	ID     int    `json:"id"`
	ImdbId string `json:"imdb_id"`
	TvdbId int    `json:"tvdb_id"`
}

type TmdbTvPosterContainer struct {
	Posters []TmdbTvPoster `json:"posters"`
}

type TmdbTvPoster struct {
	FilePath string `json:"file_path"`
	Height   int32  `json:"height"`
	Width    int32  `json:"width"`
}

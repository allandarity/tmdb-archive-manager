package database

import (
	"arhive-manager-go/src/config"
	. "arhive-manager-go/src/model"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/lib/pq"
	"log"
)

type ContentEntry struct {
	Id             int            `json:"id"`
	Title          string         `json:"title"`
	TmdbId         int            `json:"tmdbId"`
	TmdbPopularity float64        `json:"tmdbPopularity"`
	ImdbId         sql.NullString `json:"imdbId"`
	ImdbPopularity string         `json:"imdbPopularity"`
	ContentType    string         `json:"contentType"`
}

func buildConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@db:%s/%s?sslmode=disable",
		config.ApplicationConfig.Database.Username,
		config.ApplicationConfig.Database.Password,
		config.ApplicationConfig.Database.Port,
		config.ApplicationConfig.Database.Schema)
}

func OpenConnection() (*sql.DB, error) {
	db, err := sql.Open("postgres", buildConnectionString())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetFinalTmdbTvEntry(db *sql.DB) (int, error) {
	query := "select tmdb_id from content order by tmdb_id desc limit 1"

	var show TvShow

	row := db.QueryRow(query)
	if err := row.Scan(&show.ID); err != nil {
		if err == sql.ErrNoRows {
			return show.ID, err
		}
		return show.ID, err
	}
	return show.ID, nil
}

func UpdateImdbIdForGivenRow(db *sql.DB, item TmdbTvItem) {
	stmt, err := db.Prepare("update content set imdb_id = $1 where tmdb_id = $2")

	if err != nil {
		log.Fatal("Error preparing update statement:", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(item.ImdbId, item.ID)
	if err != nil {
		log.Fatal("Error executing update statement:", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal("Error getting rows affected:", err)
	}

	fmt.Printf("Rows affected: %d\n", rowsAffected)

}

func GetEntryByTmdbId(db *sql.DB, id string) (ContentEntry, error) {
	query := fmt.Sprintf("select * from content where tmdb_id = %s", id)
	var content ContentEntry

	row := db.QueryRow(query)
	var imdbPopularityNull sql.NullFloat64
	if err := row.Scan(&content.Id, &content.Title, &content.TmdbId,
		&content.TmdbPopularity, &content.ImdbId, &imdbPopularityNull, &content.ContentType); err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			return content, err
		}
		log.Println(err)
		return content, err
	}
	return content, nil
}

func BatchInsertTmdbTVData(db *sql.DB, showMap map[int]TvShow) {

	log.Printf("Starting insert with map size %d\n", len(showMap))

	txn, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := txn.Prepare(pq.CopyIn("content", "tmdb_id", "tmdb_popularity", "title", "content_type"))
	if err != nil {
		log.Fatal(err)
	}

	for _, show := range showMap {
		if (show.Title) != "" {
			floatValue := fmt.Sprintf("%.2f", show.Popularity)
			_, err = stmt.Exec(show.ID, floatValue, show.Title, "tv")
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}

	err = stmt.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = txn.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func RunMigration() (*migrate.Migrate, error) {

	m, err := migrate.New(
		"file://db/migrations", buildConnectionString())

	if err != nil {
		log.Println(err)
	}
	if err := m.Up(); err != nil {
		log.Println(err)
	}

	return m, err
}

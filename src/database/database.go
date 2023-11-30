package database

import (
	"arhive-manager-go/src/config"
	. "arhive-manager-go/src/model"
	"arhive-manager-go/src/util"
	"database/sql"
	"errors"
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

type PosterEntry struct {
	Id         int    `json:"id"`
	ContentId  int    `json:"contentId"`
	PosterUrl  string `json:"posterUrl"`
	PosterData []byte `json:"posterData"`
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
		if errors.Is(err, sql.ErrNoRows) {
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

type Entry interface{}

func UpdatePosterDatabaseEntry(db *sql.DB, id int, column string, entry Entry) {
	stmt, err := db.Prepare(fmt.Sprintf("update poster set %s = $1 where content_id = $2", column))

	if err != nil {
		log.Fatal("Error preparing update statement:", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(entry, id)
	if err != nil {
		log.Fatal("Error executing update statement:", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal("Error getting rows affected:", err)
	}

	fmt.Printf("Rows affected: %d\n", rowsAffected)
}

func GetContentEntryById(db *sql.DB, column string, id string) (ContentEntry, error) {
	query := fmt.Sprintf("select * from content where %s = %s", column, id)
	var content ContentEntry

	row := db.QueryRow(query)
	var imdbPopularityNull sql.NullFloat64
	if err := row.Scan(&content.Id, &content.Title, &content.TmdbId,
		&content.TmdbPopularity, &content.ImdbId, &imdbPopularityNull, &content.ContentType); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println(err)
			return content, err
		}
		log.Println(err)
		return content, err
	}
	return content, nil
}

func GetSpecifiedAmountOfContent(db *sql.DB, count string) ([]ContentEntry, error) {

	if !util.IsNumeric(count) {
		log.Println("must be numeric")
		return nil, errors.New("Parameter must be numeric")
	}

	query := fmt.Sprintf("SELECT * FROM content ORDER BY RANDOM() LIMIT %s", count)
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var contentEntries []ContentEntry

	for rows.Next() {
		var content ContentEntry
		var imdbPopularityNull sql.NullFloat64
		err := rows.Scan(&content.Id, &content.Title, &content.TmdbId,
			&content.TmdbPopularity, &content.ImdbId, &imdbPopularityNull, &content.ContentType)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		contentEntries = append(contentEntries, content)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	return contentEntries, nil
}

func GetPosterByContentId(db *sql.DB, id int) (PosterEntry, error) {

	exist, err := doesPosterEntryExist(db, id)

	if err != nil {
		log.Println(err)
		return PosterEntry{}, err
	}

	if exist {
		query := "SELECT * from poster where content_id = $1"

		var poster PosterEntry
		err := db.QueryRow(query, id).Scan(
			&poster.Id, &poster.ContentId, &poster.PosterUrl, &poster.PosterData)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				fmt.Println("No rows found for the given content_id:", id)
				return PosterEntry{}, err
			}
			fmt.Println("Error executing the SELECT query:", err)
			return PosterEntry{}, err
		}

		return poster, nil
	} else {
		entry, err := insertPosterEntry(db, PosterEntry{}, ContentEntry{Id: id})
		if err != nil {
			fmt.Println("Failed creating new poster entry", err)
			return PosterEntry{}, err
		}

		return entry, nil
	}
}

func doesPosterEntryExist(db *sql.DB, id int) (bool, error) {
	query := "SELECT * FROM poster WHERE content_id = $1"

	// Execute the query
	rows, err := db.Query(query, id)
	if err != nil {
		log.Println(err)
		return false, err
	}

	defer rows.Close()

	if rows.Next() {
		return true, nil
	} else {
		return false, nil
	}
}

func insertPosterEntry(db *sql.DB, entry PosterEntry, content ContentEntry) (PosterEntry, error) {
	fmt.Println("Starting Inserting Poster Data")
	query := "INSERT INTO poster (content_id, poster_url, poster_data) VALUES ($1, $2, $3) RETURNING content_id, poster_url, poster_data"

	var poster PosterEntry
	err := db.QueryRow(query, content.Id, entry.PosterUrl, entry.PosterData).Scan(
		&poster.ContentId, &poster.PosterUrl, &poster.PosterData)
	if err != nil {
		fmt.Println("Error executing the INSERT POSTER query:", err)
		return PosterEntry{}, err
	}
	fmt.Println("Finished Inserting Poster Data")
	fmt.Println(poster)
	return poster, nil
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

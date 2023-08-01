package main

import (
	"arhive-manager-go/src/config"
	"arhive-manager-go/src/database"
	"arhive-manager-go/src/model"
	"arhive-manager-go/src/routes"
	"compress/gzip"
	"fmt"
	"github.com/cavaliergopher/grab/v3"
	"io"
	"log"
	"os"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type ContentType string
type DirectoryType string

const (
	MOVIE                 ContentType   = "movie"
	TV                    ContentType   = "tv_series"
	CompressedDirectory   DirectoryType = "./compressed/"
	DecompressedDirectory DirectoryType = "./decompressed/"
)

func buildUrl(contentType ContentType) string {
	return fmt.Sprintf("http://files.tmdb.org/p/exports/%s_ids_%s.json.gz", contentType, getTimeStamp())
}

func getTimeStamp() string {
	return time.Now().Format("01_02_2006")
}

func getFileName(contentType ContentType, directoryType DirectoryType, prefix string) string {
	return fmt.Sprintf("%s%s_%s."+prefix, directoryType, getTimeStamp(), contentType)
}

func fileExists(directoryType DirectoryType, contentType ContentType, prefix string) bool {
	_, err := os.Stat(getFileName(contentType, directoryType, prefix))
	return err == nil
}

func getFile(contentType ContentType) (file string, err error) {
	resp, err := grab.Get(getFileName(contentType, CompressedDirectory, "gz"), buildUrl(contentType))
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	fmt.Println("Download saved to", resp.Filename)

	return resp.Filename, nil
}

func decompressFile(fileName string, contentType ContentType) string {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	defer file.Close()

	reader, err := gzip.NewReader(file)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	defer reader.Close()
	outFile, err := os.Create(getFileName(contentType, DecompressedDirectory, "json"))
	if err != nil {
		log.Fatal(err)
		return ""
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, reader)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	log.Printf("%s decompressed successfully!\n\n", fileName)
	return outFile.Name()
}

func main() {

	config.InitConfig()

	_, err := database.RunMigration()
	if err != nil {
		log.Printf("error on run migration is: %s", err)
	}

	contentType := TV

	if !fileExists(CompressedDirectory, contentType, "gz") {
		_, err := getFile(contentType)
		if err != nil {
			log.Fatal("Failed to find downloaded file")
		}
	} else {
		log.Println("Skipped file check - already exists")
	}

	if !fileExists(DecompressedDirectory, contentType, "json") {
		decompressFile(getFileName(contentType, CompressedDirectory, "gz"), contentType)
	} else {
		log.Println("Skipped file decompression - already exists")
	}

	showMap, lastId, err := model.ReadDownloadedFileForTvShow(getFileName(contentType, DecompressedDirectory, "json"))
	if err != nil {
		log.Fatal("Failed with map")
	}

	db, databaseOpenError := database.OpenConnection()

	if databaseOpenError != nil {
		log.Println(databaseOpenError)
	}

	lastTvShowIdInDb, _ := database.GetFinalTmdbTvEntry(db)

	if lastTvShowIdInDb != lastId {
		log.Printf("Inserting from TMBDID %d\n", lastTvShowIdInDb)
		for k := range showMap {
			if k <= lastTvShowIdInDb {
				delete(showMap, k)
			}
		}
		database.BatchInsertTmdbTVData(db, showMap)
	} else {
		log.Println("Skipped inserting we up to date")
	}

	routes.Run(db)

}

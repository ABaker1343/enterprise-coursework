package repository

import (
	"database/sql"
	"log"
	"os"

	"github.com/mattn/go-sqlite3"
)

//var repo map[string]interface{}
var repo Repository

type Repository struct {
    DB *sql.DB
    Log *log.Logger
}

type Track struct {
    Id string
    Data string
}

func Init() {
    if db, err := sql.Open("sqlite3", "DBtracks"); err != nil {
        repo = Repository {DB : db}
        if logger := log.New(os.Stdout, "", 0); err != nil {
            repo = Repository {DB : db, Log : logger}
        }
    } else {
        log.Fatal("failed to open database")
    }
}

func Create() int {
    const sql = "CREATE TABLE IF NOT EXISTS tracks(id TEXT PRIMARY KEY, data TEXT)"
    if _, err := repo.DB.Exec(sql); err != nil {
        return 0
    } else {
        return -1
    }
}

func AddNewTrack(t Track) int64 {
    const sql = "INSERT INTO tracks (id, data) VALUES (?, ?)"
    if stmt, err := repo.DB.Prepare(sql); err != nil {
        defer stmt.Close()
        if result, err := stmt.Exec(t.Id, t.Data); err != nil {
            if n, err := result.RowsAffected(); err != nil {
                return n;
            } else {
                // failed to check how many rows were affected
                repo.Log.Output(2, "failed to read row count in sql insert statement")
                return -1;
            }
        }
        // failed to execute statement
        repo.Log.Output(2, "failed to execute sql insert statement")
        return -1
    } else {
        // failed to prepare statement
        repo.Log.Output(2, "failed to prepare sql insert statement")
        return -1
    }
}

func GetTrackById(id string) (Track, int) {
    const sql = "SELECT * FROM tracks WHERE id = ?"
    if stmt, err := repo.DB.Prepare(sql); err != nil {
        defer stmt.Close()
        var t Track
        row := stmt.QueryRow(id)
        if err := row.Scan(&t.Id, &t.Data); err != nil {
            return t, 1
        } else {
            return Track{}, 0
        }
    } else {
        // failed to prepare statement
        repo.Log.Output(2, "failed to prepare sql select statement")
        return Track{}, -1
    }
}

func GetAllTracks() ([]Track, int) {
    const sql = "SELECT * FROM tracks"
    if stmt, err := repo.DB.Prepare(sql); err != nil {
        defer stmt.Close()
        tracks := make([]Track, 0)
        rows, err := stmt.Query()
        if err != nil {
            repo.Log.Output(2, "failed to query rows in sql select statement")
        }

        for rows.Next() {
            newTrack := Track{}
            err := rows.Scan(&newTrack.Id, &newTrack.Data)
            if err != nil {
                log.Output(2, "failed to scan row in sql select statement")
                return []Track{}, 0
            }
            tracks = append(tracks, newTrack)
        }
        return tracks, len(tracks)
    } else {
        // failed to prepare statement
        repo.Log.Output(2, "failed to prepare sql select statement")
        return []Track{}, -1
    }
}

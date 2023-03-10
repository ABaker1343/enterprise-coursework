package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

//var repo map[string]interface{}

type Repository struct {
    DB *sql.DB
    Log *log.Logger
}

var repo Repository

type Track struct {
    Id string
    Audio string
}

func Init() {
    if db, err := sql.Open("sqlite3", "DBtracks.db"); err == nil {
        if logger := log.New(os.Stdout, "", 0); err == nil {
            repo = Repository {DB : db, Log : logger}
        }
    } else {
        log.Fatal(err)
    }
}

func Create() int {
    const sql = "CREATE TABLE IF NOT EXISTS tracks(id TEXT PRIMARY KEY, audio TEXT)"
    if _, err := repo.DB.Exec(sql); err == nil {
        return 0
    } else {
        return -1
    }
}

func Clear() int {
    const sql = "DROP TABLE tracks"
    if _, err := repo.DB.Exec(sql); err == nil {
        return 0
    } else {
        return -1
    }
}

func AddNewTrack(t Track) int64 {
    const sql = "INSERT INTO tracks (id, audio) VALUES (?, ?)"
    if stmt, err := repo.DB.Prepare(sql); err == nil {
        defer stmt.Close()
        if result, err := stmt.Exec(t.Id, t.Audio); err == nil {
            if n, err := result.RowsAffected(); err == nil {
                return n;
            } else {
                // failed to check how many rows were affected
                fmt.Println(err)
                repo.Log.Output(2, "failed to read row count in sql insert statement")
                return -1;
            }
        } else {
            if strings.HasPrefix(err.Error(), "UNIQUE constraint failed") {
                return 0
            }
            // failed to execute statement
            repo.Log.Output(2, "failed to execute sql insert statement")
            return -1
        }
    } else {
        // failed to prepare statement
        repo.Log.Output(2, "failed to prepare sql insert statement")
        fmt.Println(err)
        return -1
    }
}

func UpdateTrack(t Track) int64 {
    const sql = "UPDATE tracks SET audio = ? WHERE id = ?"
    if stmt, err := repo.DB.Prepare(sql); err == nil {
        defer stmt.Close()
        if result, err := stmt.Exec(t.Audio, t.Id); err == nil {
            if rowsAffected, err := result.RowsAffected(); err == nil {
                return rowsAffected
            } else {
                // failed to read the rows affected
                repo.Log.Output(2, "failed to read the row count in sql update statement: " + err.Error())
                return -1
            }
        } else {
            // failed to execute sql
            repo.Log.Output(2, "failed to execute sql update statement: " + err.Error())
            return -1
        }
    } else {
        // failed to prepare statement
        repo.Log.Output(2, "failed to prepare sql update statement: " + err.Error())
        return -1
    }
}

func GetTrackById(id string) (Track, int) {
    const sql = "SELECT * FROM tracks WHERE id = ?"
    if stmt, err := repo.DB.Prepare(sql); err == nil {
        defer stmt.Close()
        var t Track
        row := stmt.QueryRow(id)
        if err := row.Scan(&t.Id, &t.Audio); err == nil {
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
    if stmt, err := repo.DB.Prepare(sql); err == nil {
        defer stmt.Close()
        tracks := make([]Track, 0)
        rows, err := stmt.Query()
        if err != nil {
            repo.Log.Output(2, "failed to query rows in sql select statement")
        }

        for rows.Next() {
            newTrack := Track{}
            err := rows.Scan(&newTrack.Id, &newTrack.Audio)
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

func DeleteTrack(id string) int64 {
    const sql = "DELETE FROM tracks WHERE id = ?"
    if stmt, err := repo.DB.Prepare(sql); err == nil {
        defer stmt.Close()
        if result, err := stmt.Exec(id); err == nil {
            if n, err := result.RowsAffected(); err == nil {
                return n
            } else {
                // couldnt check rows affected
                repo.Log.Output(2, err.Error())
                repo.Log.Output(2, "failed to read rows affected in sql delete statement")
                return -1
            }
        } else {
            // failed to execute delete statement
            repo.Log.Output(2, err.Error())
            repo.Log.Output(2, "failed to execute sql delete statement")
            return -1
        }
    } else {
        // failed to prepare statement
        repo.Log.Output(2, err.Error())
        repo.Log.Output(2, "failed to prepare sql delete statement")
        return -1
    }
}

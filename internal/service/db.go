package service

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

var db *sql.DB
var err error

func InitDb()  {
	db, err = sql.Open("sqlite", "./tracker.db")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`
		create table if not exists actions
		(
			start TEXT PRIMARY KEY,
			finish TEXT,
			name TEXT NOT NULL,
			descr TEXT
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func SaveFinish(t time.Time) {
	_, err := db.Exec("update actions set finish = ? where finish is null", t)
	if err != nil {
		log.Fatal(err)
	}
}

func SaveStart(t time.Time, name, descr string) {
	_, err := db.Exec("insert into actions (start, name, descr) values (?, ?, ?)", t, name, descr)
	if err != nil {
		log.Fatal(err)
	}
}

func getActiveJob() (string, int) {
	rows, err := db.Query(`
		select name, coalesce((unixepoch(CURRENT_TIMESTAMP) - unixepoch(start))/60, 0) t
		from actions where finish is null
	`)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
	var n string
	var t int

    for rows.Next() {
        if err := rows.Scan(&n, &t); err != nil {
            log.Fatal(err)
        }
    }
    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }
	return n, t
}

func getMinutesToday() int {
	rows, err := db.Query(`
		select coalesce(sum((unixepoch(coalesce(a.finish, CURRENT_TIMESTAMP)) - unixepoch(a.start)))/60, 0) t
		from actions a
		where date(a.start, 'localtime') = date(CURRENT_TIMESTAMP, 'localtime')
	`)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
	var t int

    for rows.Next() {
        if err := rows.Scan(&t); err != nil {
            log.Fatal(err)
        }
    }
    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }
	return t
}


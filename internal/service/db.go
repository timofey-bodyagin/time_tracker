package service

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

var db *sql.DB
var err error
var LogFile *os.File

func init() {
	initLog()
	db, err = sql.Open("sqlite", "./tracker.db")
	if err != nil {
		log.Fatal(err)
	}
	applyMigrations()
	initSettings()
	log.Println("Database initialized")
}

func SaveFinish(t time.Time) {
	_, err := db.Exec(updateActionFinishSql, t)
	if err != nil {
		log.Fatal(err)
	}
}

func SaveStart(t time.Time, name, descr string) {
	_, err := db.Exec(insertActionSql, t, name, descr)
	if err != nil {
		log.Fatal(err)
	}
}

func RecalcReport() {
	_, err = db.Exec(recalcReportSql)
	if err != nil {
		log.Fatal(err)
	}
}

func getActiveJob() (string, int) {
	rows, err := db.Query(getActiveActionSql)
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
	rows, err := db.Query(getMinutesTodaySql)
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

func GetReportData(dt time.Time) []ReportItem {
	rows, err := db.Query(getReportDataSql, dt)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var items []ReportItem
	for rows.Next() {
		var item ReportItem
		err = rows.Scan(&item.Date, &item.Name, &item.Minutes, &item.Registered)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return items
}

func SetRegistered(name string, dt time.Time) {
	_, err = db.Exec(updateReportDataRegisteredSql, name, dt)
	if err != nil {
		log.Fatal(err)
	}
}

func UpdateSetting(ident, val string) {
	_, err = db.Exec(updateSettingSql, val, ident)
	if err != nil {
		log.Fatal(err)
	}
}

type SyncWriter struct {
	file *os.File
}

func (sw SyncWriter) Write(p []byte) (n int, err error) {
	n, err = sw.file.Write(p)
	sw.file.Sync()
	return n, err
}

func initLog() {
	LogFile, err = os.OpenFile("app.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(SyncWriter{file: LogFile})
}

func applyMigrations() {
	for _, query := range initSqls {
		_, err = db.Exec(query)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func initSettings() {
	rows, err := db.Query(getSettingsSql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var item SettingsRow
		err = rows.Scan(&item.Ident, &item.Val)
		if err != nil {
			log.Fatal(err)
		}
		SettingsRows = append(SettingsRows, item)
		switch item.Ident {
		case GitlabUrlSetting:
			Settings.GitlabUrl = item.Val
		case GitlabTokenSetting:
			Settings.GitlabToken = item.Val
		case RecentItemsSetting:
			Settings.RecentItems = strings.Split(item.Val, ",")
		case RecentCountInRowSetting:
			Settings.RecentCountInRow, err = strconv.Atoi(item.Val)
			if err != nil {
				Settings.RecentCountInRow = 3
			}
		case OtherIssueSetting:
			Settings.OtherIssue = item.Val
		}
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}

package service

import (
	"log"
	"time"
)

type RefreshData struct {
	CurrentJob  string
	CurrentTime int
	DayTime     int
}

type SettingsRow struct {
	Ident string
	Val   string
}

type SettingsData struct {
	GitlabUrl        string
	GitlabToken      string
	RecentItems      []string
	RecentCountInRow int
	OtherIssue       string
}

var Settings SettingsData
var SettingsRows []SettingsRow

func init() {
	log.Println("Init mainService")
}

func Init(callback func(RefreshData)) {
	refresh(callback)
	go func() {
		for {
			time.Sleep(time.Minute)
			refresh(callback)
		}
	}()
	RecalcReport()
}

func refresh(callback func(RefreshData)) {
	var data RefreshData
	n, t := getActiveJob()
	if n != "" {
		data.CurrentJob = n
		data.CurrentTime = t
	}
	data.DayTime = getMinutesToday()
	callback(data)
}

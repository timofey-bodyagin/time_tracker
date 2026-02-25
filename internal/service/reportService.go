package service

import (
	"log"
	"time"
)

type ReportItem struct {
	Date       time.Time
	Name       string
	Minutes    int
	Registered bool
}

func init() {
	log.Println("Init reportService")
}

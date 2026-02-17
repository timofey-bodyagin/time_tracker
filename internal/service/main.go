package service

import (
	"fmt"
	"log"
	"time"
	"tracker/internal/forms"
)

func Start(args ...string) {
	var s string
	if len(args) == 0 {
		s, err = forms.StartJobValue.Get()
		if err != nil {
			log.Println(err)
		}
	} else {
		s = args[0]
	}
	forms.CurrentJob.Set(s)
	forms.StartJobValue.Set("")
	forms.CurrentTime.Set("0 мин")
	currentTime := time.Now()
	SaveFinish(currentTime)
	SaveStart(currentTime, s, "")
}

func Stop() {
	forms.CurrentJob.Set("")
	forms.CurrentTime.Set("")
	SaveFinish(time.Now())
}

func StartRefresh() {
	refresh()
	go func() {
		for {
			time.Sleep(time.Minute)
			refresh()
		}
	}()
}

func refresh() {
	n, t := getActiveJob()
	if (n != "") {
		forms.CurrentJob.Set(n)
		forms.CurrentTime.Set(fmt.Sprintf("%d мин", t))
	} 
	t = getMinutesToday()
	forms.DayTime.Set(fmt.Sprintf("%d мин", t))
}
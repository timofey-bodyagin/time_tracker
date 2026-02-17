package main

import (
	"fyne.io/fyne/v2/app"

	_ "github.com/glebarez/go-sqlite"

	"tracker/internal/forms"
	"tracker/internal/service"
)

func main() {
	service.InitDb()

	a := app.New()
	w := forms.InitMainForm(a)
	forms.OnStart = service.Start
	forms.OnStop = service.Stop
	service.StartRefresh()
	
	w.ShowAndRun()
}

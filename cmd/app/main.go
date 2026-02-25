package main

import (
	"fyne.io/fyne/v2/app"

	_ "github.com/glebarez/go-sqlite"

	"tracker/internal/forms"
	"tracker/internal/service"
)

func main() {
	defer service.LogFile.Close()
	a := app.New()
	w := forms.InitMainForm(a)
	service.Init(forms.OnRefresh)

	w.ShowAndRun()
}

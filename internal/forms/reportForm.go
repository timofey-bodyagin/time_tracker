package forms

import (
	"fmt"
	"log"
	"time"

	"tracker/internal/graphql"
	"tracker/internal/service"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var ReportWindow fyne.Window
var notRegisteredIcon = theme.UploadIcon()
var registeredIcon = theme.CheckButtonCheckedIcon()

func showReportWindow(a fyne.App) {

	reportData := service.GetReportData(time.Now())
	list := container.NewVBox()
	for _, item := range reportData {
		check := widget.NewButtonWithIcon("", registeredIcon, func() {})
		check.Resize(fyne.NewSize(20, 10))
		check.Disable()
		if (!item.Registered) {
			check.SetIcon(notRegisteredIcon)
			check.Enable()
			check.OnTapped = func() {
					sendToGitlab(item.Name, item.Minutes, item.Date.Format("2006-01-02"))
					check.SetIcon(registeredIcon)
					check.Disable()
				}
		}
		list.Add(container.NewHBox(
			widget.NewLabel(item.Date.Format("2006-01-02")),
			widget.NewLabel(item.Name),
			layout.NewSpacer(),
			widget.NewLabel(fmt.Sprintf(fmtMinutes, item.Minutes)),
			check,
		))

	}
	scr := container.NewVScroll(list)

	ReportWindow = a.NewWindow("Отчет")
	ReportWindow.SetContent(container.NewPadded(scr))
	ReportWindow.Resize(fyne.NewSize(800, 800))
	ReportWindow.SetOnClosed(func() {
		ReportWindow = nil
	})
	ReportWindow.Show()
}

func sendToGitlab(issue string, minutes int, date string) {
	graphql.AddSpendTime(issue, minutes, date)
	
	dt, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Println("Error parsing date:", err)
		return
	}
	service.SetRegistered(issue, dt)
}

func init() {
	log.Println("Init reportForm")
}
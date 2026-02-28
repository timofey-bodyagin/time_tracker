package forms

import (
	"fmt"
	"log"
	"strings"
	"time"

	"tracker/internal/graphql"
	"tracker/internal/service"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var ReportWindow fyne.Window
var notRegisteredIcon = theme.UploadIcon()
var registeredIcon = theme.CheckButtonCheckedIcon()
var currentMonthLabel = binding.NewString()
var currentMonth = truncateToMonth(time.Now())
var dataContainer *fyne.Container
var reportData []service.ReportItem = []service.ReportItem{}
var searchString = binding.NewString()

func showReportWindow(a fyne.App) {
	dataContainer = container.NewVBox()
	refreshReport()

	scr := container.NewVScroll(dataContainer)

	ReportWindow = a.NewWindow("Отчет")
	ReportWindow.SetContent(
		container.NewBorder(
			initReportFormToolbar(),
			nil,
			nil,
			nil,
			container.NewPadded(scr),
		),
	)
	ReportWindow.Resize(fyne.NewSize(800, 800))
	ReportWindow.SetOnClosed(func() {
		ReportWindow = nil
	})
	ReportWindow.Show()
}

func sendToGitlab(issue string, minutes int, date string) {
	err := graphql.AddSpendTime(issue, minutes, date)
	if err != nil {
		dialog.ShowError(err, ReportWindow)
		return
	}
	
	dt, err := time.Parse("2006-01-02", date)
	if err != nil {
		dialog.ShowError(err, ReportWindow)
		return
	}
	service.SetRegistered(issue, dt)
}

func init() {
	log.Println("Init reportForm")
}

func initReportFormToolbar() *fyne.Container {
	recalcButton := widget.NewButtonWithIcon("Пересчитать", theme.ViewRefreshIcon(), func() {
		service.RecalcReport()
		refreshReport()
	})
	prevButton := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		currentMonth = currentMonth.AddDate(0, -1, 0)
		refreshMonthLabel()
		refreshReport()
	})
	nextButton := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		currentMonth = currentMonth.AddDate(0, 1, 0)
		refreshMonthLabel()
		refreshReport()
	})
	label := widget.NewLabelWithData(currentMonthLabel)
	refreshMonthLabel()	
	searchEntry := NewFixedWidthEntry(searchString, 250)
	searchEntry.OnChanged = func(s string) {
		filter()
		if s == "" {
			searchEntry.ActionItem.Hide()
		} else {
			searchEntry.ActionItem.Show()
		}
	}
	clearBtn := widget.NewButtonWithIcon("", theme.ContentClearIcon(), func() {
		searchString.Set("")
	})
	clearBtn.Importance = widget.LowImportance
	searchEntry.ActionItem = clearBtn
	searchEntry.ActionItem.Hide()
	searchEntry.SetIcon(theme.SearchIcon())
	searchEntry.SetPlaceHolder("Поиск")
	searchEntry.Refresh()
	return container.NewHBox(
		recalcButton,
		layout.NewSpacer(),
		prevButton,
		label,
		nextButton,
		layout.NewSpacer(),
		searchEntry,
	)
}
func refreshReport() {
	dataContainer.RemoveAll()

	reportData = service.GetReportData(currentMonth)
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
		row := container.NewHBox(
			widget.NewLabel(formatDate(item.Date)),
			widget.NewLabel(item.Name),
			layout.NewSpacer(),
			widget.NewLabel(fmt.Sprintf(fmtMinutes, item.Minutes)),
			check,
		)
		if isFiltered(item) {
			row.Hide()
		}
		dataContainer.Add(row)
	}
	dataContainer.Refresh()
}

func refreshMonthLabel() {
	currentMonthLabel.Set(currentMonth.Format("January 2006"))
}

func isFiltered(item service.ReportItem) bool {
	str, _ := searchString.Get()
	return !(strings.Contains(item.Name, str) || strings.Contains(formatDate(item.Date), str))
}

func filter() {
	for i, item := range reportData {
		row := dataContainer.Objects[i]
		if isFiltered(item) {
			row.Hide()
		} else {
			row.Show()
		}
	}
}

func formatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func truncateToMonth(t time.Time) time.Time {
	t = t.UTC()
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
}
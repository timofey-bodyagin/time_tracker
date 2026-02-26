package forms

import (
	// "errors"
	"fmt"
	"log"
	"slices"
	"time"
	_ "tracker/internal/graphql"
	"tracker/internal/service"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const currentJobHeader = "Активная задача:"
const currentTimeHeader = "Время:"
const dayTimeHeader = "Время за день:"
const stopButtonText = "\n Стоп \n"
const startButtonText = "Старт"
const jobPlaceHolder = "Задача"
const fmtMinutes = "%d мин"

var currentJob = binding.NewString()
var currentTime = binding.NewString()
var dayTime = binding.NewString()
var startJobValue = binding.NewString()
var recentUsedContainer *fyne.Container
var startJobEdit *widget.Entry
var errorPopup *widget.PopUp
var lastJobButton *widget.Button

func InitMainForm(a fyne.App) fyne.Window {
	w := a.NewWindow("Тайм-трекер")
	w.Resize(fyne.NewSize(650, 100))
	w.SetMaster()

	initRecentUsedContainer()

	infoContainer := container.NewAdaptiveGrid(2,
		widget.NewLabel(currentJobHeader),
		widget.NewLabelWithData(currentJob),
		widget.NewLabel(currentTimeHeader),
		widget.NewLabelWithData(currentTime),
		widget.NewLabel(dayTimeHeader),
		widget.NewLabelWithData(dayTime),
	)

	errorPopup = widget.NewPopUp(widget.NewLabel("Необходимо заполнить поле"), w.Canvas())

	startJobEdit = widget.NewEntryWithData(startJobValue)
	startJobEdit.SetPlaceHolder(jobPlaceHolder)
	startJobEdit.OnSubmitted = func(s string) { startCustom() }

	startButton := widget.NewButton(startButtonText, startCustom)
	startContainer := container.NewAdaptiveGrid(2, startJobEdit, startButton)

	stopButton := widget.NewButton(stopButtonText, stop)
	stopContainer := container.NewAdaptiveGrid(1, stopButton)

	mainContainer := container.NewPadded(
		container.NewVBox(
			infoContainer,
			widget.NewSeparator(),
			startContainer,
			widget.NewSeparator(),
			recentUsedContainer,
			widget.NewSeparator(),
			stopContainer,
			widget.NewSeparator(),
		))
	w.SetContent(container.NewBorder(
		initToolbar(a, w),
		nil,
		nil,
		nil,
		mainContainer,
	))

	return w
}

func OnRefresh (data service.RefreshData) {
	currentJob.Set(data.CurrentJob)
	if data.CurrentJob == "" {
		currentTime.Set("")
	} else {
		currentTime.Set(fmt.Sprintf(fmtMinutes, data.CurrentTime))
	}
	dayTime.Set(fmt.Sprintf(fmtMinutes, data.DayTime))
}

func initToolbar(a fyne.App, w fyne.Window) *widget.Toolbar {
	reportAction := widget.NewToolbarAction(theme.DocumentIcon(), func() {
		if ReportWindow == nil {
			showReportWindow(a)
		} else {
			ReportWindow.RequestFocus()
		}
		
	})
	settingAction := widget.NewToolbarAction(theme.SettingsIcon(), func() {
		showSettingsForm(w)
	})
	return widget.NewToolbar(
		reportAction,
		widget.NewToolbarSpacer(),
		settingAction,
	)
}

func startCustom() {
	txt := startJobEdit.Text
	if txt == "" {
		errorPopup.ShowAtRelativePosition(fyne.NewPos(10, -errorPopup.MinSize().Height-10), startJobEdit)
	} else {
		start(txt)
	}
}

func start(name string) {
	updateLastJobButton()
	currentJob.Set(name)
	currentTime.Set(fmt.Sprintf(fmtMinutes, 0))
	startJobValue.Set("")
	currTime := time.Now()
	service.SaveFinish(currTime)
	service.SaveStart(currTime, name, "")
}

func stop() {
	updateLastJobButton()
	currentJob.Set("")
	currentTime.Set("")
	service.SaveFinish(time.Now())
}

func initRecentUsedContainer() {
	recentUsedContainer = container.NewAdaptiveGrid(service.Settings.RecentCountInRow)
	for _, val := range service.Settings.RecentItems {
		txt := fmt.Sprintf("\n%s\n", val)
		recentUsedContainer.Add(widget.NewButton(txt, func() {
			start(val)
		}))
	}
}

func updateLastJobButton() {
	txt, _ := currentJob.Get()
	if txt != "" && !slices.Contains(service.Settings.RecentItems, txt) {
		if lastJobButton == nil {
			lastJobButton = widget.NewButton(txt, func() {
				start(txt)
			})
			recentUsedContainer.Add(lastJobButton)
		} else {
			lastJobButton.Text = txt
			lastJobButton.OnTapped = func() {
				start(txt)
			}
			lastJobButton.Refresh()
		}
	}
}

func init() {
	log.Println("Init mainForm")
}


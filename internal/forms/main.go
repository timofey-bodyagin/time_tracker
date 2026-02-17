package forms

import (
	// "errors"
	"fmt"
	"slices"
	"time"

	// "tracker/internal/service"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

const currentJobHeader = "Активная задача:"
const currentTimeHeader = "Время:"
const dayTimeHeader = "Время за день:"
const stopButtonText = "\n Стоп \n"
const startButtonText = "Старт"
const jobPlaceHolder = "Задача"

var recentUsedNames = []string{
	"Дейлик",
	"Статус",
	"Демо",
	"Ретро",
	"Предпланирование",
}
var CurrentJob = binding.NewString()
var CurrentTime = binding.NewString()
var DayTime = binding.NewString()
var StartJobValue = binding.NewString()
var OnStart func(args ...string)
var OnStop func()
var recentUsedContainer *fyne.Container
var startJobEdit *widget.Entry
var errorPopup *widget.PopUp
var lastJobButton *widget.Button

func InitMainForm(a fyne.App) fyne.Window {
	w := a.NewWindow("Тайм-трекер")
	w.Resize(fyne.NewSize(400, 100))
	initRecentUsedContainer()

	infoContainer := container.NewAdaptiveGrid(2,
		widget.NewLabel(currentJobHeader),
		widget.NewLabelWithData(CurrentJob),
		widget.NewLabel(currentTimeHeader),
		widget.NewLabelWithData(CurrentTime),
		widget.NewLabel(dayTimeHeader),
		widget.NewLabelWithData(DayTime),
	)

	startJobEdit = widget.NewEntryWithData(StartJobValue)
	startJobEdit.SetPlaceHolder(jobPlaceHolder)

	errorPopup = widget.NewPopUp(widget.NewLabel("Необходимо заполнить поле"), w.Canvas())

	startJobEdit.OnSubmitted = func(s string) { start() }

	startButton := widget.NewButton(startButtonText, start)
	startContainer := container.NewAdaptiveGrid(2, startJobEdit, startButton)
	startJobEdit.Validate()

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
	w.SetContent(mainContainer)

	return w
}

func start() {
	txt := startJobEdit.Text
	if txt == "" {
		errorPopup.ShowAtRelativePosition(fyne.NewPos(10, -errorPopup.MinSize().Height-10), startJobEdit)
	} else {
		currTxt, _ := CurrentJob.Get()
		updateLastJobButton(currTxt)
		OnStart()
	}
}

func startRecent(txt string) {
	currentTxt, _ := CurrentJob.Get()
	updateLastJobButton(currentTxt)
	OnStart(txt)
}

func stop() {
	txt, _ := CurrentJob.Get()
	updateLastJobButton(txt)
	OnStop()
}

func initRecentUsedContainer() {
	recentUsedContainer = container.NewAdaptiveGrid(3)
	for _, val := range recentUsedNames {
		txt := fmt.Sprintf("\n%s\n", val)
		recentUsedContainer.Add(widget.NewButton(txt, func() {
			startRecent(val)
		}))
	}
}

func updateLastJobButton(txt string) {
	if txt != "" && !slices.Contains(recentUsedNames, txt) {
		if lastJobButton == nil {
			lastJobButton = widget.NewButton(txt, func() {
				startRecent(txt)
			})
			recentUsedContainer.Add(lastJobButton)
		} else {
			lastJobButton.Text = txt
			lastJobButton.OnTapped = func() {
				startRecent(txt)
			}
			lastJobButton.Refresh()
		}
	}
}

func nvl(t *time.Time) time.Time {
	if t == nil {
		return time.Now()
	}
	return *t
}

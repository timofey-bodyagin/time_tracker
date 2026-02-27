package forms

import (
	// "errors"
	"fmt"
	"log"
	"slices"
	"sort"
	"strings"
	"time"
	"tracker/internal/graphql"
	_ "tracker/internal/graphql"
	"tracker/internal/service"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	currentJobHeader  = "Активная задача:"
	currentTimeHeader = "Время:"
	dayTimeHeader     = "Время за день:"
	stopButtonText    = "\n Стоп \n"
	startButtonText   = "Старт"
	jobPlaceHolder    = "Задача"
	fmtMinutes        = "%d мин"
	maxCountLastItems = 5
)

var (
	currentJob          = graphql.IssueInfo{}
	currentJobLabel     = widget.NewLabel("")
	currentTimeLabel    = widget.NewLabel("")
	dayTimeLabel        = widget.NewLabel("")
	startJobValue       = binding.NewString()
	recentUsedContainer *fyne.Container
	startJobEdit        *widget.Entry
	errorPopup          *widget.PopUp
	lastItemsArray      = []LastItem{}
	lastItemsContainer  = container.NewVBox()
)

func InitMainForm(a fyne.App) fyne.Window {
	w := a.NewWindow("Тайм-трекер")
	w.Resize(fyne.NewSize(650, 100))
	w.SetMaster()

	currentJobLabel.Wrapping = fyne.TextWrapWord

	initRecentUsedContainer()

	infoContainer := &widget.Form{}
	infoContainer.Append(currentJobHeader, currentJobLabel)
	infoContainer.Append(currentTimeHeader, currentTimeLabel)
	infoContainer.Append(dayTimeHeader, dayTimeLabel)

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
			lastItemsContainer,
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

func OnRefresh(data service.RefreshData) {
	if currentJob.Iid != data.CurrentJob {
		if isNumeric(data.CurrentJob) {
			currentJob = graphql.GetIssueInfo(data.CurrentJob)
		}
		currentJob.Iid = data.CurrentJob
		currentJobLabel.SetText(currentJob.Str())
	}
	if currentJob.Iid == "" {
		currentTimeLabel.SetText("")
	} else {
		currentTimeLabel.SetText(fmt.Sprintf(fmtMinutes, data.CurrentTime))
	}
	dayTimeLabel.SetText(fmt.Sprintf(fmtMinutes, data.DayTime))
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
	currentJob = graphql.GetIssueInfo(name)
	currentJob.Iid = strings.TrimSpace(name)
	currentJobLabel.SetText(currentJob.Str())
	currentTimeLabel.SetText(fmt.Sprintf(fmtMinutes, 0))
	startJobValue.Set("")
	currTime := time.Now()
	service.SaveFinish(currTime)
	service.SaveStart(currTime, name, "")
	lastItemsArray = slices.DeleteFunc(lastItemsArray, func(item LastItem) bool {
		return item.info.Iid == name
	})
	refreshLastItemsContainer()
}

func stop() {
	updateLastJobButton()
	currentJob = graphql.IssueInfo{}
	currentJobLabel.SetText("")
	currentTimeLabel.SetText("")
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
	if currentJob.Iid == "" || slices.Contains(service.Settings.RecentItems, currentJob.Iid) {
		return
	}

	contains := false
	oldestItemIndex := 0
	t := time.Now()

	for i, item := range lastItemsArray {
		if item.info.Iid == currentJob.Iid {
			contains = true
			lastItemsArray[i].t = time.Now()
		}
		if item.t.Before(t) {
			oldestItemIndex = i
			t = item.t
		}
	}

	if !contains {
		id := currentJob.Iid
		label := widget.NewLabel(currentJob.Str())
		label.Wrapping = fyne.TextWrapWord
		label.Alignment = fyne.TextAlignCenter

		button := widget.NewButton("", func() {
			start(id)
		})
		content := container.NewStack(button, label)

		lastItem := LastItem{
			t:      time.Now(),
			info:   currentJob,
			button: content,
		}

		if len(lastItemsArray) == maxCountLastItems {
			lastItemsArray[oldestItemIndex] = lastItem
		} else {
			lastItemsArray = append(lastItemsArray, lastItem)
		}
	}
	refreshLastItemsContainer()
}

func refreshLastItemsContainer() {
	sort.Slice(lastItemsArray, func(i, j int) bool {
		return lastItemsArray[i].t.After(lastItemsArray[j].t)
	})
	lastItemsContainer.RemoveAll()
	for _, item := range lastItemsArray {
		lastItemsContainer.Add(item.button)
	}
}

func init() {
	log.Println("Init mainForm")
}

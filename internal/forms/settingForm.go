package forms

import (
	"tracker/internal/service"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

const (
	title = "Настройки"
	okButtonText = "Сохранить"
	cancelButtonText = "Отменить"
)

func showSettingsForm(w fyne.Window) {
	itemContainer := container.NewAdaptiveGrid(2)
	var values = []binding.String{}
	for i, item := range service.SettingsRows {
		values = append(values, binding.NewString())
		values[i].Set(item.Val)
		itemContainer.Add(widget.NewLabel(item.Ident))
		entry := widget.NewEntryWithData(values[i])
		itemContainer.Add(entry)
	}
	mainContainer := container.NewPadded(itemContainer)
	dialog.ShowCustomConfirm(title, okButtonText, cancelButtonText, mainContainer, func(b bool) {onClose(b, values)}, w)
}

func onClose(isOk bool, values []binding.String) {
	if isOk {
		for i, item := range service.SettingsRows {
			val, _ := values[i].Get()
			service.SettingsRows[i].Val = val
			service.UpdateSetting(item.Ident, val)
		}
	}
}
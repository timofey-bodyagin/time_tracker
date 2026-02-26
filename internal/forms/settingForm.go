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
	form := &widget.Form{}
	var values = []binding.String{}
	for i, item := range service.SettingsRows {
		values = append(values, binding.NewString())
		values[i].Set(item.Val)
		form.Append(item.Ident, NewFixedWidthEntry(values[i], 400))
	}
	
	mainContainer := container.NewPadded(form)
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
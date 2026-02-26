package forms

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type FixedWidthEntry struct {
	widget.Entry
	width float32
}

func (f *FixedWidthEntry) MinSize() fyne.Size {
	return fyne.NewSize(f.width, f.Entry.MinSize().Height)
}

func NewFixedWidthEntry(data binding.String, width float32) *FixedWidthEntry {
	e := &FixedWidthEntry{width: width}
	e.Bind(data)
	e.Validator = nil
	e.ExtendBaseWidget(e)
	return e
}
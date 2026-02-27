package forms

import (
	"regexp"
	"time"
	"tracker/internal/graphql"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type LastItem struct {
	t time.Time
	info graphql.IssueInfo
	button *fyne.Container
}

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

func isNumeric(s string) bool {
	var re = regexp.MustCompile(`^[0-9]+$`)
	return re.MatchString(s)
}
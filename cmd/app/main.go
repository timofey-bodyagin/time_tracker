package main

import (
	"fyne.io/fyne/v2/app"

	_ "github.com/glebarez/go-sqlite"

	"tracker/internal/forms"
	"tracker/internal/service"
)

func main() {
	service.InitDb()

	a := app.New()
	w := forms.InitMainForm(a)
	forms.OnStart = service.Start
	forms.OnStop = service.Stop
	service.StartRefresh()
	
	// rows, err := db.Query("SELECT * FROM tmp")
    // if err != nil {
    //     log.Fatal(err)
    // }
    // defer rows.Close()
    // // Loop through rows, using Scan to assign column data to struct fields.
    // for rows.Next() {
    //     var id int
    //     if err := rows.Scan(&id); err != nil {
    //         log.Fatal(err)
    //     }
    //     labelsContainer.Add(widget.NewLabel(strconv.Itoa(id)))
    // }
    // if err := rows.Err(); err != nil {
    //     log.Fatal(err)
    // }

	w.ShowAndRun()
}

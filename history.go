package main

import (
	"errors"
)

// historical record
type hRecord struct {
	screen activeScreen
	id     int
}

func historyInit() {
	addMessage(message{isError: false, message: "History system intialized"})
}

func historyAdd(rec hRecord) {
	app.hist = append(app.hist, rec)
}

// go to the last projects or items screen we were on
func historyGoToLast(screens activeScreen) error {
	for i := len(app.hist) - 1; i >= 0; i-- {
		rec := app.hist[i]

		if rec.screen&screens != 0 {
			return historyGoTo(rec)
		}
	}

	return errors.New("Screen not found")
}

func historyGoTo(rec hRecord) error {
	switch rec.screen {
	case items:
		// set the project via its ID
		if rec.id == -1 {
			app.ui.project = getToday()
		} else if rec.id == -2 {
			app.ui.project = getTomorrow()
		} else {
			for key, value := range app.todoist.Projects {
				if value.ID == rec.id {
					app.ui.project = &app.todoist.Projects[key]
					break
				}
			}
		}

		showScreen(rec.screen)
		return nil
	case projects:
		showScreen(rec.screen)
		return nil
	}

	return errors.New("Unable to go to history")
}

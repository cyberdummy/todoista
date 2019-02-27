package main

import (
	"github.com/gdamore/tcell"
)

func bindsInit() {
	app.ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			if app.ui.screen == addItem || app.ui.screen == updateItem {
				break
			}

			switch event.Rune() {
			case 's':
				DoSync()
				break
			case 'a':
				showScreen(addItem)
				break
			case 'u':
				if app.ui.screen == items {
					showScreen(updateItem)
					break
				}
			case 'p':
				showScreen(projects)
				break
			case 'm':
				showScreen(messages)
				break
			case 'q':
				app.ui.app.Stop()
				break
			}
		}

		return event
	})
}

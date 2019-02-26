// The todoista program a CLI UI for working with todoist.
package main

import (
	"github.com/cyberdummy/todoista/todoist"
)

type todoista struct {
	todoist *todoist.Todoist
	msgs []message
	cfg config
	ui userInterface
}

var app todoista


func main() {
	var err error
	// Read in config from various sources
	getConfig()

	app.todoist,err = todoist.New(app.cfg.apiKey)

	if err != nil {
		panic("Unable to create new todoist")
	}

	// Creates the UI app, sets up the binds
	uiInit()

	messagesInit()

	showScreen(projects)
	DoSync()

	uiRun()

	messagesShutdown()
}

func SetUiMessage(message string) {
	app.ui.msg.SetText(message)
	app.ui.app.Draw()
}

func DoSync() {
	var err error
	SetUiMessage("Syncing...")
	app.todoist, err = app.todoist.ReadSync()

	if err != nil {
		SetUiMessage("Sync failed! [red]"+err.Error())
		addMessage(message{message: err.Error(), isError: true, })
		return
	}

	switch app.ui.screen {
	case projects:
		buildProjectIdx()
		break
	case items:
		buildItemIdx()
		break
	}

	SetUiMessage("Sync Complete")
	app.ui.app.Draw()
}

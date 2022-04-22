// The todoista program a CLI UI for working with todoist.
package main

import (
	"github.com/cyberdummy/todoista/todoist"
)

type todoista struct {
	todoist *todoist.Todoist
	msgs    []message
	hist    []hRecord
	cfg     config
	ui      userInterface
}

var app todoista

func main() {
	var err error
	// Read in config from various sources
	getConfig()

	app.todoist, err = todoist.New(app.cfg.apiKey)

	if err != nil {
		panic("Unable to create new todoist")
	}

	// Creates the UI app, sets up the binds
	uiInit()

	messagesInit()
	historyInit()

	showScreen(projects)
	doSync()

	uiRun()

	messagesShutdown()
}

func setUIMessage(msg string) {
	app.ui.msg.SetText(msg)
	addMessage(message{isError: false, message: msg})
	//app.ui.app.Draw()
}

func doSync() {
	var err error
	setUIMessage("Syncing...")
	app.todoist, err = app.todoist.ReadSync()

	if err != nil {
		setUIMessage("Sync failed! [red]" + err.Error())
		addMessage(message{message: err.Error(), isError: true})
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

	setUIMessage("Sync Complete")
	//app.ui.app.Draw()
}

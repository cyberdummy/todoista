package main

import (
	"github.com/cyberdummy/todoista/todoist"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type activeScreen int

const (
	projects activeScreen = 1 << iota
	items
	messages
	addItem
	updateItem
)

type userInterface struct {
	app       *tview.Application
	idxLayout *tview.Flex
	idx       *tview.Table
	status    *tview.TextView
	msg       *tview.TextView
	screen    activeScreen
	project   *todoist.Project
	items     []*todoist.Item
}

// uiInit called at startup after the config is read
func uiInit() {
	app.ui.app = tview.NewApplication()

	bindsInit()
	createStatus()
	createMessage()
	createIndexLayout()
}

// showScreen make sure the specified screen layout is displayed.
func showScreen(show activeScreen) {
	if app.ui.screen == show {
		setUIMessage("Already on that screen")
		return
	}

	switch show {
	case projects:
		showProjectsUI() // projects.go
		historyAdd(hRecord{screen: show})
		break
	case items:
		showItemsUI() // items.go
		historyAdd(hRecord{screen: show, id: app.ui.project.ID})
		break
	case messages:
		showMessagesUI() // messages.go
		break
	case addItem:
		showAddItem()
		app.ui.screen = show
		return
	case updateItem:
		// figure out selected item
		row, _ := app.ui.idx.GetSelection()
		showUpdateItem(app.ui.items[row])
		app.ui.screen = show
		return
		return
	}

	app.ui.screen = show

	app.ui.app.SetRoot(app.ui.idxLayout, true)
	app.ui.app.SetFocus(app.ui.idx)
	app.ui.app.Draw()
}

// uiRun starts the main loop thats draws the UI.
func uiRun() {
	if err := app.ui.app.Run(); err != nil {
		panic(err)
	}
}

// createStatus create the status element
func createStatus() {
	app.ui.status = tview.NewTextView().
		SetDynamicColors(true).
		SetText("This is the status bar")
}

func createMessage() {
	app.ui.msg = tview.NewTextView().
		SetDynamicColors(true).
		SetText("-").
		SetChangedFunc(func() {
			app.ui.app.Draw()
		})
}

// createIndexLayout creates the layout that contains an index table.
func createIndexLayout() {
	app.ui.idx = tview.NewTable()
	app.ui.idx.SetBackgroundColor(tcell.ColorGray)
	app.ui.idx.SetSelectedStyle(tcell.ColorBlack, tcell.ColorOlive, 0)

	app.ui.idxLayout = tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(app.ui.status, 1, 1, false).
			AddItem(app.ui.idx, 0, 3, false).
			AddItem(app.ui.msg, 1, 1, false), 0, 2, false)
}

func createFormLayout(form *tview.Form) {
	form.SetBackgroundColor(tcell.ColorGray)

	layout := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(app.ui.status, 1, 1, false).
			AddItem(form, 0, 3, false).
			AddItem(app.ui.msg, 1, 1, false), 0, 2, false)

	app.ui.app.SetRoot(layout, true)
	app.ui.app.SetFocus(form)
	app.ui.app.Draw()
}

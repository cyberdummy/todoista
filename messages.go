package main

import (
	"log"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Defines a message
type message struct {
	isError bool
	message string
	time    time.Time
}

// messagesInit called when the app starts up, to perform any intialization
// required.
func messagesInit() {
	addMessage(message{isError: false, message: "Message subsystem intialized"})
}

// messagesShutdown is called then the app is closed, prints any error messages
// in the stack to stderr.
func messagesShutdown() {
	for _, value := range app.msgs {
		if value.isError {
			log.Println(value.message)
		}
	}
}

// addMessage adds a message to the stack.
func addMessage(message message) {
	message.time = time.Now()
	app.msgs = append(app.msgs, message)
}

// showMessagesUI called when we want to display the interface for view the
// message stack.
func showMessagesUI() {
	app.ui.idx.Clear()
	app.ui.idx.SetSelectable(false, false)

	for key, value := range app.msgs {
		cell := tview.NewTableCell("[" + value.time.Format("2006/01/02 15:04:05") + "]")
		cell.SetAlign(tview.AlignLeft)
		cell.SetBackgroundColor(tcell.ColorGray)
		cell.SetTextColor(tcell.ColorDefault)
		app.ui.idx.SetCell(key, 0, cell)

		app.ui.idx.SetCell(key, 0, cell)
		cell = tview.NewTableCell(value.message)
		cell.SetAlign(tview.AlignLeft)
		cell.SetBackgroundColor(tcell.ColorGray)
		if value.isError {
			cell.SetTextColor(tcell.ColorRed)
		} else {
			cell.SetTextColor(tcell.ColorDefault)
		}

		cell.SetExpansion(1)
		app.ui.idx.SetCell(key, 1, cell)
	}

	app.ui.idx.SetSelectable(true, false)
	app.ui.idx.Select(0, 0)
}

package main

import (
	"github.com/cyberdummy/todoista/todoist"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func showItemsUi() {
	app.ui.status.SetText(app.ui.project.Name)
	app.ui.idx.SetSelectable(false, false)
	// When a user selects a project
	app.ui.idx.SetSelectedFunc(func(row int, column int) {
		SetUiMessage("Completing Task")
		app.todoist.ItemComplete(*app.ui.items[row])
		SetUiMessage("Task Completed")
		DoSync()
	})

	// Build the project table rows
	buildItemIdx()

	app.ui.idx.SetSelectable(true, false)
	app.ui.idx.Select(0,0)
}

func buildItemIdx() {
	app.ui.idx.Clear()

	row := 0

	var items []todoist.Item
	app.ui.items = nil

	if app.ui.project.GetItems != nil {
		items = app.ui.project.GetItems()
	} else {
		items = app.todoist.Items
	}

	for key,value := range items {
		if app.ui.project.GetItems == nil && value.ProjectId != app.ui.project.ID {
			continue
		}

		cell := tview.NewTableCell(value.Content)
		cell.SetAlign(tview.AlignLeft)
		cell.SetBackgroundColor(tcell.ColorGray)
		cell.SetTextColor(tcell.ColorDefault)
		cell.SetExpansion(1)
		app.ui.idx.SetCell(row, 0, cell)

		row++
		app.ui.items = append(app.ui.items, &items[key])
	}
}

// addItem displays the form for adding an item.
func showAddItem() {
	var form *tview.Form

	// make the project drop down
	dd := make([]string, len(app.todoist.Projects))

	for key,value := range app.todoist.Projects {
		dd[key] = value.Name
	}

	form = tview.NewForm().
	AddInputField("Task", "", 0, nil, nil).
	AddInputField("Date", "tomorrow", 0, nil, nil).
	AddDropDown("Project", dd, 0, nil).
	AddButton("Save", func() {
		SetUiMessage("Saving item...")

		idx,_ := form.GetFormItem(2).(*tview.DropDown).GetCurrentOption()

		err := app.todoist.ItemAdd(
			form.GetFormItem(0).(*tview.InputField).GetText(),
			form.GetFormItem(1).(*tview.InputField).GetText(),
			app.todoist.Projects[idx].ID)

		if err != nil {
			SetUiMessage("Add item failed! [red]"+err.Error())
			addMessage(message{message: err.Error(), isError: true, })
			return
		}

		showScreen(projects)
		DoSync()
	}).
	AddButton("Quit", func() {
		showScreen(projects)
	})

	form.SetBorder(true).SetTitle("Add Item").SetTitleAlign(tview.AlignLeft)

	createFormLayout(form)
}

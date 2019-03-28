package main

import (
	"strconv"
	"strings"

	"github.com/cyberdummy/todoista/todoist"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func showItemsUI() {
	app.ui.status.SetText("|" + app.ui.project.Name + "|")
	app.ui.idx.SetSelectable(false, false)
	// When a user selects a project
	app.ui.idx.SetSelectedFunc(func(row int, column int) {
		if len(app.ui.items) < (row + 1) {
			setUIMessage("No item selected")
			return
		}

		setUIMessage("Completing Task")
		app.todoist.ItemComplete(*app.ui.items[row])
		setUIMessage("Task Completed")
		doSync()

		// if was last item in list select prev
		if len(app.ui.items) < (row+1) && row != 0 {
			app.ui.idx.Select((row - 1), 0)
		}
	})

	// Build the project table rows
	buildItemIdx()

	app.ui.idx.SetSelectable(true, false)
	app.ui.idx.Select(0, 0)
}

func buildItemIdx() {
	var cell *tview.TableCell
	var sb strings.Builder

	app.ui.idx.Clear()

	row := 0

	var items []todoist.Item
	app.ui.items = nil

	if app.ui.project.GetItems != nil {
		items = app.ui.project.GetItems()
	} else {
		items = app.todoist.Items
	}

	for key, value := range items {
		if app.ui.project.GetItems == nil && value.ProjectID != app.ui.project.ID {
			continue
		}

		if !value.DueDate.IsZero() &&
			value.DueDate.Hour() != 23 &&
			app.ui.project.ID < 0 {
			sb.WriteString(value.DueDate.Format("15:00 "))
		}

		sb.WriteString(value.Content)

		cell = tview.NewTableCell(sb.String())
		cell.SetAlign(tview.AlignLeft)
		cell.SetBackgroundColor(tcell.ColorGray)
		cell.SetTextColor(tcell.ColorDefault)
		cell.SetExpansion(1)
		app.ui.idx.SetCell(row, 0, cell)

		row++
		app.ui.items = append(app.ui.items, &items[key])
		sb.Reset()
	}
}

func itemDelete() {
	setUIMessage("Deleting Item..")
	// find selected one..
	row, _ := app.ui.idx.GetSelection()

	if len(app.ui.items) < (row + 1) {
		setUIMessage("No item selected")
		return
	}

	app.todoist.ItemDelete(app.ui.items[row])
	setUIMessage("Item Deleted")
	doSync()

	// if was last item in list select prev
	if len(app.ui.items) < (row+1) && row != 0 {
		app.ui.idx.Select((row - 1), 0)
	}
}

func itemForm(save func()) *tview.Form {
	var form *tview.Form

	dd := make([]string, len(app.todoist.Projects))

	for key, value := range app.todoist.Projects {
		dd[key] = value.Name
	}

	form = tview.NewForm().
		AddInputField("Task", "", 0, nil, nil).
		AddInputField("Date", "tomorrow", 0, nil, nil).
		AddDropDown("Project", dd, 0, nil).
		AddButton("Save", save).
		AddButton("Quit", func() {
			showScreen(projects)
		})

	return form
}

// addItem displays the form for adding an item.
func showAddItem() {
	var form *tview.Form

	app.ui.status.SetText("Add Item")

	// make the project drop down
	form = itemForm(func() {
		setUIMessage("Adding item...")

		idx, _ := form.GetFormItem(2).(*tview.DropDown).GetCurrentOption()

		err := app.todoist.ItemAdd(
			form.GetFormItem(0).(*tview.InputField).GetText(),
			form.GetFormItem(1).(*tview.InputField).GetText(),
			app.todoist.Projects[idx].ID)

		if err != nil {
			setUIMessage("Add item failed! [red]" + err.Error())
			addMessage(message{message: err.Error(), isError: true})
			return
		}

		err = historyGoToLast((items | projects))

		if err != nil {
			showScreen(projects)
		}

		doSync()
	})

	createFormLayout(form)
}

func showUpdateItem(item *todoist.Item) {
	var form *tview.Form

	app.ui.status.SetText("Edit Item")

	// make the project drop down
	form = itemForm(func() {
		setUIMessage("Updating item...")

		idx, _ := form.GetFormItem(2).(*tview.DropDown).GetCurrentOption()

		err := app.todoist.ItemUpdate(
			item,
			form.GetFormItem(0).(*tview.InputField).GetText(),
			form.GetFormItem(1).(*tview.InputField).GetText(),
			app.todoist.Projects[idx].ID)

		if err != nil {
			setUIMessage("Update item failed! [red]" + err.Error())
			addMessage(message{message: err.Error(), isError: true})
			return
		}

		addMessage(message{message: "Updated item" + strconv.Itoa(item.ID)})
		err = historyGoToLast(items)

		if err != nil {
			showScreen(projects)
		}

		doSync()
	})

	form.GetFormItem(0).(*tview.InputField).SetText(item.Content)
	form.GetFormItem(1).(*tview.InputField).SetText(item.DateString)
	// Select Project
	for key, value := range app.todoist.Projects {
		if value.ID == item.ProjectID {
			form.GetFormItem(2).(*tview.DropDown).SetCurrentOption(key)
			break
		}
	}

	createFormLayout(form)
}

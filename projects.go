package main

import (
	"time"

	"github.com/cyberdummy/todoista/todoist"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func getToday() *todoist.Project {
	return &todoist.Project{
		Name: "Today",
		ID:   -1,
		GetItems: func() []todoist.Item {
			var arr []todoist.Item

			now := time.Now()
			endDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
			for _, item := range app.todoist.Items {
				if item.DueDate.IsZero() {
					continue
				}

				if item.DueDate.Before(endDay) || item.DueDate.Equal(endDay) {
					arr = append(arr, item)
				}
			}

			return arr
		},
	}
}

func getTomorrow() *todoist.Project {
	return &todoist.Project{
		Name: "Tomorrow",
		ID:   -2,
		GetItems: func() []todoist.Item {
			var arr []todoist.Item

			tomorrow := time.Now().AddDate(0, 0, 1)

			for _, item := range app.todoist.Items {
				if item.DueDate.IsZero() {
					continue
				}

				if item.DueDate.Year() == tomorrow.Year() && item.DueDate.YearDay() == tomorrow.YearDay() {
					arr = append(arr, item)
				}
			}

			return arr
		},
	}
}

func showProjectsUI() {
	app.ui.status.SetText("Project Selection")

	app.ui.idx.SetSelectable(false, false)
	// When a user selects a project
	app.ui.idx.SetSelectedFunc(func(row int, column int) {
		if row == 0 {
			app.ui.project = getToday()
		} else if row == 1 {
			app.ui.project = getTomorrow()
		} else {
			app.ui.project = &app.todoist.Projects[(row - 2)]
		}

		showScreen(items)
	})

	// Build the project table rows
	buildProjectIdx()

	app.ui.idx.SetSelectable(true, false)
	app.ui.idx.Select(0, 0)
}

func buildProjectIdx() {
	app.ui.idx.Clear()

	cell := tview.NewTableCell("Today")
	cell.SetAlign(tview.AlignLeft)
	cell.SetBackgroundColor(tcell.ColorGray)
	cell.SetTextColor(tcell.ColorDefault)
	cell.SetExpansion(1)
	app.ui.idx.SetCell(0, 0, cell)

	cell = tview.NewTableCell("Tomorrow")
	cell.SetAlign(tview.AlignLeft)
	cell.SetBackgroundColor(tcell.ColorGray)
	cell.SetTextColor(tcell.ColorDefault)
	cell.SetExpansion(1)
	app.ui.idx.SetCell(1, 0, cell)

	for key, value := range app.todoist.Projects {
		cell := tview.NewTableCell(value.Name)
		cell.SetAlign(tview.AlignLeft)
		cell.SetBackgroundColor(tcell.ColorGray)
		cell.SetTextColor(tcell.ColorDefault)
		cell.SetExpansion(1)
		app.ui.idx.SetCell((key + 2), 0, cell)
	}
}

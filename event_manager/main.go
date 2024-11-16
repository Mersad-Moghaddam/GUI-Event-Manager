package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Event struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

var events = []Event{
	{ID: "1", Title: "Event 1", Description: "This is event 1", Date: time.Now()},
	{ID: "2", Title: "Event 2", Description: "This is event 2", Date: time.Now()},
	{ID: "3", Title: "Event 3", Description: "This is event 3", Date: time.Now()},
}

func setWindowSize(window fyne.Window, width, height int) {
	window.Resize(fyne.NewSize(float32(width), float32(height)))
}
func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Event Manager")
	setWindowSize(myWindow, 800, 600)

	viewAllEventsWindow := myApp.NewWindow("All Events")
	setWindowSize(viewAllEventsWindow, 800, 600)

	viewEventByIDWindow := myApp.NewWindow("View Event by ID")
	setWindowSize(viewEventByIDWindow, 800, 600)

	addEventWindow := myApp.NewWindow("Add Event")
	setWindowSize(addEventWindow, 800, 600)

	updateEventWindow := myApp.NewWindow("Update Event")
	setWindowSize(updateEventWindow, 800, 600)

	deleteEventWindow := myApp.NewWindow("Delete Event")
	setWindowSize(deleteEventWindow, 800, 600)

	eventList := widget.NewList(
		func() int { return len(events) },
		func() fyne.CanvasObject { return widget.NewLabel("Template") },
		func(i int, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(events[i].Title)
		},
	)

	viewAllEventsButton := widget.NewButton("View All Events", func() {
		viewAllEventsWindow := myApp.NewWindow("All Events")
		eventList := widget.NewList(
			func() int { return len(events) },
			func() fyne.CanvasObject { return widget.NewLabel("Template") },
			func(i int, item fyne.CanvasObject) {
				item.(*widget.Label).SetText(events[i].Title + " - " + events[i].Description)
			},
		)
		viewAllEventsWindow.SetContent(eventList)
		viewAllEventsWindow.Show()
	})

	viewEventByIDButton := widget.NewButton("View Event by ID", func() {
		viewEventByIDWindow := myApp.NewWindow("View Event by ID")
		idEntry := widget.NewEntry()
		viewEventButton := widget.NewButton("View Event", func() {
			for _, event := range events {
				if event.ID == idEntry.Text {
					viewEventWindow := myApp.NewWindow("Event")
					eventLabel := widget.NewLabel("Title: " + event.Title + "\nDescription: " + event.Description + "\nDate: " + event.Date.String())
					viewEventWindow.SetContent(eventLabel)
					viewEventWindow.Show()
					return
				}
			}
			viewEventWindow := myApp.NewWindow("Event")
			eventLabel := widget.NewLabel("Event not found")
			viewEventWindow.SetContent(eventLabel)
			viewEventWindow.Show()
		})
		viewEventByIDForm := container.NewVBox(
			widget.NewLabel("ID:"),
			idEntry,
			viewEventButton,
		)
		viewEventByIDWindow.SetContent(viewEventByIDForm)
		viewEventByIDWindow.Show()
	})

	addEventButton := widget.NewButton("Add Event", func() {
		addEventWindow := myApp.NewWindow("Add Event")
		titleEntry := widget.NewEntry()
		descriptionEntry := widget.NewEntry()
		dateEntry := widget.NewEntry()
		addEventButton := widget.NewButton("Add Event", func() {
			event := Event{
				ID:          "4",
				Title:       titleEntry.Text,
				Description: descriptionEntry.Text,
				Date:        time.Now(),
			}
			events = append(events, event)
			eventList.Refresh()
			addEventWindow.Close()
		})
		addEventForm := container.NewVBox(
			widget.NewLabel("Title:"),
			titleEntry,
			widget.NewLabel("Description:"),
			descriptionEntry,
			widget.NewLabel("Date:"),
			dateEntry,
			addEventButton,
		)
		addEventWindow.SetContent(addEventForm)
		addEventWindow.Show()
	})

	updateEventButton := widget.NewButton("Update Event", func() {
		updateEventWindow := myApp.NewWindow("Update Event")
		idEntry := widget.NewEntry()
		titleEntry := widget.NewEntry()
		descriptionEntry := widget.NewEntry()
		updateEventButton := widget.NewButton("Update Event", func() {
			for i, event := range events {
				if event.ID == idEntry.Text {
					events[i].Title = titleEntry.Text
					events[i].Description = descriptionEntry.Text
					eventList.Refresh()
					updateEventWindow.Close()
					return
				}
			}
			updateEventWindow := myApp.NewWindow("Event")
			eventLabel := widget.NewLabel("Event not found")
			updateEventWindow.SetContent(eventLabel)
			updateEventWindow.Show()
		})
		updateEventForm := container.NewVBox(
			widget.NewLabel("ID:"),
			idEntry,
			widget.NewLabel("Title:"),
			titleEntry,
			widget.NewLabel("Description:"),
			descriptionEntry,
			updateEventButton,
		)
		updateEventWindow.SetContent(updateEventForm)
		updateEventWindow.Show()
	})

	deleteEventButton := widget.NewButton("Delete Event", func() {
		deleteEventWindow := myApp.NewWindow("Delete Event")
		idEntry := widget.NewEntry()
		deleteEventButton := widget.NewButton("Delete Event", func() {
			for i, event := range events {
				if event.ID == idEntry.Text {
					events = append(events[:i], events[i+1:]...)
					eventList.Refresh()
					deleteEventWindow.Close()
					return
				}
			}
			deleteEventWindow := myApp.NewWindow("Event")
			eventLabel := widget.NewLabel("Event not found")
			deleteEventWindow.SetContent(eventLabel)
			deleteEventWindow.Show()
		})
		deleteEventForm := container.NewVBox(
			widget.NewLabel("ID:"),
			idEntry,
			deleteEventButton,
		)
		deleteEventWindow.SetContent(deleteEventForm)
		deleteEventWindow.Show()
	})

	content := container.NewVBox(
		eventList,
		viewAllEventsButton,
		viewEventByIDButton,
		addEventButton,
		updateEventButton,
		deleteEventButton,
	)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

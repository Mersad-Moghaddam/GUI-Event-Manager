package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/gorilla/mux"
)

// Event struct represents an event with basic information.
type Event struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

// Mock database to store events
var events []Event

// Handler to get all events (API)
func getEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

// Handler to get a single event by ID (API)
func getEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range events {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	http.Error(w, "Event not found", http.StatusNotFound)
}

// Handler to create a new event (API)
func createEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var event Event
	_ = json.NewDecoder(r.Body).Decode(&event)
	events = append(events, event)
	json.NewEncoder(w).Encode(event)
}

// Handler to update an event by ID (API)
func updateEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for i, item := range events {
		if item.ID == params["id"] {
			events = append(events[:i], events[i+1:]...)
			var updatedEvent Event
			_ = json.NewDecoder(r.Body).Decode(&updatedEvent)
			updatedEvent.ID = params["id"]
			events = append(events, updatedEvent)
			json.NewEncoder(w).Encode(updatedEvent)
			return
		}
	}
	http.Error(w, "Event not found", http.StatusNotFound)
}

// Handler to delete an event by ID (API)
func deleteEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for i, item := range events {
		if item.ID == params["id"] {
			events = append(events[:i], events[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Event not found", http.StatusNotFound)
}

// CLI Functions
func showMenu() {
	fmt.Println("\n--- Event Scheduling CLI ---")
	fmt.Println("1. View all events")
	fmt.Println("2. View an event by ID")
	fmt.Println("3. Add a new event")
	fmt.Println("4. Update an event by ID")
	fmt.Println("5. Delete an event by ID")
	fmt.Println("6. Exit CLI")
	fmt.Print("Select an option: ")
}

func viewAllEventsCLI() {
	if len(events) == 0 {
		fmt.Println("No events available.")
		return
	}
	fmt.Println("\nList of Events:")
	for _, event := range events {
		fmt.Printf("ID: %s | Title: %s | Description: %s | Date: %s\n", event.ID, event.Title, event.Description, event.Date.Format("2006-01-02 15:04"))
	}
}

func viewEventByIDCLI() {
	id := input("Enter event ID: ")
	for _, event := range events {
		if event.ID == id {
			fmt.Printf("ID: %s | Title: %s | Description: %s | Date: %s\n", event.ID, event.Title, event.Description, event.Date.Format("2006-01-02 15:04"))
			return
		}
	}
	fmt.Println("Event not found.")
}

func addEventCLI() {
	dateStr := input("Enter event date (YYYY-MM-DD HH:MM): ")
	date, _ := time.Parse("2006-01-02 15:04", dateStr)
	event := Event{
		ID:          input("Enter event ID: "),
		Title:       input("Enter event title: "),
		Description: input("Enter event description: "),
		Date:        date,
	}
	events = append(events, event)
	fmt.Println("Event added successfully.")
}

func updateEventByIDCLI() {
	id := input("Enter event ID to update: ")
	for i, event := range events {
		if event.ID == id {
			dateStr := input("Enter new event date (YYYY-MM-DD HH:MM): ")
			date, _ := time.Parse("2006-01-02 15:04", dateStr)
			events[i].Title = input("Enter new title: ")
			events[i].Description = input("Enter new description: ")
			events[i].Date = date
			fmt.Println("Event updated successfully.")
			return
		}
	}
	fmt.Println("Event not found.")
}

func deleteEventByIDCLI() {
	id := input("Enter event ID to delete: ")
	for i, event := range events {
		if event.ID == id {
			events = append(events[:i], events[i+1:]...)
			fmt.Println("Event deleted successfully.")
			return
		}
	}
	fmt.Println("Event not found.")
}

// Helper functions for CLI
func input(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

// GUI Functions
func startGUI() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Event Manager")

	// Event List
	eventList := widget.NewList(
		func() int { return len(events) },
		func() fyne.CanvasObject { return widget.NewLabel("Template") },
		func(i int, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(events[i].Title)
		},
	)

	// GUI Layout
	content := container.NewVBox(
		eventList,
	)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

func main() {
	// Start API server
	go func() {
		router := mux.NewRouter()
		router.HandleFunc("/events", getEvents).Methods("GET")
		router.HandleFunc("/events/{id}", getEvent).Methods("GET")
		router.HandleFunc("/events", createEvent).Methods("POST")
		router.HandleFunc("/events/{id}", updateEvent).Methods("PUT")
		router.HandleFunc("/events/{id}", deleteEvent).Methods("DELETE")
		log.Fatal(http.ListenAndServe(":8080", router))
	}()

	// Start GUI
	go startGUI()

	// Start CLI
	startCLI()
}

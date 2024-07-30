package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Event represents a calendar event
type Event struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Date      time.Time `json:"date"`
	UserID    int       `json:"user_id"`
}

var events = make(map[int]Event)
var nextID = 1

func main() {
	http.HandleFunc("/create_event", createEventHandler)
	http.HandleFunc("/update_event", updateEventHandler)
	http.HandleFunc("/delete_event", deleteEventHandler)
	http.HandleFunc("/events_for_day", eventsForDayHandler)
	http.HandleFunc("/events_for_week", eventsForWeekHandler)
	http.HandleFunc("/events_for_month", eventsForMonthHandler)
	http.Handle("/", loggingMiddleware(http.DefaultServeMux))

	port := ":8080"
	fmt.Printf("Starting server at port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// loggingMiddleware logs the details of each request
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
		log.Printf("Completed in %v", time.Since(start))
	})
}

// createEventHandler handles the creation of events
func createEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	title := r.FormValue("title")
	userID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		http.Error(w, `{"error": "Invalid user_id"}`, http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", r.FormValue("date"))
	if err != nil {
		http.Error(w, `{"error": "Invalid date format, should be YYYY-MM-DD"}`, http.StatusBadRequest)
		return
	}

	event := Event{
		ID:     nextID,
		Title:  title,
		Date:   date,
		UserID: userID,
	}
	nextID++
	events[event.ID] = event

	response := map[string]interface{}{"result": event}
	jsonResponse(w, response)
}

// updateEventHandler handles the updating of events
func updateEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil || id <= 0 {
		http.Error(w, `{"error": "Invalid event id"}`, http.StatusBadRequest)
		return
	}

	event, ok := events[id]
	if !ok {
		http.Error(w, `{"error": "Event not found"}`, http.StatusNotFound)
		return
	}

	title := r.FormValue("title")
	userID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		http.Error(w, `{"error": "Invalid user_id"}`, http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", r.FormValue("date"))
	if err != nil {
		http.Error(w, `{"error": "Invalid date format, should be YYYY-MM-DD"}`, http.StatusBadRequest)
		return
	}

	event.Title = title
	event.Date = date
	event.UserID = userID
	events[id] = event

	response := map[string]interface{}{"result": event}
	jsonResponse(w, response)
}

// deleteEventHandler handles the deletion of events
func deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil || id <= 0 {
		http.Error(w, `{"error": "Invalid event id"}`, http.StatusBadRequest)
		return
	}

	if _, ok := events[id]; !ok {
		http.Error(w, `{"error": "Event not found"}`, http.StatusNotFound)
		return
	}

	delete(events, id)
	response := map[string]interface{}{"result": "Event deleted"}
	jsonResponse(w, response)
}

// eventsForDayHandler handles the retrieval of events for a specific day
func eventsForDayHandler(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid date format, should be YYYY-MM-DD"}`, http.StatusBadRequest)
		return
	}

	var result []Event
	for _, event := range events {
		if event.Date.Year() == date.Year() && event.Date.YearDay() == date.YearDay() {
			result = append(result, event)
		}
	}

	jsonResponse(w, map[string]interface{}{"result": result})
}

// eventsForWeekHandler handles the retrieval of events for a specific week
func eventsForWeekHandler(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid date format, should be YYYY-MM-DD"}`, http.StatusBadRequest)
		return
	}

	var result []Event
	startOfWeek := date.AddDate(0, 0, -int(date.Weekday()))
	for _, event := range events {
		if event.Date.After(startOfWeek) && event.Date.Before(startOfWeek.AddDate(0, 0, 7)) {
			result = append(result, event)
		}
	}

	jsonResponse(w, map[string]interface{}{"result": result})
}

// eventsForMonthHandler handles the retrieval of events for a specific month
func eventsForMonthHandler(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse("2006-01", dateStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid date format, should be YYYY-MM"}`, http.StatusBadRequest)
		return
	}

	var result []Event
	for _, event := range events {
		if event.Date.Year() == date.Year() && event.Date.Month() == date.Month() {
			result = append(result, event)
		}
	}

	jsonResponse(w, map[string]interface{}{"result": result})
}

// jsonResponse serializes the response to JSON and writes it to the ResponseWriter
func jsonResponse(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

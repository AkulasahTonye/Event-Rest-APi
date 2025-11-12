package models

import (
	"time"

	"example.com/api-testing/db"
)

type Event struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"  binding:"required"`
	Description string    `json:"description" binding:"required"`
	Location    string    `json:"location" binding:"required"`
	DateTime    time.Time `json:"dateTime" binding:"required"`
	UserID      int64     `json:"userId"`
}

// Hooks for testing
var (
	EventSaveFunc          func(e *Event) error
	EventUpdateFunc        func(e *Event) error
	EventDeleteFunc        func(e *Event) error
	GetEventByIDFunc       func(id int64) (*Event, error)
	GetAllEventsFunc       func() ([]Event, error)
	CancelRegistrationFunc func(e *Event, userID int64) error
	RegisterFunc           func(e *Event, userID int64) error
)

var (
	GetAllEvents = getAllEvents
	GetEventByID = getEventByID
)

func (e *Event) Save() error {
	if EventSaveFunc != nil {
		return EventSaveFunc(e)
	}

	query := `
	INSERT INTO events(name, description, location, dateTime, user_id) 
	VALUES (?, ?, ?, ?, ?)`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID)
	if err != nil {

		return err
	}
	id, err := result.LastInsertId()
	e.ID = id
	return err

}

func getAllEvents() ([]Event, error) {
	if GetAllEventsFunc != nil {
		return GetAllEventsFunc()
	}

	query := "SELECT * FROM events"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event

	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)

		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func getEventByID(id int64) (*Event, error) {
	if GetEventByIDFunc != nil {
		return GetEventByIDFunc(id)
	}

	query := "SELECT * FROM events WHERE id =?"
	row := db.DB.QueryRow(query, id)

	var event Event
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
	if err != nil {
		return nil, err
	}
	return &event, err
}

func (e *Event) Update() error {
	if EventUpdateFunc != nil {
		return EventUpdateFunc(e)
	}

	query := `
     Update events  
SET name = ?, description = ?, location = ?,dateTime =?
WHERE id = ?
`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.ID)
	return err

}

func (e *Event) Delete() error {
	if EventDeleteFunc != nil {
		return EventDeleteFunc(e)
	}

	query := "DELETE FROM main.events WHERE events.id = ?"
	Stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}
	defer Stmt.Close()

	_, err = Stmt.Exec(e.ID)
	return err
}
func (e *Event) Register(userId int64) error {
	if CancelRegistrationFunc != nil {
		return CancelRegistrationFunc(e, userId)
	}

	query := `INSERT INTO registration(event_id, user_id) VALUES (?,?)`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(e.ID, userId)

	return err
}

func (e *Event) CancelRegistration(userId int64) error {
	if CancelRegistrationFunc != nil {
		return CancelRegistrationFunc(e, userId)
	}

	query := `DELETE FROM registration WHERE event_id = ? AND user_id = ?`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(e.ID, userId)

	return err
}

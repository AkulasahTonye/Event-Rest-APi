package models

import (
	"testing"
	"time"

	"example.com/api-testing/db"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func setUpMockDB(t *testing.T) (sqlmock.Sqlmock, func()) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a mock sql database", err)
	}
	db.DB = mockDB // assign mock DB to global db.DB

	return mock, func() {
		mockDB.Close()
	}
}

func TestSaveEvent(t *testing.T) {
	mock, closeFn := setUpMockDB(t)
	defer closeFn()

	events := Event{
		Name:        "Test Event",
		Description: "Test Save Function",
		Location:    "Online",
		DateTime:    time.Now(),
		UserID:      1,
	}
	mock.ExpectPrepare("INSERT INTO events").ExpectExec().WithArgs(events.Name, events.Description, events.Location, events.DateTime, events.UserID).WillReturnResult(sqlmock.NewResult(1, 1))

	err := events.Save()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), events.ID)
}

func TestGetAllEvents(t *testing.T) {
	mock, closeFn := setUpMockDB(t)
	defer closeFn()

	rows := sqlmock.NewRows([]string{"id", "name", "description", "location", "dateTime", "user_id"}).
		AddRow(1, "Event A", "Desc A", "Loc A", time.Now(), 2).
		AddRow(2, "Event B", "Desc B", "Loc B", time.Now(), 3)

	mock.ExpectQuery("SELECT \\* FROM events").WillReturnRows(rows)

	events, err := getAllEvents()
	assert.NoError(t, err)
	assert.Len(t, events, 2)
	assert.Equal(t, "Event A", events[0].Name)
}

func TestGetEventByID(t *testing.T) {
	mock, closeFn := setUpMockDB(t)
	defer closeFn()

	row := sqlmock.NewRows([]string{"id", "name", "description", "location", "dateTime", "user_id"}).
		AddRow(1, "Event X", "Desc X", "Loc X", time.Now(), 2)

	mock.ExpectQuery("SELECT \\* FROM events WHERE id =\\?").
		WithArgs(int64(1)).
		WillReturnRows(row)

	event, err := getEventByID(1)
	assert.NoError(t, err)
	assert.Equal(t, "Event X", event.Name)
}

func TestUpdateEvent(t *testing.T) {
	mock, closeFn := setUpMockDB(t)
	defer closeFn()

	event := Event{
		ID:          1,
		Name:        "Updated Name",
		Description: "Updated Desc",
		Location:    "New Location",
		DateTime:    time.Now(),
	}

	mock.ExpectPrepare("Update events").
		ExpectExec().
		WithArgs(event.Name, event.Description, event.Location, event.DateTime, event.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := event.Update()
	assert.NoError(t, err)
}

func TestDeleteEvent(t *testing.T) {
	mock, closeFn := setUpMockDB(t)
	defer closeFn()

	event := Event{ID: 10}

	mock.ExpectPrepare("DELETE FROM main.events").
		ExpectExec().
		WithArgs(event.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := event.Delete()
	assert.NoError(t, err)
}

func TestRegisterEvent(t *testing.T) {
	mock, closeFn := setUpMockDB(t)
	defer closeFn()

	event := Event{ID: 5}

	mock.ExpectPrepare("INSERT INTO registration").
		ExpectExec().
		WithArgs(event.ID, int64(7)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := event.Register(7)
	assert.NoError(t, err)
}

func TestCancelRegistration(t *testing.T) {
	mock, closeFn := setUpMockDB(t)
	defer closeFn()

	event := Event{ID: 8}

	mock.ExpectPrepare("DELETE FROM registration").
		ExpectExec().
		WithArgs(event.ID, int64(2)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := event.CancelRegistration(2)
	assert.NoError(t, err)
}

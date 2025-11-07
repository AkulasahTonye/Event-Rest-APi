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

	mock.ExpectQuery("SELECT * FROM events").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "location", "dateTime", "userId"}).AddRow(1, "Test Event", "Test Save Function", "Online", time.Now(), 1))

}

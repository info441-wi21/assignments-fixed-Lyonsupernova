package users

import (
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

type RowObj struct {
	id        int64
	email     string
	passHash  []byte
	userName  string
	firstName string
	lastName  string
	photoURL  string
}

func TestGetByID(t *testing.T) {
	cases := []struct {
		name        string
		rows        []RowObj
		idToGet     int64
		expectError bool
	}{
		{
			"Working case",
			[]RowObj{
				RowObj{
					int64(1),
					"test@test.com",
					[]byte("passhash123"),
					"username",
					"firstname",
					"lastname",
					"photourl",
				},
			},
			int64(1),
			false,
		},
		{
			"No rows found",
			[]RowObj{},
			int64(2),
			true,
		},
		{
			"Large id",
			[]RowObj{
				RowObj{
					int64(1251241251),
					"test@test.com",
					[]byte("passhash123"),
					"username",
					"firstname",
					"lastname",
					"photourl",
				},
			},
			int64(1251241251),
			false,
		},
	}

	for _, c := range cases {
		// initialize new db mock for each test
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		expectedRow := &RowObj{}

		// add rows to mock db
		rows := sqlmock.NewRows([]string{"ID", "Email", "PassHash", "UserName", "FirstName", "LastName", "PhotoURL"})
		for _, row := range c.rows {
			if c.idToGet == row.id {
				expectedRow.id = row.id
				expectedRow.email = row.email
				expectedRow.passHash = row.passHash
				expectedRow.userName = row.userName
				expectedRow.firstName = row.firstName
				expectedRow.lastName = row.lastName
				expectedRow.photoURL = row.photoURL
			}
			rows.AddRow(row.id,
				row.email,
				row.passHash,
				row.userName,
				row.firstName,
				row.lastName,
				row.photoURL)
		}

		// Set up expected query
		mock.ExpectQuery("^select (.+) from Users where ID = ").WithArgs(c.idToGet).WillReturnRows(rows)

		// Create MySQLStore using the mock db
		mysqlstore := &MySQLStore{db}

		// Run the query
		user, errG := mysqlstore.GetByID(c.idToGet)
		if !c.expectError && errG != nil {
			t.Errorf("Unexpected error occured on test [%s]: %v", c.name, errG)
		}

		if c.expectError && errG == nil {
			t.Errorf("Expected error but got %v", errG)
		}

		// Test if user matches the correct row
		if !c.expectError && (user.ID != expectedRow.id || user.Email != expectedRow.email ||
			!reflect.DeepEqual(user.PassHash, expectedRow.passHash) || user.UserName != expectedRow.userName ||
			user.FirstName != expectedRow.firstName || user.LastName != expectedRow.lastName ||
			user.PhotoURL != expectedRow.photoURL) {
			t.Errorf("Error, invalid match in test [%s]", c.name)
		}

		// Test a mock error case
		mock.ExpectQuery("^select (.+) from Users where ID = ").WithArgs(c.idToGet).
			WillReturnError(errors.New("some error"))
		user, errG = mysqlstore.GetByID(c.idToGet)
		if user != nil || errG == nil {
			t.Errorf("Expected error [%v] but got [%v] instead", errors.New("some error"), errG)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}
	}
}

func TestGetByEmail(t *testing.T) {
	cases := []struct {
		name        string
		rows        []RowObj
		emailToGet  string
		expectError bool
	}{
		{
			"Working case",
			[]RowObj{
				RowObj{
					int64(1),
					"test@test.com",
					[]byte("passhash123"),
					"username",
					"firstname",
					"lastname",
					"photourl",
				},
			},
			"test@test.com",
			false,
		},
		{
			"No rows found",
			[]RowObj{},
			"test@test.com",
			true,
		},
	}

	for _, c := range cases {
		// initialize new db mock for each test
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		expectedRow := &RowObj{}

		// add rows to mock db
		rows := sqlmock.NewRows([]string{"ID", "Email", "PassHash", "UserName", "FirstName", "LastName", "PhotoURL"})
		for _, row := range c.rows {
			if c.emailToGet == row.email {
				expectedRow.id = row.id
				expectedRow.email = row.email
				expectedRow.passHash = row.passHash
				expectedRow.userName = row.userName
				expectedRow.firstName = row.firstName
				expectedRow.lastName = row.lastName
				expectedRow.photoURL = row.photoURL
			}
			rows.AddRow(row.id,
				row.email,
				row.passHash,
				row.userName,
				row.firstName,
				row.lastName,
				row.photoURL)
		}

		// Set up expected query
		mock.ExpectQuery("^select (.+) from Users where Email = ").WithArgs(c.emailToGet).WillReturnRows(rows)

		// Create MySQLStore using the mock db
		mysqlstore := &MySQLStore{db}

		// Run the query
		user, errG := mysqlstore.GetByEmail(c.emailToGet)
		if !c.expectError && errG != nil {
			t.Errorf("Unexpected error occured on test [%s]: %v", c.name, errG)
		}

		if c.expectError && errG == nil {
			t.Errorf("Expected error but got %v", errG)
		}

		// Test if user matches the correct row
		if !c.expectError && (user.ID != expectedRow.id || user.Email != expectedRow.email ||
			!reflect.DeepEqual(user.PassHash, expectedRow.passHash) || user.UserName != expectedRow.userName ||
			user.FirstName != expectedRow.firstName || user.LastName != expectedRow.lastName ||
			user.PhotoURL != expectedRow.photoURL) {
			t.Errorf("Error, invalid match in test [%s]", c.name)
		}

		// Test a mock error case
		mock.ExpectQuery("^select (.+) from Users where Email = ").WithArgs(c.emailToGet).
			WillReturnError(errors.New("some error"))
		user, errG = mysqlstore.GetByEmail(c.emailToGet)
		if user != nil || errG == nil {
			t.Errorf("Expected error [%v] but got [%v] instead", errors.New("some error"), errG)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}
	}
}

func TestGetByUserName(t *testing.T) {
	cases := []struct {
		name          string
		rows          []RowObj
		userNameToGet string
		expectError   bool
	}{
		{
			"Working case",
			[]RowObj{
				RowObj{
					int64(1),
					"test@test.com",
					[]byte("passhash123"),
					"username",
					"firstname",
					"lastname",
					"photourl",
				},
			},
			"username",
			false,
		},
		{
			"No rows found",
			[]RowObj{},
			"username",
			true,
		},
	}

	for _, c := range cases {
		// initialize new db mock for each test
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		expectedRow := &RowObj{}

		// add rows to mock db
		rows := sqlmock.NewRows([]string{"ID", "Email", "PassHash", "UserName", "FirstName", "LastName", "PhotoURL"})
		for _, row := range c.rows {
			if c.userNameToGet == row.userName {
				expectedRow.id = row.id
				expectedRow.email = row.email
				expectedRow.passHash = row.passHash
				expectedRow.userName = row.userName
				expectedRow.firstName = row.firstName
				expectedRow.lastName = row.lastName
				expectedRow.photoURL = row.photoURL
			}
			rows.AddRow(row.id,
				row.email,
				row.passHash,
				row.userName,
				row.firstName,
				row.lastName,
				row.photoURL)
		}

		// Set up expected query
		mock.ExpectQuery("^select (.+) from Users where UserName = ").WithArgs(c.userNameToGet).WillReturnRows(rows)

		// Create MySQLStore using the mock db
		mysqlstore := &MySQLStore{db}

		// Run the query
		user, errG := mysqlstore.GetByUserName(c.userNameToGet)
		if !c.expectError && errG != nil {
			t.Errorf("Unexpected error occured on test [%s]: %v", c.name, errG)
		}

		if c.expectError && errG == nil {
			t.Errorf("Expected error but got %v", errG)
		}

		// Test if user matches the correct row
		if !c.expectError && (user.ID != expectedRow.id || user.Email != expectedRow.email ||
			!reflect.DeepEqual(user.PassHash, expectedRow.passHash) || user.UserName != expectedRow.userName ||
			user.FirstName != expectedRow.firstName || user.LastName != expectedRow.lastName ||
			user.PhotoURL != expectedRow.photoURL) {
			t.Errorf("Error, invalid match in test [%s]", c.name)
		}

		// Test a mock error case
		mock.ExpectQuery("^select (.+) from Users where UserName = ").WithArgs(c.userNameToGet).
			WillReturnError(errors.New("some error"))
		user, errG = mysqlstore.GetByUserName(c.userNameToGet)
		if user != nil || errG == nil {
			t.Errorf("Expected error [%v] but got [%v] instead", errors.New("some error"), errG)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}
	}
}

func TestInsert(t *testing.T) {
	cases := []struct {
		name        string
		user        User
		idToSet     int64
		expectError bool
	}{
		{
			"Working case",
			User{
				Email:     "test@test.com",
				PassHash:  []byte("passhash"),
				UserName:  "username",
				FirstName: "firstname",
				LastName:  "lastname",
				PhotoURL:  "photourl",
			},
			1,
			false,
		},
		{
			"Empty user",
			User{},
			1,
			false,
		},
	}

	for _, c := range cases {
		// initialize new db mock for each test
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		// Create MySQLStore using the mock db
		mysqlstore := &MySQLStore{db}

		mock.ExpectExec("^insert into Users(.+) values ").WillReturnResult(sqlmock.NewResult(c.idToSet, 1))
		user, err := mysqlstore.Insert(&c.user)
		if !c.expectError && err != nil {
			t.Errorf("Unexpected error occured on test [%s]: [%v]", c.name, err)
		}

		if c.expectError && err == nil {
			t.Errorf("Error expected but got [%v]", err)
		}

		if user.ID != c.idToSet {
			t.Errorf("ID not as expected. Expected [%v] but got [%v]", c.idToSet, user.ID)
		}

		mock.ExpectExec("^insert into Users(.+) values ").WillReturnError(errors.New("some error"))
		user, errG := mysqlstore.Insert(&c.user)
		if user != nil || errG == nil {
			t.Errorf("Expected error [%v] but got [%v] instead", errors.New("some error"), errG)
		}

		mock.ExpectExec("^insert into Users(.+) values ").WillReturnResult(sqlmock.NewErrorResult(errors.New("some error")))
		user, errG = mysqlstore.Insert(&c.user)
		if user != nil || errG == nil {
			t.Errorf("Expected error [%v] but got [%v] instead", errors.New("some error"), errG)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}
	}
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mysqlstore := &MySQLStore{db}

	// Set up first test for a working case
	mock.ExpectExec("^update Users set FirstName(.+), LastName(.+) where ID").WillReturnResult(sqlmock.NewResult(1, 1))
	row := &RowObj{
		int64(1),
		"email@email.com",
		[]byte("passhash"),
		"username",
		"A",
		"B",
		"photourl",
	}
	rows := sqlmock.NewRows([]string{"ID", "Email", "PassHash", "UserName", "FirstName", "LastName", "PhotoURL"}).AddRow(row.id,
		row.email,
		row.passHash,
		row.userName,
		row.firstName,
		row.lastName,
		row.photoURL)
	mock.ExpectQuery("^select (.+) from Users where ID = ").WithArgs(1).WillReturnRows(rows)

	user, errG := mysqlstore.Update(1, &Updates{FirstName: "A", LastName: "B"})
	if err != nil {
		t.Errorf("Unexpected error occured. Got [%v]", errG)
	}

	if user.FirstName != "A" || user.LastName != "B" {
		t.Errorf("Update failed. Expected [%v, %v]. Received [%v, %v].", "A", "B", user.FirstName, user.LastName)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}

	// Check for SQL function execution error
	mock.ExpectExec("^update Users set FirstName(.+), LastName(.+) where ID").WillReturnError(errors.New("some error"))
	user, errG = mysqlstore.Update(1, &Updates{FirstName: "A", LastName: "B"})
	if errG == nil {
		t.Errorf("Expected Error but got [%v]", err)
	}

	// Check for rows affected error
	mock.ExpectExec("^update Users set FirstName(.+), LastName(.+) where ID").WillReturnResult(sqlmock.NewErrorResult(errors.New("some error")))
	user, errG = mysqlstore.Update(1, &Updates{FirstName: "A", LastName: "B"})
	if errG == nil {
		t.Errorf("Expected Error but got [%v]", err)
	}

	// Check for multiple rows affected for some reason
	mock.ExpectExec("^update Users set FirstName(.+), LastName(.+) where ID").WillReturnResult(sqlmock.NewResult(1, 3))
	user, errG = mysqlstore.Update(1, &Updates{FirstName: "A", LastName: "B"})
	if errG == nil {
		t.Errorf("Expected Error but got [%v]", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mysqlstore := &MySQLStore{db}

	// Set up first test for a working case
	mock.ExpectExec("^delete from Users where ID").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	err = mysqlstore.Delete(1)
	if err != nil {
		t.Errorf("Unexpected error [%v]", err)
	}

	mock.ExpectExec("^delete from Users where ID").WillReturnError(errors.New("some error"))
	err = mysqlstore.Delete(1)
	if err == nil {
		t.Errorf("Expected Error but got [%v]", err)
	}

	mock.ExpectExec("^delete from Users where ID").WillReturnResult(sqlmock.NewErrorResult(errors.New("some error")))
	err = mysqlstore.Delete(1)
	if err == nil {
		t.Errorf("Expected Error but got [%v]", err)
	}

	mock.ExpectExec("^delete from Users where ID").WillReturnResult(sqlmock.NewResult(1, 3))
	err = mysqlstore.Delete(1)
	if err == nil {
		t.Errorf("Expected Error but got [%v]", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

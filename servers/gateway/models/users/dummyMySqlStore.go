package users

import (
	"time"
)

// DummyMySQLStore strcut
type DummyMySQLStore struct {
}

// NewDummySQL method
func NewDummySQL() *DummyMySQLStore {
	return &DummyMySQLStore{}
}

// GetByEmail method
func (du *DummyMySQLStore) GetByEmail(email string) (*User, error) {
	return nil, nil
}

// GetByID method
func (du *DummyMySQLStore) GetByID(id int64) (*User, error) {
	return nil, nil
}

// GetByUserName method
func (du *DummyMySQLStore) GetByUserName(username string) (*User, error) {
	return nil, nil

}

// Delete method
func (du *DummyMySQLStore) Delete(id int64) error {
	return nil
}

// Update method
func (du *DummyMySQLStore) Update(id int64, updates *Updates) (*User, error) {
	return nil, nil
}

// Insert method
func (du *DummyMySQLStore) Insert(user *User) (*User, error) {
	if user.ID == 0 {
		user.ID = 1
		return user, nil
	}
	return nil, nil
}

// LogSignIn method
func (du *DummyMySQLStore) LogSignIn(user *User, dateTime time.Time, clientIP string) error {
	return nil
}

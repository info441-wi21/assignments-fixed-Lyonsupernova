package models

// a type to represent a single channel
type channel struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json: description`
	Private     bool    `json: private`
	Members     []*User `json: members` // not sure what is a 'member'
	CreatedAt   string  `json: createdAt`
	Creator     string  `json: creator`
	editedAt    string  `json: editedAt`
}

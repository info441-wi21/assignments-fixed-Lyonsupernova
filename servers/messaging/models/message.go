package models

// a type to represent a single message
type message struct {
	ID        int64  `json:"id"`
	ChannelID int64  `json:channelID`
	Body      string `json:body`
	CreatedAt string `json:createdAt`
	Creator   string `json:creator`
	EditedAt  string `json:editedAt`
}

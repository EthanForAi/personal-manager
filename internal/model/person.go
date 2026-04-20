package model

type Person struct {
	UserID string `json:"userid"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
}

type DeleteResponse struct {
	Deleted bool `json:"deleted"`
}

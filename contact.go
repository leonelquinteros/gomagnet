package gomagnet

// Contact data
type Contact struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Question  string `json:"question"`
	CreatedAt string `json:"created_at"`
}

// ContactResponse data
type ContactResponse struct {
	Status  bool      `json:"status"`
	Message int       `json:"message"`
	Data    []Contact `json:"data"`
}

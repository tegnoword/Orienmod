package domain

type Studen struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"lastname"`
	Email    string `json:"email"`
}

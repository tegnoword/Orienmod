package domain

type Course struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Section     string `json:"section"`
	Description string `json:"description"`
}

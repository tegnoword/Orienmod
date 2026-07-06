package domain

type Grande struct {
	ID            int    `json:"id"`
	IdStuden      int    `json:"id_studen"`
	IdTask        int    `json:"id_task"`
	GradeObtained int    `json:"grade_obtained"`
	Status        string `json:"status"`
}

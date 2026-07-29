package graphql

import (
	"fmt"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/tegnoword/orienmod/internal/core/domain"
)

func (r *Resolver) CoursesResolver(p graphql.ResolveParams) (interface{}, error) {
	email, ok := p.Context.Value("email").(string)
	if !ok || email == "" {
		return nil, fmt.Errorf("usuario no autenticado")
	}

	courses, err := r.GoogleAdapter.GetAllCourses(p.Context, email)
	if err != nil {
		return nil, err
	}

	// Convertir a []interface{} para GraphQL
	result := make([]interface{}, len(courses))
	for i, c := range courses {
		result[i] = c
	}
	return result, nil
}

// CourseResolver - Resuelve la query course
func (r *Resolver) CourseResolver(p graphql.ResolveParams) (interface{}, error) {
	email, ok := p.Context.Value("email").(string)
	if !ok || email == "" {
		return nil, fmt.Errorf("usuario no autenticado")
	}

	id, ok := p.Args["id"].(string)
	if !ok || id == "" {
		return nil, fmt.Errorf("id es requerido")
	}

	courses, err := r.GoogleAdapter.GetAllCourses(p.Context, email)
	if err != nil {
		return nil, err
	}

	for _, c := range courses {
		if c.ID == id {
			return c, nil
		}
	}
	return nil, fmt.Errorf("curso no encontrado")
}

// SearchCoursesResolver - Resuelve la query searchCourses
func (r *Resolver) SearchCoursesResolver(p graphql.ResolveParams) (interface{}, error) {
	email, ok := p.Context.Value("email").(string)
	if !ok || email == "" {
		return nil, fmt.Errorf("usuario no autenticado")
	}

	query, ok := p.Args["query"].(string)
	if !ok {
		return nil, fmt.Errorf("query es requerido")
	}

	courses, err := r.GoogleAdapter.GetAllCourses(p.Context, email)
	if err != nil {
		return nil, err
	}

	var results []domain.Course
	for _, c := range courses {
		if contains(c.Name, query) || contains(c.Description, query) {
			results = append(results, c)
		}
	}

	result := make([]interface{}, len(results))
	for i, c := range results {
		result[i] = c
	}
	return result, nil
}

// StudentsResolver - Resuelve la query students
func (r *Resolver) StudentsResolver(p graphql.ResolveParams) (interface{}, error) {
	email, ok := p.Context.Value("email").(string)
	if !ok || email == "" {
		return nil, fmt.Errorf("usuario no autenticado")
	}

	courseID, ok := p.Args["courseId"].(string)
	if !ok || courseID == "" {
		return nil, fmt.Errorf("courseId es requerido")
	}

	students, err := r.GoogleAdapter.GetStudentsByCourse(p.Context, email, courseID)
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(students))
	for i, s := range students {
		result[i] = s
	}
	return result, nil
}

// SearchStudentsResolver - Resuelve la query searchStudents
func (r *Resolver) SearchStudentsResolver(p graphql.ResolveParams) (interface{}, error) {
	email, ok := p.Context.Value("email").(string)
	if !ok || email == "" {
		return nil, fmt.Errorf("usuario no autenticado")
	}

	courseID, ok := p.Args["courseId"].(string)
	if !ok || courseID == "" {
		return nil, fmt.Errorf("courseId es requerido")
	}

	query, ok := p.Args["query"].(string)
	if !ok {
		return nil, fmt.Errorf("query es requerido")
	}

	students, err := r.GoogleAdapter.GetStudentsByCourse(p.Context, email, courseID)
	if err != nil {
		return nil, err
	}

	var results []domain.Student
	for _, s := range students {
		if contains(s.Name, query) || contains(s.Email, query) {
			results = append(results, s)
		}
	}

	result := make([]interface{}, len(results))
	for i, s := range results {
		result[i] = s
	}
	return result, nil
}

// TasksResolver - Resuelve la query tasks
func (r *Resolver) TasksResolver(p graphql.ResolveParams) (interface{}, error) {
	email, ok := p.Context.Value("email").(string)
	if !ok || email == "" {
		return nil, fmt.Errorf("usuario no autenticado")
	}

	courseID, ok := p.Args["courseId"].(string)
	if !ok || courseID == "" {
		return nil, fmt.Errorf("courseId es requerido")
	}

	tasks, err := r.GoogleAdapter.GetTasksByCourse(p.Context, email, courseID)
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(tasks))
	for i, t := range tasks {
		result[i] = t
	}
	return result, nil
}

// TaskResolver - Resuelve la query task
func (r *Resolver) TaskResolver(p graphql.ResolveParams) (interface{}, error) {
	email, ok := p.Context.Value("email").(string)
	if !ok || email == "" {
		return nil, fmt.Errorf("usuario no autenticado")
	}

	taskID, ok := p.Args["id"].(string)
	if !ok || taskID == "" {
		return nil, fmt.Errorf("id es requerido")
	}

	courseID, ok := p.Args["courseId"].(string)
	if !ok || courseID == "" {
		return nil, fmt.Errorf("courseId es requerido")
	}

	task, err := r.GoogleAdapter.GetTaskByID(p.Context, email, courseID, taskID)
	if err != nil {
		return nil, err
	}
	return task, nil
}

// TaskSubmissionsResolver - Resuelve la query taskSubmissions
func (r *Resolver) TaskSubmissionsResolver(p graphql.ResolveParams) (interface{}, error) {
	email, ok := p.Context.Value("email").(string)
	if !ok || email == "" {
		return nil, fmt.Errorf("usuario no autenticado")
	}

	taskID, ok := p.Args["taskId"].(string)
	if !ok || taskID == "" {
		return nil, fmt.Errorf("taskId es requerido")
	}

	courseID, ok := p.Args["courseId"].(string)
	if !ok || courseID == "" {
		return nil, fmt.Errorf("courseId es requerido")
	}

	submissions, err := r.GoogleAdapter.GetTaskSubmissions(p.Context, email, courseID, taskID)
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(submissions))
	for i, s := range submissions {
		result[i] = s
	}
	return result, nil
}

// CheckAuthResolver - Resuelve la query checkAuth
func (r *Resolver) CheckAuthResolver(p graphql.ResolveParams) (interface{}, error) {
	email, ok := p.Args["email"].(string)
	if !ok || email == "" {
		return map[string]interface{}{
			"authenticated": false,
			"email":         email,
			"valid":         false,
			"error":         "email requerido",
		}, nil
	}

	token, err := r.TokenStore.GetToken(p.Context, email)
	if err != nil {
		return map[string]interface{}{
			"authenticated": false,
			"email":         email,
			"valid":         false,
			"error":         err.Error(),
		}, nil
	}

	expiry := token.Expiry.Format(time.RFC3339)
	return map[string]interface{}{
		"authenticated": true,
		"email":         email,
		"tokenExpiry":   expiry,
		"valid":         token.Valid(),
	}, nil
}

// ============================================
// MUTATIONS
// ============================================

// CreateCourseResolver - Resuelve la mutation createCourse
func (r *Resolver) CreateCourseResolver(p graphql.ResolveParams) (interface{}, error) {
	email, ok := p.Context.Value("email").(string)
	if !ok || email == "" {
		return nil, fmt.Errorf("usuario no autenticado")
	}

	input, ok := p.Args["input"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("input inválido")
	}

	name, _ := input["name"].(string)
	section, _ := input["section"].(string)
	description, _ := input["description"].(string)

	if name == "" {
		return nil, fmt.Errorf("name es requerido")
	}

	course := domain.Course{
		Name:        name,
		Section:     section,
		Description: description,
	}

	err := r.GoogleAdapter.SaveCourse(p.Context, email, course)
	if err != nil {
		return nil, err
	}

	return course, nil
}

// UpdateCourseResolver - Resuelve la mutation updateCourse
func (r *Resolver) UpdateCourseResolver(p graphql.ResolveParams) (interface{}, error) {
	email, ok := p.Context.Value("email").(string)
	if !ok || email == "" {
		return nil, fmt.Errorf("usuario no autenticado")
	}

	input, ok := p.Args["input"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("input inválido")
	}

	id, _ := input["id"].(string)
	name, _ := input["name"].(string)
	section, _ := input["section"].(string)
	description, _ := input["description"].(string)

	if id == "" {
		return nil, fmt.Errorf("id es requerido")
	}
	if name == "" {
		return nil, fmt.Errorf("name es requerido")
	}

	course := domain.Course{
		ID:          id,
		Name:        name,
		Section:     section,
		Description: description,
	}

	err := r.GoogleAdapter.UpdateCourse(p.Context, email, course)
	if err != nil {
		return nil, err
	}

	return course, nil
}

// DeleteCourseResolver - Resuelve la mutation deleteCourse
func (r *Resolver) DeleteCourseResolver(p graphql.ResolveParams) (interface{}, error) {
	email, ok := p.Context.Value("email").(string)
	if !ok || email == "" {
		return false, fmt.Errorf("usuario no autenticado")
	}

	id, ok := p.Args["id"].(string)
	if !ok || id == "" {
		return false, fmt.Errorf("id es requerido")
	}

	err := r.GoogleAdapter.DeleteCourse(p.Context, email, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

// SyncCourseResolver - Resuelve la mutation syncCourse
func (r *Resolver) SyncCourseResolver(p graphql.ResolveParams) (interface{}, error) {
	email, ok := p.Context.Value("email").(string)
	if !ok || email == "" {
		return 0, fmt.Errorf("usuario no autenticado")
	}

	id, ok := p.Args["id"].(string)
	if !ok || id == "" {
		return 0, fmt.Errorf("id es requerido")
	}

	students, err := r.GoogleAdapter.GetStudentsByCourse(p.Context, email, id)
	if err != nil {
		return 0, err
	}
	return len(students), nil
}

// AddStudentResolver - Resuelve la mutation addStudent
func (r *Resolver) AddStudentResolver(p graphql.ResolveParams) (interface{}, error) {
	email, ok := p.Context.Value("email").(string)
	if !ok || email == "" {
		return nil, fmt.Errorf("usuario no autenticado")
	}

	input, ok := p.Args["input"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("input inválido")
	}

	courseID, _ := input["courseId"].(string)
	studentID, _ := input["studentId"].(string)
	studentEmail, _ := input["email"].(string)
	name, _ := input["name"].(string)

	if courseID == "" || studentID == "" || studentEmail == "" || name == "" {
		return nil, fmt.Errorf("todos los campos son requeridos")
	}

	student := domain.Student{
		ID:       studentID,
		Email:    studentEmail,
		Name:     name,
		CourseID: courseID,
	}

	err := r.GoogleAdapter.SaveStudent(p.Context, email, courseID, student)
	if err != nil {
		return nil, err
	}
	return student, nil
}

// DeleteStudentResolver - Resuelve la mutation deleteStudent
func (r *Resolver) DeleteStudentResolver(p graphql.ResolveParams) (interface{}, error) {
	email, ok := p.Context.Value("email").(string)
	if !ok || email == "" {
		return false, fmt.Errorf("usuario no autenticado")
	}

	courseID, ok := p.Args["courseId"].(string)
	if !ok || courseID == "" {
		return false, fmt.Errorf("courseId es requerido")
	}

	studentID, ok := p.Args["studentId"].(string)
	if !ok || studentID == "" {
		return false, fmt.Errorf("studentId es requerido")
	}

	err := r.GoogleAdapter.DeleteStudent(p.Context, email, courseID, studentID)
	if err != nil {
		return false, err
	}
	return true, nil
}

// CreateTaskResolver - Resuelve la mutation createTask
func (r *Resolver) CreateTaskResolver(p graphql.ResolveParams) (interface{}, error) {
	email, ok := p.Context.Value("email").(string)
	if !ok || email == "" {
		return nil, fmt.Errorf("usuario no autenticado")
	}

	input, ok := p.Args["input"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("input inválido")
	}

	courseID, _ := input["courseId"].(string)
	title, _ := input["title"].(string)
	description, _ := input["description"].(string)
	dueDateStr, _ := input["dueDate"].(string)
	maxPoints, _ := input["maxPoints"].(float64)
	workType, _ := input["workType"].(string)

	if courseID == "" || title == "" {
		return nil, fmt.Errorf("courseId y title son requeridos")
	}

	dueDate, err := time.Parse(time.RFC3339, dueDateStr)
	if err != nil {
		return nil, fmt.Errorf("formato de fecha inválido: %v", err)
	}

	task := domain.Task{
		CourseID:    courseID,
		Title:       title,
		Description: description,
		DueDate:     dueDate,
		MaxPoints:   maxPoints,
		WorkType:    workType,
	}

	taskID, err := r.GoogleAdapter.CreateTask(p.Context, email, task)
	if err != nil {
		return nil, err
	}

	task.ID = taskID
	return task, nil
}

// UpdateTaskResolver - Resuelve la mutation updateTask
func (r *Resolver) UpdateTaskResolver(p graphql.ResolveParams) (interface{}, error) {
	email, ok := p.Context.Value("email").(string)
	if !ok || email == "" {
		return nil, fmt.Errorf("usuario no autenticado")
	}

	input, ok := p.Args["input"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("input inválido")
	}

	id, _ := input["id"].(string)
	courseID, _ := input["courseId"].(string)
	title, _ := input["title"].(string)
	description, _ := input["description"].(string)
	dueDateStr, _ := input["dueDate"].(string)
	maxPoints, _ := input["maxPoints"].(float64)
	state, _ := input["state"].(string)
	workType, _ := input["workType"].(string)

	if id == "" || courseID == "" || title == "" {
		return nil, fmt.Errorf("id, courseId y title son requeridos")
	}

	dueDate, err := time.Parse(time.RFC3339, dueDateStr)
	if err != nil {
		return nil, fmt.Errorf("formato de fecha inválido: %v", err)
	}

	task := domain.Task{
		ID:          id,
		CourseID:    courseID,
		Title:       title,
		Description: description,
		DueDate:     dueDate,
		MaxPoints:   maxPoints,
		State:       state,
		WorkType:    workType,
	}

	err = r.GoogleAdapter.UpdateTask(p.Context, email, task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

// DeleteTaskResolver - Resuelve la mutation deleteTask
func (r *Resolver) DeleteTaskResolver(p graphql.ResolveParams) (interface{}, error) {
	email, ok := p.Context.Value("email").(string)
	if !ok || email == "" {
		return false, fmt.Errorf("usuario no autenticado")
	}

	id, ok := p.Args["id"].(string)
	if !ok || id == "" {
		return false, fmt.Errorf("id es requerido")
	}

	courseID, ok := p.Args["courseId"].(string)
	if !ok || courseID == "" {
		return false, fmt.Errorf("courseId es requerido")
	}

	err := r.GoogleAdapter.DeleteTask(p.Context, email, courseID, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GradeTaskResolver - Resuelve la mutation gradeTask
func (r *Resolver) GradeTaskResolver(p graphql.ResolveParams) (interface{}, error) {
	email, ok := p.Context.Value("email").(string)
	if !ok || email == "" {
		return false, fmt.Errorf("usuario no autenticado")
	}

	input, ok := p.Args["input"].(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("input inválido")
	}

	courseID, _ := input["courseId"].(string)
	taskID, _ := input["taskId"].(string)
	submissionID, _ := input["submissionId"].(string)
	grade, _ := input["grade"].(float64)

	if courseID == "" || taskID == "" || submissionID == "" {
		return false, fmt.Errorf("courseId, taskId y submissionId son requeridos")
	}

	err := r.GoogleAdapter.GradeTask(p.Context, email, courseID, taskID, submissionID, grade)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ============================================
// FUNCIÓN AUXILIAR
// ============================================

func contains(s, substr string) bool {
	if s == "" || substr == "" {
		return false
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

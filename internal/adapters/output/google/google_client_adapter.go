package google

import (
	"context"
	"fmt"
	"log"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/classroom/v1"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	"github.com/tegnoword/orienmod/internal/core/domain"
	"github.com/tegnoword/orienmod/internal/core/ports"
)

type GoogleClientAdapter struct {
	oauthConfig   *oauth2.Config
	tokenStore    ports.TokenRepository
	spreadsheetID string
	limiter       *RateLimiter // ✅ NUEVO: Rate Limiting
}

func NewGoogleClientAdapter(config *oauth2.Config, store ports.TokenRepository, spreadsheetID string) *GoogleClientAdapter {
	return &GoogleClientAdapter{
		oauthConfig:   config,
		tokenStore:    store,
		spreadsheetID: spreadsheetID,
		limiter:       NewRateLimiter(), // ✅ INICIALIZAR RATE LIMITER
	}
}

// ✅ MÉTODO PARA CONFIGURAR RATE LIMITER DESDE FUERA (OPCIONAL)
func (a *GoogleClientAdapter) SetRateLimiter(limiter *RateLimiter) {
	a.limiter = limiter
}

func (a *GoogleClientAdapter) getClassroomService(ctx context.Context, clientEmail string) (*classroom.Service, error) {
	token, err := a.tokenStore.GetToken(ctx, clientEmail)
	if err != nil {
		return nil, fmt.Errorf("no se encontraron credenciales para el cliente %s: %w", clientEmail, err)
	}

	client := a.oauthConfig.Client(ctx, token)
	srv, err := classroom.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("error al inicializar servicio de Classroom: %w", err)
	}

	return srv, nil
}

func (a *GoogleClientAdapter) getSheetsService(ctx context.Context, clientEmail string) (*sheets.Service, error) {
	token, err := a.tokenStore.GetToken(ctx, clientEmail)
	if err != nil {
		return nil, fmt.Errorf("no se encontraron credenciales para el cliente %s: %w", clientEmail, err)
	}

	client := a.oauthConfig.Client(ctx, token)
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("error al inicializar servicio de Sheets: %w", err)
	}

	return srv, nil
}

// =========================================================================
// MÉTODOS PARA SHEETS
// =========================================================================

func (a *GoogleClientAdapter) SaveEvent(ctx context.Context, clientEmail string, event domain.ClassroomEvent) error {
	// ✅ RATE LIMITING
	if err := a.limiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limit: %w", err)
	}

	srv, err := a.getSheetsService(ctx, clientEmail)
	if err != nil {
		return err
	}

	values := [][]interface{}{
		{
			event.ResourceID,
			event.EventType,
			event.ReceivedAt.Format(time.RFC3339),
		},
	}

	valueRange := &sheets.ValueRange{Values: values}

	_, err = srv.Spreadsheets.Values.Append(a.spreadsheetID, "Eventos!A:C", valueRange).
		ValueInputOption("RAW").
		Do()

	if err != nil {
		return fmt.Errorf("error al guardar evento en Sheets: %w", err)
	}

	return nil
}

func (a *GoogleClientAdapter) SaveTasks(ctx context.Context, clientEmail string, studentID string, tasks []domain.Task) error {
	// ✅ RATE LIMITING
	if err := a.limiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limit: %w", err)
	}

	srv, err := a.getSheetsService(ctx, clientEmail)
	if err != nil {
		return err
	}

	var values [][]interface{}
	for _, t := range tasks {
		values = append(values, []interface{}{
			studentID,
			t.ID,
			t.Title,
			t.Description,
		})
	}

	valueRange := &sheets.ValueRange{Values: values}

	_, err = srv.Spreadsheets.Values.Append(a.spreadsheetID, "Tareas!A:D", valueRange).
		ValueInputOption("RAW").
		Do()

	if err != nil {
		return fmt.Errorf("error al registrar tareas en Sheets: %w", err)
	}

	return nil
}

func (a *GoogleClientAdapter) SaveGrade(ctx context.Context, clientEmail string, grade domain.Grade) error {
	// ✅ RATE LIMITING
	if err := a.limiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limit: %w", err)
	}

	srv, err := a.getSheetsService(ctx, clientEmail)
	if err != nil {
		return err
	}

	values := [][]interface{}{
		{grade.StudentID, grade.TaskID, grade.Score},
	}

	valueRange := &sheets.ValueRange{Values: values}

	_, err = srv.Spreadsheets.Values.Append(a.spreadsheetID, "Notas!A:C", valueRange).
		ValueInputOption("RAW").
		Do()

	if err != nil {
		return fmt.Errorf("error al guardar nota en Sheets: %w", err)
	}

	return nil
}

// =========================================================================
// MÉTODOS PARA CURSOS
// =========================================================================

func (a *GoogleClientAdapter) SaveCourse(ctx context.Context, clientEmail string, course domain.Course) error {
	// ✅ RATE LIMITING
	if err := a.limiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limit: %w", err)
	}

	srv, err := a.getClassroomService(ctx, clientEmail)
	if err != nil {
		return err
	}

	classroomCourse := &classroom.Course{
		Name:        course.Name,
		Section:     course.Section,
		Description: course.Description,
		OwnerId:     "me",
	}

	_, err = srv.Courses.Create(classroomCourse).Do()
	if err != nil {
		return fmt.Errorf("error al crear el curso en Classroom: %w", err)
	}

	return nil
}

func (a *GoogleClientAdapter) GetAllCourses(ctx context.Context, clientEmail string) ([]domain.Course, error) {
	// ✅ RATE LIMITING
	if err := a.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	srv, err := a.getClassroomService(ctx, clientEmail)
	if err != nil {
		return nil, err
	}

	res, err := srv.Courses.List().PageSize(20).Do()
	if err != nil {
		return nil, fmt.Errorf("error al listar cursos de Google Classroom: %w", err)
	}

	if res == nil || res.Courses == nil {
		return []domain.Course{}, nil
	}

	var courses []domain.Course
	for _, c := range res.Courses {
		if c == nil {
			continue
		}
		courses = append(courses, domain.Course{
			ID:          c.Id,
			Name:        c.Name,
			Section:     c.Section,
			Description: c.Description,
		})
	}

	return courses, nil
}

func (a *GoogleClientAdapter) UpdateCourse(ctx context.Context, email string, course domain.Course) error {
	// ✅ RATE LIMITING
	if err := a.limiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limit: %w", err)
	}

	srv, err := a.getClassroomService(ctx, email)
	if err != nil {
		return err
	}

	classroomCourse := &classroom.Course{
		Name:        course.Name,
		Section:     course.Section,
		Description: course.Description,
	}

	_, err = srv.Courses.Update(course.ID, classroomCourse).Do()
	if err != nil {
		return fmt.Errorf("error al actualizar curso: %w", err)
	}

	log.Printf("✅ Curso actualizado: %s", course.ID)
	return nil
}

func (a *GoogleClientAdapter) DeleteCourse(ctx context.Context, email string, courseID string) error {
	// ✅ RATE LIMITING
	if err := a.limiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limit: %w", err)
	}

	srv, err := a.getClassroomService(ctx, email)
	if err != nil {
		return err
	}

	course := &classroom.Course{
		CourseState: "DELETED",
	}
	_, err = srv.Courses.Update(courseID, course).Do()
	if err != nil {
		return fmt.Errorf("error al eliminar curso: %w", err)
	}

	log.Printf("🗑️ Curso eliminado: %s", courseID)
	return nil
}

// =========================================================================
// MÉTODOS PARA ESTUDIANTES
// =========================================================================

func (a *GoogleClientAdapter) SaveStudent(ctx context.Context, clientEmail string, courseID string, student domain.Student) error {
	// ✅ RATE LIMITING
	if err := a.limiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limit: %w", err)
	}

	srv, err := a.getClassroomService(ctx, clientEmail)
	if err != nil {
		return err
	}

	studentObj := &classroom.Student{
		UserId: student.ID,
	}

	_, err = srv.Courses.Students.Create(courseID, studentObj).Do()
	if err != nil {
		return fmt.Errorf("error al agregar estudiante a la clase: %w", err)
	}

	return nil
}

func (a *GoogleClientAdapter) GetStudentsByCourse(ctx context.Context, clientEmail string, courseID string) ([]domain.Student, error) {
	// ✅ RATE LIMITING
	if err := a.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	srv, err := a.getClassroomService(ctx, clientEmail)
	if err != nil {
		return nil, err
	}

	res, err := srv.Courses.Students.List(courseID).Do()
	if err != nil {
		return nil, fmt.Errorf("error al obtener estudiantes del curso %s: %w", courseID, err)
	}

	if res == nil || res.Students == nil {
		return []domain.Student{}, nil
	}

	var students []domain.Student
	for _, s := range res.Students {
		if s == nil {
			continue
		}

		fullName := ""
		if s.Profile != nil && s.Profile.Name != nil {
			fullName = s.Profile.Name.FullName
		}

		email := ""
		if s.Profile != nil {
			email = s.Profile.EmailAddress
		}

		students = append(students, domain.Student{
			ID:       s.UserId,
			Name:     fullName,
			Email:    email,
			CourseID: courseID,
		})
	}

	return students, nil
}

func (a *GoogleClientAdapter) DeleteStudent(ctx context.Context, email string, courseID string, studentID string) error {
	// ✅ RATE LIMITING
	if err := a.limiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limit: %w", err)
	}

	srv, err := a.getClassroomService(ctx, email)
	if err != nil {
		return err
	}

	_, err = srv.Courses.Students.Delete(courseID, studentID).Do()
	if err != nil {
		return fmt.Errorf("error al eliminar estudiante: %w", err)
	}

	log.Printf("🗑️ Estudiante eliminado: %s del curso %s", studentID, courseID)
	return nil
}

// =========================================================================
// MÉTODOS PARA TAREAS
// =========================================================================

func (a *GoogleClientAdapter) CreateTask(ctx context.Context, email string, task domain.Task) (string, error) {
	// ✅ RATE LIMITING
	if err := a.limiter.Wait(ctx); err != nil {
		return "", fmt.Errorf("rate limit: %w", err)
	}

	srv, err := a.getClassroomService(ctx, email)
	if err != nil {
		return "", err
	}

	if task.CourseID == "" {
		return "", fmt.Errorf("course_id es requerido")
	}
	if task.Title == "" {
		return "", fmt.Errorf("title es requerido")
	}

	validWorkTypes := map[string]bool{
		"ASSIGNMENT":               true,
		"SHORT_ANSWER":             true,
		"MULTIPLE_CHOICE_QUESTION": true,
	}

	workType := task.WorkType
	if workType == "" {
		workType = "ASSIGNMENT"
	} else if !validWorkTypes[workType] {
		log.Printf("⚠️ WorkType inválido: %s, usando ASSIGNMENT", workType)
		workType = "ASSIGNMENT"
	}

	classroomTask := &classroom.CourseWork{
		Title:       task.Title,
		Description: task.Description,
		MaxPoints:   task.MaxPoints,
		State:       "PUBLISHED",
		WorkType:    workType,
	}

	if !task.DueDate.IsZero() {
		classroomTask.DueDate = &classroom.Date{
			Year:  int64(task.DueDate.Year()),
			Month: int64(task.DueDate.Month()),
			Day:   int64(task.DueDate.Day()),
		}
		classroomTask.DueTime = &classroom.TimeOfDay{
			Hours:   int64(task.DueDate.Hour()),
			Minutes: int64(task.DueDate.Minute()),
		}
	}

	log.Printf("📝 Creando tarea: %s (WorkType: %s)", task.Title, workType)

	resp, err := srv.Courses.CourseWork.Create(task.CourseID, classroomTask).Do()
	if err != nil {
		return "", fmt.Errorf("error al crear tarea: %w", err)
	}

	log.Printf("✅ Tarea creada: %s", resp.Id)
	return resp.Id, nil
}

func (a *GoogleClientAdapter) GetTasksByCourse(ctx context.Context, email string, courseID string) ([]domain.Task, error) {
	// ✅ RATE LIMITING
	if err := a.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	srv, err := a.getClassroomService(ctx, email)
	if err != nil {
		return nil, err
	}

	resp, err := srv.Courses.CourseWork.List(courseID).Do()
	if err != nil {
		return nil, fmt.Errorf("error al listar tareas: %w", err)
	}

	if resp == nil || resp.CourseWork == nil {
		return []domain.Task{}, nil
	}

	var tasks []domain.Task
	for _, t := range resp.CourseWork {
		task := domain.Task{
			ID:          t.Id,
			CourseID:    courseID,
			Title:       t.Title,
			Description: t.Description,
			MaxPoints:   t.MaxPoints,
			State:       t.State,
			WorkType:    t.WorkType,
		}

		if t.DueDate != nil {
			task.DueDate = time.Date(
				int(t.DueDate.Year),
				time.Month(t.DueDate.Month),
				int(t.DueDate.Day),
				23, 59, 59, 0,
				time.UTC,
			)
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (a *GoogleClientAdapter) GetTaskByID(ctx context.Context, email string, courseID string, taskID string) (*domain.Task, error) {
	// ✅ RATE LIMITING
	if err := a.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	srv, err := a.getClassroomService(ctx, email)
	if err != nil {
		return nil, err
	}

	resp, err := srv.Courses.CourseWork.Get(courseID, taskID).Do()
	if err != nil {
		return nil, fmt.Errorf("error al obtener tarea: %w", err)
	}

	task := &domain.Task{
		ID:          resp.Id,
		CourseID:    courseID,
		Title:       resp.Title,
		Description: resp.Description,
		MaxPoints:   resp.MaxPoints,
		State:       resp.State,
		WorkType:    resp.WorkType,
	}

	if resp.DueDate != nil {
		task.DueDate = time.Date(
			int(resp.DueDate.Year),
			time.Month(resp.DueDate.Month),
			int(resp.DueDate.Day),
			23, 59, 59, 0,
			time.UTC,
		)
	}

	return task, nil
}

func (a *GoogleClientAdapter) UpdateTask(ctx context.Context, email string, task domain.Task) error {
	// ✅ RATE LIMITING
	if err := a.limiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limit: %w", err)
	}

	srv, err := a.getClassroomService(ctx, email)
	if err != nil {
		return err
	}

	validWorkTypes := map[string]bool{
		"ASSIGNMENT":               true,
		"SHORT_ANSWER":             true,
		"MULTIPLE_CHOICE_QUESTION": true,
	}

	workType := task.WorkType
	if workType == "" {
		workType = "ASSIGNMENT"
	} else if !validWorkTypes[workType] {
		log.Printf("⚠️ WorkType inválido: %s, usando ASSIGNMENT", workType)
		workType = "ASSIGNMENT"
	}

	classroomTask := &classroom.CourseWork{
		Title:       task.Title,
		Description: task.Description,
		MaxPoints:   task.MaxPoints,
		State:       task.State,
		WorkType:    workType,
	}

	if !task.DueDate.IsZero() {
		classroomTask.DueDate = &classroom.Date{
			Year:  int64(task.DueDate.Year()),
			Month: int64(task.DueDate.Month()),
			Day:   int64(task.DueDate.Day()),
		}
		classroomTask.DueTime = &classroom.TimeOfDay{
			Hours:   int64(task.DueDate.Hour()),
			Minutes: int64(task.DueDate.Minute()),
		}
	}

	_, err = srv.Courses.CourseWork.Patch(task.CourseID, task.ID, classroomTask).Do()
	if err != nil {
		return fmt.Errorf("error al actualizar tarea: %w", err)
	}

	log.Printf("✅ Tarea actualizada: %s", task.ID)
	return nil
}

func (a *GoogleClientAdapter) DeleteTask(ctx context.Context, email string, courseID string, taskID string) error {
	// ✅ RATE LIMITING
	if err := a.limiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limit: %w", err)
	}

	srv, err := a.getClassroomService(ctx, email)
	if err != nil {
		return err
	}

	_, err = srv.Courses.CourseWork.Delete(courseID, taskID).Do()
	if err != nil {
		return fmt.Errorf("error al eliminar tarea: %w", err)
	}

	log.Printf("🗑️ Tarea eliminada: %s", taskID)
	return nil
}

// =========================================================================
// MÉTODOS PARA ENTREGAS Y CALIFICACIONES
// =========================================================================

func (a *GoogleClientAdapter) GetTaskSubmissions(ctx context.Context, email string, courseID string, taskID string) ([]domain.TaskSubmission, error) {
	// ✅ RATE LIMITING
	if err := a.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	srv, err := a.getClassroomService(ctx, email)
	if err != nil {
		return nil, err
	}

	resp, err := srv.Courses.CourseWork.StudentSubmissions.List(courseID, taskID).Do()
	if err != nil {
		return nil, fmt.Errorf("error al listar entregas: %w", err)
	}

	if resp == nil || resp.StudentSubmissions == nil {
		return []domain.TaskSubmission{}, nil
	}

	var submissions []domain.TaskSubmission
	for _, s := range resp.StudentSubmissions {
		submission := domain.TaskSubmission{
			ID:        s.Id,
			CourseID:  courseID,
			TaskID:    taskID,
			StudentID: s.UserId,
			State:     s.State,
		}

		if s.AssignedGrade != 0 {
			submission.Grade = s.AssignedGrade
		}

		submissions = append(submissions, submission)
	}

	return submissions, nil
}

func (a *GoogleClientAdapter) GradeTask(ctx context.Context, email string, courseID string, taskID string, submissionID string, grade float64) error {
	// ✅ RATE LIMITING
	if err := a.limiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limit: %w", err)
	}

	srv, err := a.getClassroomService(ctx, email)
	if err != nil {
		return err
	}

	submission := &classroom.StudentSubmission{
		AssignedGrade: grade,
	}

	_, err = srv.Courses.CourseWork.StudentSubmissions.Patch(courseID, taskID, submissionID, submission).Do()
	if err != nil {
		return fmt.Errorf("error al calificar tarea: %w", err)
	}

	log.Printf("✅ Tarea calificada: %s con %.2f puntos", taskID, grade)
	return nil
}

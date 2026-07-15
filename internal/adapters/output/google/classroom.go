package google

import (
	"context"
	"fmt"
	"os"

	"github.com/tegnoword/orienmod/internal/core/domain"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/classroom/v1"
	"google.golang.org/api/option"
)

type ClassroomAdapter struct {
	service *classroom.Service
}

func NewClassroomAdapter(ctx context.Context, credentialsPath string) (*ClassroomAdapter, error) {
	b, err := os.ReadFile(credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("no se pudo leer el archivo de credenciales: %w", err)
	}

	config, err := google.JWTConfigFromJSON(b,
		classroom.ClassroomRostersReadonlyScope,
		classroom.ClassroomCoursesReadonlyScope,
		classroom.ClassroomCourseworkStudentsScope,
	)
	if err != nil {
		return nil, fmt.Errorf("error al parsear JWT desde JSON: %w", err)
	}

	client := config.Client(ctx)
	srv, err := classroom.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("no se pudo recuperar el servicio de Classroom: %w", err)
	}

	return &ClassroomAdapter{service: srv}, nil
}

func (a *ClassroomAdapter) SubmitGrade(ctx context.Context, courseID string, courseworkID string, studentSubmissionID string, grade float64) error {
	submission := &classroom.StudentSubmission{
		AssignedGrade: grade,
		DraftGrade:    grade,
	}

	_, err := a.service.Courses.CourseWork.StudentSubmissions.
		Patch(courseID, courseworkID, studentSubmissionID, submission).
		UpdateMask("assignedGrade,draftGrade").
		Context(ctx).
		Do()

	if err != nil {
		return fmt.Errorf("error al subir la nota a Google Classroom: %w", err)
	}

	return nil
}

func (a *ClassroomAdapter) GetCourses(ctx context.Context) ([]domain.Course, error) {
	res, err := a.service.Courses.List().
		CourseStates("ACTIVE").
		Context(ctx).
		Do()

	if err != nil {
		return nil, fmt.Errorf("error al listar cursos de Classroom: %w", err)
	}

	var courses []domain.Course
	for _, c := range res.Courses {
		courses = append(courses, domain.Course{
			ID:          c.Id,
			Name:        c.Name,
			Section:     c.Section,
			Description: c.Description,
		})
	}

	return courses, nil
}

func (a *ClassroomAdapter) GetStudentsByCourse(ctx context.Context, courseID string) ([]domain.Studen, error) {
	res, err := a.service.Courses.Students.List(courseID).
		Context(ctx).
		Do()

	if err != nil {
		return nil, fmt.Errorf("error al listar estudiantes del curso %s: %w", err)
	}

	var students []domain.Studen
	for _, s := range res.Students {

		students = append(students, domain.Studen{
			ID:    s.UserId,
			Name:  s.Profile.Name.FullName,
			Email: s.Profile.EmailAddress,
		})
	}

	return students, nil
}

// internal/graphql/server.go
package graphql

import (
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/tegnoword/orienmod/internal/adapters/output/google"
	"github.com/tegnoword/orienmod/internal/core/ports"
	"golang.org/x/oauth2"
)

func NewGraphQLHandler(
	googleAdapter *google.GoogleClientAdapter,
	tokenStore ports.TokenRepository,
	oauthConfig *oauth2.Config,
) http.Handler {
	resolver := NewResolver(googleAdapter, tokenStore, oauthConfig)
	rootQuery := graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			// ✅ 1. checkAuth - Con email
			"checkAuth": &graphql.Field{
				Type: AuthCheckResponseType,
				Args: graphql.FieldConfigArgument{
					"email": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolver.CheckAuthResolver(p)
				},
			},

			// ✅ 2. courses - AHORA CON EMAIL
			"courses": &graphql.Field{
				Type: graphql.NewList(CourseType),
				Args: graphql.FieldConfigArgument{
					"email": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolver.CoursesResolver(p)
				},
			},

			// ✅ 3. course - Con id y email
			"course": &graphql.Field{
				Type: CourseType,
				Args: graphql.FieldConfigArgument{
					"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"email": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolver.CourseResolver(p)
				},
			},

			// ✅ 4. searchCourses - Con query y email
			"searchCourses": &graphql.Field{
				Type: graphql.NewList(CourseType),
				Args: graphql.FieldConfigArgument{
					"query": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"email": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolver.SearchCoursesResolver(p)
				},
			},

			// ✅ 5. students - Con courseId y email
			"students": &graphql.Field{
				Type: graphql.NewList(StudentType),
				Args: graphql.FieldConfigArgument{
					"courseId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"email":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolver.StudentsResolver(p)
				},
			},

			// ✅ 6. searchStudents - Con courseId, query y email
			"searchStudents": &graphql.Field{
				Type: graphql.NewList(StudentType),
				Args: graphql.FieldConfigArgument{
					"courseId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"query":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"email":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolver.SearchStudentsResolver(p)
				},
			},

			// ✅ 7. tasks - Con courseId y email
			"tasks": &graphql.Field{
				Type: graphql.NewList(TaskType),
				Args: graphql.FieldConfigArgument{
					"courseId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"email":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolver.TasksResolver(p)
				},
			},

			// ✅ 8. task - Con id, courseId y email
			"task": &graphql.Field{
				Type: TaskType,
				Args: graphql.FieldConfigArgument{
					"id":       &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"courseId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"email":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolver.TaskResolver(p)
				},
			},

			// ✅ 9. taskSubmissions - Con taskId, courseId y email
			"taskSubmissions": &graphql.Field{
				Type: graphql.NewList(TaskSubmissionType),
				Args: graphql.FieldConfigArgument{
					"taskId":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"courseId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"email":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolver.TaskSubmissionsResolver(p)
				},
			},
		},
	}

	// ============================================
	// MUTATIONS
	// ============================================

	rootMutation := graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"createCourse": &graphql.Field{
				Type: CourseType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(CreateCourseInputType)},
					"email": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: resolver.CreateCourseResolver,
			},
			"updateCourse": &graphql.Field{
				Type: CourseType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(UpdateCourseInputType)},
					"email": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: resolver.UpdateCourseResolver,
			},
			"deleteCourse": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"email": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: resolver.DeleteCourseResolver,
			},
			"syncCourse": &graphql.Field{
				Type: graphql.Int,
				Args: graphql.FieldConfigArgument{
					"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"email": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: resolver.SyncCourseResolver,
			},
			"addStudent": &graphql.Field{
				Type: StudentType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(AddStudentInputType)},
					"email": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: resolver.AddStudentResolver,
			},
			"deleteStudent": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"courseId":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"studentId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"email":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: resolver.DeleteStudentResolver,
			},
			"createTask": &graphql.Field{
				Type: TaskType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(CreateTaskInputType)},
					"email": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: resolver.CreateTaskResolver,
			},
			"updateTask": &graphql.Field{
				Type: TaskType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(UpdateTaskInputType)},
					"email": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: resolver.UpdateTaskResolver,
			},
			"deleteTask": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"id":       &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"courseId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"email":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: resolver.DeleteTaskResolver,
			},
			"gradeTask": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(GradeTaskInputType)},
					"email": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: resolver.GradeTaskResolver,
			},
		},
	}

	schemaConfig := graphql.SchemaConfig{
		Query:    graphql.NewObject(rootQuery),
		Mutation: graphql.NewObject(rootMutation),
	}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		panic(err)
	}

	return handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		GraphiQL:   true,
		Playground: true,
	})
}

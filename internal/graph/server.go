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
			"checkAuth": &graphql.Field{
				Type: AuthCheckResponseType,
				Args: graphql.FieldConfigArgument{
					"email": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolver.CheckAuthResolver(p)
				},
			},
			"courses": &graphql.Field{
				Type: graphql.NewList(CourseType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolver.CoursesResolver(p)
				},
			},
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
	// CONFIGURAR MUTATIONS
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
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolver.CreateCourseResolver(p)
				},
			},
			"updateCourse": &graphql.Field{
				Type: CourseType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(UpdateCourseInputType)},
					"email": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolver.UpdateCourseResolver(p)
				},
			},
			"deleteCourse": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"email": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolver.DeleteCourseResolver(p)
				},
			},
			"syncCourse": &graphql.Field{
				Type: graphql.Int,
				Args: graphql.FieldConfigArgument{
					"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"email": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolver.SyncCourseResolver(p)
				},
			},
			"addStudent": &graphql.Field{
				Type: StudentType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(AddStudentInputType)},
					"email": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolver.AddStudentResolver(p)
				},
			},
			"deleteStudent": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"courseId":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"studentId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"email":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolver.DeleteStudentResolver(p)
				},
			},
			"createTask": &graphql.Field{
				Type: TaskType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(CreateTaskInputType)},
					"email": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolver.CreateTaskResolver(p)
				},
			},
			"updateTask": &graphql.Field{
				Type: TaskType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(UpdateTaskInputType)},
					"email": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolver.UpdateTaskResolver(p)
				},
			},
			"deleteTask": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"id":       &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"courseId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"email":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolver.DeleteTaskResolver(p)
				},
			},
			"gradeTask": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(GradeTaskInputType)},
					"email": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolver.GradeTaskResolver(p)
				},
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
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})
}

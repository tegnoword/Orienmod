package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/tegnoword/orienmod/internal/adapters/output/google"
	"github.com/tegnoword/orienmod/internal/core/domain"
	"github.com/tegnoword/orienmod/internal/core/ports"
	"golang.org/x/oauth2"
)

type Resolver struct {
	GoogleAdapter *google.GoogleClientAdapter
	TokenStore    ports.TokenRepository
	OAuthConfig   *oauth2.Config
}

func NewResolver(
	googleAdapter *google.GoogleClientAdapter,
	tokenStore ports.TokenRepository,
	oauthConfig *oauth2.Config,
) *Resolver {
	return &Resolver{
		GoogleAdapter: googleAdapter,
		TokenStore:    tokenStore,
		OAuthConfig:   oauthConfig,
	}
}

var CourseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Course",
	Fields: graphql.Fields{
		"id":          &graphql.Field{Type: graphql.String},
		"name":        &graphql.Field{Type: graphql.String},
		"section":     &graphql.Field{Type: graphql.String},
		"description": &graphql.Field{Type: graphql.String},
		"ownerId":     &graphql.Field{Type: graphql.String},
		"students": &graphql.Field{
			Type: graphql.NewList(StudentType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				course, ok := p.Source.(domain.Course)
				if !ok {
					return nil, nil
				}
				email, _ := p.Context.Value("email").(string)
				if email == "" {
					return nil, nil
				}
				resolver, _ := p.Context.Value("resolver").(*Resolver)
				if resolver == nil {
					return nil, nil
				}
				return resolver.GoogleAdapter.GetStudentsByCourse(p.Context, email, course.ID)
			},
		},
		"tasks": &graphql.Field{
			Type: graphql.NewList(TaskType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				course, ok := p.Source.(domain.Course)
				if !ok {
					return nil, nil
				}
				email, _ := p.Context.Value("email").(string)
				if email == "" {
					return nil, nil
				}
				resolver, _ := p.Context.Value("resolver").(*Resolver)
				if resolver == nil {
					return nil, nil
				}
				return resolver.GoogleAdapter.GetTasksByCourse(p.Context, email, course.ID)
			},
		},
	},
})

var StudentType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Student",
	Fields: graphql.Fields{
		"id":       &graphql.Field{Type: graphql.String},
		"name":     &graphql.Field{Type: graphql.String},
		"email":    &graphql.Field{Type: graphql.String},
		"courseId": &graphql.Field{Type: graphql.String},
	},
})

var TaskType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Task",
	Fields: graphql.Fields{
		"id":          &graphql.Field{Type: graphql.String},
		"courseId":    &graphql.Field{Type: graphql.String},
		"title":       &graphql.Field{Type: graphql.String},
		"description": &graphql.Field{Type: graphql.String},
		"dueDate":     &graphql.Field{Type: graphql.String},
		"maxPoints":   &graphql.Field{Type: graphql.Float},
		"state":       &graphql.Field{Type: graphql.String},
		"workType":    &graphql.Field{Type: graphql.String},
		"submissions": &graphql.Field{
			Type: graphql.NewList(TaskSubmissionType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				task, ok := p.Source.(domain.Task)
				if !ok {
					return nil, nil
				}
				email, _ := p.Context.Value("email").(string)
				if email == "" {
					return nil, nil
				}
				resolver, _ := p.Context.Value("resolver").(*Resolver)
				if resolver == nil {
					return nil, nil
				}
				return resolver.GoogleAdapter.GetTaskSubmissions(p.Context, email, task.CourseID, task.ID)
			},
		},
	},
})

var TaskSubmissionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "TaskSubmission",
	Fields: graphql.Fields{
		"id":        &graphql.Field{Type: graphql.String},
		"taskId":    &graphql.Field{Type: graphql.String},
		"studentId": &graphql.Field{Type: graphql.String},
		"state":     &graphql.Field{Type: graphql.String},
		"grade":     &graphql.Field{Type: graphql.Float},
	},
})

var AuthResponseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "AuthResponse",
	Fields: graphql.Fields{
		"status":  &graphql.Field{Type: graphql.String},
		"message": &graphql.Field{Type: graphql.String},
		"email":   &graphql.Field{Type: graphql.String},
		"token":   &graphql.Field{Type: TokenInfoType},
	},
})

var TokenInfoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "TokenInfo",
	Fields: graphql.Fields{
		"accessToken": &graphql.Field{Type: graphql.String},
		"tokenType":   &graphql.Field{Type: graphql.String},
		"expiry":      &graphql.Field{Type: graphql.String},
	},
})

var AuthCheckResponseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "AuthCheckResponse",
	Fields: graphql.Fields{
		"authenticated": &graphql.Field{Type: graphql.Boolean},
		"email":         &graphql.Field{Type: graphql.String},
		"tokenExpiry":   &graphql.Field{Type: graphql.String},
		"valid":         &graphql.Field{Type: graphql.Boolean},
		"error":         &graphql.Field{Type: graphql.String},
	},
})

var CreateCourseInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "CreateCourseInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"name":        &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"section":     &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"description": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
	},
})

var UpdateCourseInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "UpdateCourseInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"id":          &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"name":        &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"section":     &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"description": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
	},
})

var CreateTaskInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "CreateTaskInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"courseId":    &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"title":       &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"description": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"dueDate":     &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"maxPoints":   &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Float)},
		"workType":    &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
	},
})

var UpdateTaskInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "UpdateTaskInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"id":          &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"courseId":    &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"title":       &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"description": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"dueDate":     &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"maxPoints":   &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Float)},
		"state":       &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"workType":    &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
	},
})

var AddStudentInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "AddStudentInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"courseId":  &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"studentId": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"email":     &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"name":      &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
	},
})

var GradeTaskInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "GradeTaskInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"courseId":     &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"taskId":       &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"submissionId": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"grade":        &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Float)},
	},
})

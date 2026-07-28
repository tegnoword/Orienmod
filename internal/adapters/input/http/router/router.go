package router

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"github.com/tegnoword/orienmod/internal/adapters/output/google"
	"github.com/tegnoword/orienmod/internal/core/domain"
	apperrors "github.com/tegnoword/orienmod/internal/shared/error"
)

// ============================================
// INTERFACES Y ESTRUCTURAS
// ============================================

type TokenStore interface {
	SaveToken(ctx context.Context, email string, token *oauth2.Token) error
	GetToken(ctx context.Context, email string) (*oauth2.Token, error)
	DeleteToken(ctx context.Context, email string) error
}

type Router struct {
	oauthConfig   *oauth2.Config
	tokenStore    TokenStore
	googleAdapter *google.GoogleClientAdapter
}

func NewRouter(
	oauthConfig *oauth2.Config,
	tokenStore TokenStore,
	googleAdapter *google.GoogleClientAdapter,
) *Router {
	return &Router{
		oauthConfig:   oauthConfig,
		tokenStore:    tokenStore,
		googleAdapter: googleAdapter,
	}
}

// ============================================
// SERVIDOR PRINCIPAL - SERVEHTTP
// ============================================

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Configurar CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-User-Email, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, QUERY, OPTIONS")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	email := r.Header.Get("X-User-Email")
	log.Printf("📨 Petición: %s %s | Email: %s", r.Method, r.URL.Path, email)

	switch {
	// HEALTH CHECK
	case r.URL.Path == "/health" && r.Method == http.MethodGet:
		rt.handleHealthCheck(w, r)

	// AUTENTICACIÓN
	case r.URL.Path == "/api/v1/auth/google" && r.Method == http.MethodGet:
		rt.handleLogin(w, r)

	case r.URL.Path == "/api/v1/auth/google/callback" && r.Method == http.MethodGet:
		rt.handleCallback(w, r)

	case r.URL.Path == "/api/v1/auth/logout" && r.Method == http.MethodPost:
		rt.requireAuth(w, r, email, rt.handleLogout)

	case r.URL.Path == "/api/v1/auth/check" && r.Method == http.MethodGet:
		rt.handleCheckAuth(w, r)

	case r.URL.Path == "/api/v1/auth/refresh" && r.Method == http.MethodPost:
		rt.requireAuth(w, r, email, rt.handleRefreshToken)

	case r.URL.Path == "/api/v1/auth/force" && r.Method == http.MethodGet:
		rt.handleForceAuth(w, r)

	// CLASSROOM - CURSOS
	case r.URL.Path == "/api/v1/courses" && r.Method == http.MethodGet:
		rt.requireAuth(w, r, email, rt.handleGetCourses)

	case r.URL.Path == "/api/v1/courses" && r.Method == http.MethodPost:
		rt.requireAuth(w, r, email, rt.handleCreateCourse)

	case r.URL.Path == "/api/v1/courses/search" && r.Method == http.MethodGet:
		rt.requireAuth(w, r, email, rt.handleSearchCourses)

	case strings.HasPrefix(r.URL.Path, "/api/v1/courses/") && r.Method == http.MethodPut:
		rt.requireAuth(w, r, email, rt.handleUpdateCourse)

	case strings.HasPrefix(r.URL.Path, "/api/v1/courses/") && r.Method == http.MethodDelete:
		rt.requireAuth(w, r, email, rt.handleDeleteCourse)

	case strings.HasPrefix(r.URL.Path, "/api/v1/courses/") && r.Method == http.MethodGet:
		if !strings.Contains(r.URL.Path[15:], "/") {
			rt.requireAuth(w, r, email, rt.handleGetCourseByID)
		} else {
			rt.requireAuth(w, r, email, rt.handleCourseRoutes)
		}

	// CLASSROOM - ESTUDIANTES
	case r.URL.Path == "/api/v1/students/search" && r.Method == http.MethodGet:
		rt.requireAuth(w, r, email, rt.handleSearchStudents)

	case strings.HasPrefix(r.URL.Path, "/api/v1/courses/") && strings.Contains(r.URL.Path, "/students/") && r.Method == http.MethodDelete:
		rt.requireAuth(w, r, email, rt.handleDeleteStudent)

	// CLASSROOM - TAREAS
	case r.URL.Path == "/api/v1/tasks" && r.Method == http.MethodPost:
		rt.requireAuth(w, r, email, rt.handleCreateTask)

	case r.URL.Path == "/api/v1/tasks" && r.Method == http.MethodGet:
		rt.requireAuth(w, r, email, rt.handleGetTasks)

	case strings.HasPrefix(r.URL.Path, "/api/v1/tasks/") && r.Method == http.MethodPut:
		rt.requireAuth(w, r, email, rt.handleUpdateTask)

	case strings.HasPrefix(r.URL.Path, "/api/v1/tasks/") && r.Method == http.MethodDelete:
		rt.requireAuth(w, r, email, rt.handleDeleteTask)

	case strings.HasPrefix(r.URL.Path, "/api/v1/tasks/") && r.Method == http.MethodGet:
		if !strings.Contains(r.URL.Path[14:], "/") {
			rt.requireAuth(w, r, email, rt.handleGetTaskByID)
		} else {
			rt.requireAuth(w, r, email, rt.handleTaskRoutes)
		}

	// WEBHOOK
	case r.URL.Path == "/api/v1/classroom/webhook" && r.Method == http.MethodPost:
		rt.handleWebhook(w, r)

	case r.URL.Path == "/api/v1/classroom/webhook" && r.Method == http.MethodGet:
		rt.handleWebhookVerification(w, r)

	default:
		log.Printf("❌ Ruta no encontrada: %s %s", r.Method, r.URL.Path)
		http.NotFound(w, r)
	}
}

// ============================================
// FUNCIONES DE RESPUESTA CON ERRORES GLOBALES
// ============================================

// respondJSON - Envía una respuesta JSON
func (rt *Router) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("❌ Error al codificar JSON: %v", err)
	}
}

// respondError - Envía una respuesta de error usando AppError
func (rt *Router) respondError(w http.ResponseWriter, err *apperrors.AppError) {
	response := err.ToResponse()
	log.Printf("❌ Error: %s (Code: %d)", err.Message, err.Code)
	rt.respondJSON(w, err.Code, response)
}

// respondHTTPError - Envía un error HTTP simple (para casos sin AppError)
func (rt *Router) respondHTTPError(w http.ResponseWriter, status int, message string) {
	rt.respondJSON(w, status, map[string]string{"error": message})
}

// ============================================
// HANDLERS DE AUTENTICACIÓN
// ============================================

func (rt *Router) handleLogin(w http.ResponseWriter, r *http.Request) {
	if rt.oauthConfig.ClientID == "" {
		rt.respondError(w, apperrors.NewBadRequest("OAuth2 no configurado correctamente"))
		return
	}

	redirectURL := rt.oauthConfig.RedirectURL
	if redirectURL == "" {
		redirectURL = "http://localhost:8080/api/v1/auth/google/callback"
	}

	authURL := rt.oauthConfig.AuthCodeURL(
		"state-token-"+time.Now().Format("20060102150405"),
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce,
		oauth2.SetAuthURLParam("redirect_uri", redirectURL),
	)

	log.Printf("🌐 Redirigiendo a Google")
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (rt *Router) handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		rt.respondError(w, apperrors.NewBadRequest("Código de autorización no encontrado"))
		return
	}

	redirectURL := rt.oauthConfig.RedirectURL
	if redirectURL == "" {
		redirectURL = "http://localhost:8080/api/v1/auth/google/callback"
	}

	token, err := rt.oauthConfig.Exchange(
		r.Context(),
		code,
		oauth2.SetAuthURLParam("redirect_uri", redirectURL),
	)
	if err != nil {
		log.Printf("❌ Error al intercambiar token: %v", err)
		rt.respondError(w, apperrors.NewGoogleAPIError(err))
		return
	}

	email, err := rt.getUserEmail(r.Context(), token)
	if err != nil {
		log.Printf("❌ Error al obtener email: %v", err)
		rt.respondError(w, apperrors.NewInternalError(err))
		return
	}

	if err := rt.tokenStore.SaveToken(r.Context(), email, token); err != nil {
		log.Printf("❌ Error al guardar token para %s: %v", email, err)
		rt.respondError(w, apperrors.NewInternalError(err))
		return
	}

	log.Printf("✅ Autenticación exitosa | Email: %s", email)

	rt.respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Autenticación exitosa",
		"email":   email,
		"token": map[string]interface{}{
			"access_token": token.AccessToken,
			"token_type":   token.TokenType,
			"expiry":       token.Expiry.Format(time.RFC3339),
		},
	})
}

func (rt *Router) handleLogout(w http.ResponseWriter, r *http.Request, email string) {
	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rt.respondError(w, apperrors.NewBadRequestWithDetails("Error al parsear body", err.Error()))
		return
	}

	if req.Email == "" {
		req.Email = email
	}

	if req.Email == "" {
		rt.respondError(w, apperrors.NewBadRequest("Email es requerido"))
		return
	}

	if err := rt.tokenStore.DeleteToken(r.Context(), req.Email); err != nil {
		log.Printf("❌ Error al eliminar token para %s: %v", req.Email, err)
		rt.respondError(w, apperrors.NewInternalError(err))
		return
	}

	log.Printf("✅ Sesión cerrada | Email: %s", req.Email)
	rt.respondJSON(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Sesión cerrada correctamente",
	})
}

func (rt *Router) handleCheckAuth(w http.ResponseWriter, r *http.Request) {
	email := r.Header.Get("X-User-Email")
	if email == "" {
		rt.respondError(w, apperrors.NewBadRequest("X-User-Email requerido"))
		return
	}

	token, err := rt.tokenStore.GetToken(r.Context(), email)
	if err != nil {
		rt.respondJSON(w, http.StatusOK, map[string]interface{}{
			"authenticated": false,
			"email":         email,
			"error":         err.Error(),
		})
		return
	}

	rt.respondJSON(w, http.StatusOK, map[string]interface{}{
		"authenticated": true,
		"email":         email,
		"token_expiry":  token.Expiry.Format(time.RFC3339),
		"valid":         token.Valid(),
	})
}

func (rt *Router) handleRefreshToken(w http.ResponseWriter, r *http.Request, email string) {
	token, err := rt.tokenStore.GetToken(r.Context(), email)
	if err != nil {
		rt.respondError(w, apperrors.NewUnauthorized("Token no encontrado"))
		return
	}

	if token.Valid() {
		rt.respondJSON(w, http.StatusOK, map[string]interface{}{
			"status":  "success",
			"message": "Token aún válido",
			"expiry":  token.Expiry.Format(time.RFC3339),
		})
		return
	}

	newToken, err := rt.oauthConfig.TokenSource(r.Context(), token).Token()
	if err != nil {
		log.Printf("❌ Error al refrescar token: %v", err)
		rt.respondError(w, apperrors.NewGoogleAPIError(err))
		return
	}

	if err := rt.tokenStore.SaveToken(r.Context(), email, newToken); err != nil {
		log.Printf("❌ Error al guardar nuevo token: %v", err)
		rt.respondError(w, apperrors.NewInternalError(err))
		return
	}

	rt.respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":       "success",
		"message":      "Token refrescado",
		"access_token": newToken.AccessToken,
		"expiry":       newToken.Expiry.Format(time.RFC3339),
	})
}

func (rt *Router) handleForceAuth(w http.ResponseWriter, r *http.Request) {
	email := r.Header.Get("X-User-Email")
	if email != "" {
		if err := rt.tokenStore.DeleteToken(r.Context(), email); err != nil {
			log.Printf("⚠️ Error al eliminar token: %v", err)
		} else {
			log.Printf("🔄 Token eliminado para: %s", email)
		}
	}

	redirectURL := rt.oauthConfig.RedirectURL
	if redirectURL == "" {
		redirectURL = "http://localhost:8080/api/v1/auth/google/callback"
	}

	authURL := rt.oauthConfig.AuthCodeURL(
		"force-auth-"+time.Now().Format("20060102150405"),
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce,
		oauth2.SetAuthURLParam("redirect_uri", redirectURL),
	)

	log.Printf("🔄 Redirigiendo a Google para forzar re-autenticación")
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// ============================================
// HANDLERS DE CLASSROOM - CURSOS
// ============================================

func (rt *Router) handleGetCourses(w http.ResponseWriter, r *http.Request, email string) {
	log.Printf("📚 Obteniendo cursos | Email: %s", email)

	courses, err := rt.googleAdapter.GetAllCourses(r.Context(), email)
	if err != nil {
		log.Printf("❌ Error al obtener cursos para %s: %v", email, err)
		rt.respondError(w, apperrors.NewGoogleAPIError(err))
		return
	}

	rt.respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"courses": courses,
		"count":   len(courses),
	})
}

func (rt *Router) handleGetCourseByID(w http.ResponseWriter, r *http.Request, email string) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/courses/")
	courseID := strings.Split(path, "/")[0]

	log.Printf("📚 Obteniendo curso | Email: %s | ID: %s", email, courseID)

	courses, err := rt.googleAdapter.GetAllCourses(r.Context(), email)
	if err != nil {
		rt.respondError(w, apperrors.NewGoogleAPIError(err))
		return
	}

	for _, c := range courses {
		if c.ID == courseID {
			rt.respondJSON(w, http.StatusOK, map[string]interface{}{
				"status": "success",
				"course": c,
			})
			return
		}
	}

	rt.respondError(w, apperrors.NewNotFound("Curso no encontrado"))
}

func (rt *Router) handleCreateCourse(w http.ResponseWriter, r *http.Request, email string) {
	var req struct {
		Name        string `json:"name"`
		Section     string `json:"section"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rt.respondError(w, apperrors.NewBadRequestWithDetails("Error al parsear body", err.Error()))
		return
	}

	if req.Name == "" {
		rt.respondError(w, apperrors.NewValidationError("name", "es requerido"))
		return
	}

	log.Printf("📚 Creando curso | Email: %s | Nombre: %s", email, req.Name)

	course := domain.Course{
		Name:        req.Name,
		Section:     req.Section,
		Description: req.Description,
	}

	if err := rt.googleAdapter.SaveCourse(r.Context(), email, course); err != nil {
		log.Printf("❌ Error al crear curso para %s: %v", email, err)
		rt.respondError(w, apperrors.NewGoogleAPIError(err))
		return
	}

	rt.respondJSON(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Curso creado exitosamente",
	})
}

func (rt *Router) handleUpdateCourse(w http.ResponseWriter, r *http.Request, email string) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/courses/")
	courseID := strings.Split(path, "/")[0]

	var req struct {
		Name        string `json:"name"`
		Section     string `json:"section"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rt.respondError(w, apperrors.NewBadRequestWithDetails("Error al parsear body", err.Error()))
		return
	}

	if req.Name == "" {
		rt.respondError(w, apperrors.NewValidationError("name", "es requerido"))
		return
	}

	course := domain.Course{
		ID:          courseID,
		Name:        req.Name,
		Section:     req.Section,
		Description: req.Description,
	}

	if err := rt.googleAdapter.UpdateCourse(r.Context(), email, course); err != nil {
		log.Printf("❌ Error al actualizar curso: %v", err)
		rt.respondError(w, apperrors.NewGoogleAPIError(err))
		return
	}

	rt.respondJSON(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Curso actualizado exitosamente",
	})
}

func (rt *Router) handleDeleteCourse(w http.ResponseWriter, r *http.Request, email string) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/courses/")
	courseID := strings.Split(path, "/")[0]

	log.Printf("🗑️ Eliminando curso | Email: %s | ID: %s", email, courseID)

	if err := rt.googleAdapter.DeleteCourse(r.Context(), email, courseID); err != nil {
		log.Printf("❌ Error al eliminar curso: %v", err)
		rt.respondError(w, apperrors.NewGoogleAPIError(err))
		return
	}

	rt.respondJSON(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Curso eliminado exitosamente",
	})
}

func (rt *Router) handleSearchCourses(w http.ResponseWriter, r *http.Request, email string) {
	var req struct {
		Query string `json:"query"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rt.respondError(w, apperrors.NewBadRequestWithDetails("Error al parsear body", err.Error()))
		return
	}

	log.Printf("🔍 Buscando cursos | Email: %s | Query: %s", email, req.Query)

	courses, err := rt.googleAdapter.GetAllCourses(r.Context(), email)
	if err != nil {
		log.Printf("❌ Error al buscar cursos para %s: %v", email, err)
		rt.respondError(w, apperrors.NewGoogleAPIError(err))
		return
	}

	var results []domain.Course
	for _, c := range courses {
		if req.Query == "" ||
			strings.Contains(strings.ToLower(c.Name), strings.ToLower(req.Query)) ||
			strings.Contains(strings.ToLower(c.Description), strings.ToLower(req.Query)) {
			results = append(results, c)
		}
	}

	rt.respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"results": results,
		"count":   len(results),
	})
}

// ============================================
// HANDLERS DE CLASSROOM - ESTUDIANTES
// ============================================

func (rt *Router) handleSearchStudents(w http.ResponseWriter, r *http.Request, email string) {
	var req struct {
		CourseID string `json:"course_id"`
		Query    string `json:"query"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rt.respondError(w, apperrors.NewBadRequestWithDetails("Error al parsear body", err.Error()))
		return
	}

	if req.CourseID == "" {
		rt.respondError(w, apperrors.NewValidationError("course_id", "es requerido"))
		return
	}

	log.Printf("🔍 Buscando estudiantes | Email: %s | Course: %s | Query: %s",
		email, req.CourseID, req.Query)

	students, err := rt.googleAdapter.GetStudentsByCourse(r.Context(), email, req.CourseID)
	if err != nil {
		log.Printf("❌ Error al obtener estudiantes para %s: %v", email, err)
		rt.respondError(w, apperrors.NewGoogleAPIError(err))
		return
	}

	var results []domain.Student
	for _, s := range students {
		if req.Query == "" ||
			strings.Contains(strings.ToLower(s.Name), strings.ToLower(req.Query)) ||
			strings.Contains(strings.ToLower(s.Email), strings.ToLower(req.Query)) {
			results = append(results, s)
		}
	}

	rt.respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"results": results,
		"count":   len(results),
	})
}

func (rt *Router) handleDeleteStudent(w http.ResponseWriter, r *http.Request, email string) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/courses/")
	parts := strings.Split(path, "/")

	if len(parts) < 3 || parts[1] != "students" {
		http.NotFound(w, r)
		return
	}

	courseID := parts[0]
	studentID := parts[2]

	log.Printf("🗑️ Eliminando estudiante | Email: %s | Curso: %s | Estudiante: %s",
		email, courseID, studentID)

	if err := rt.googleAdapter.DeleteStudent(r.Context(), email, courseID, studentID); err != nil {
		log.Printf("❌ Error al eliminar estudiante: %v", err)
		rt.respondError(w, apperrors.NewGoogleAPIError(err))
		return
	}

	rt.respondJSON(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Estudiante eliminado exitosamente",
	})
}

// ============================================
// HANDLERS DE CLASSROOM - TAREAS
// ============================================

func (rt *Router) handleCreateTask(w http.ResponseWriter, r *http.Request, email string) {
	var req struct {
		CourseID    string  `json:"course_id"`
		Title       string  `json:"title"`
		Description string  `json:"description"`
		DueDate     string  `json:"due_date"`
		MaxPoints   float64 `json:"max_points"`
		WorkType    string  `json:"work_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rt.respondError(w, apperrors.NewBadRequestWithDetails("Error al parsear body", err.Error()))
		return
	}

	if req.CourseID == "" {
		rt.respondError(w, apperrors.NewValidationError("course_id", "es requerido"))
		return
	}
	if req.Title == "" {
		rt.respondError(w, apperrors.NewValidationError("title", "es requerido"))
		return
	}
	if len(req.Title) > 250 {
		rt.respondError(w, apperrors.NewValidationError("title", "no puede exceder 250 caracteres"))
		return
	}
	if req.MaxPoints < 0 {
		rt.respondError(w, apperrors.NewValidationError("max_points", "no puede ser negativo"))
		return
	}

	var dueDate time.Time
	var err error
	if req.DueDate != "" {
		dueDate, err = time.Parse(time.RFC3339, req.DueDate)
		if err != nil {
			rt.respondError(w, apperrors.NewBadRequestWithDetails("Formato de fecha inválido", "Usa ISO 8601 (ej: 2026-07-30T23:59:59Z)"))
			return
		}
		if dueDate.Before(time.Now()) {
			rt.respondError(w, apperrors.NewValidationError("due_date", "no puede ser en el pasado"))
			return
		}
	}

	workType := req.WorkType
	if workType == "" {
		workType = "ASSIGNMENT"
	}
	validWorkTypes := map[string]bool{
		"ASSIGNMENT": true, "SHORT_ANSWER": true, "MULTIPLE_CHOICE_QUESTION": true,
	}
	if !validWorkTypes[workType] {
		rt.respondError(w, apperrors.NewValidationError("work_type", "inválido: "+workType))
		return
	}

	log.Printf("📝 Creando tarea | Email: %s | Curso: %s | Título: %s | WorkType: %s",
		email, req.CourseID, req.Title, workType)

	task := domain.Task{
		CourseID:    req.CourseID,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     dueDate,
		MaxPoints:   req.MaxPoints,
		WorkType:    workType,
	}

	taskID, err := rt.googleAdapter.CreateTask(r.Context(), email, task)
	if err != nil {
		log.Printf("❌ Error al crear tarea: %v", err)
		rt.respondError(w, apperrors.NewGoogleAPIError(err))
		return
	}

	log.Printf("✅ Tarea creada | ID: %s", taskID)

	rt.respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Tarea creada exitosamente",
		"task_id": taskID,
	})
}

func (rt *Router) handleGetTasks(w http.ResponseWriter, r *http.Request, email string) {
	courseID := r.URL.Query().Get("course_id")
	if courseID == "" {
		rt.respondError(w, apperrors.NewBadRequest("course_id es requerido"))
		return
	}

	log.Printf("📚 Obteniendo tareas | Email: %s | Curso: %s", email, courseID)

	tasks, err := rt.googleAdapter.GetTasksByCourse(r.Context(), email, courseID)
	if err != nil {
		log.Printf("❌ Error al obtener tareas: %v", err)
		rt.respondError(w, apperrors.NewGoogleAPIError(err))
		return
	}

	rt.respondJSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"tasks":  tasks,
		"count":  len(tasks),
	})
}

func (rt *Router) handleGetTaskByID(w http.ResponseWriter, r *http.Request, email string) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/tasks/")
	taskID := strings.Split(path, "/")[0]

	courseID := r.URL.Query().Get("course_id")
	if courseID == "" {
		rt.respondError(w, apperrors.NewBadRequest("course_id es requerido (query param)"))
		return
	}

	log.Printf("📚 Obteniendo tarea | Email: %s | ID: %s", email, taskID)

	task, err := rt.googleAdapter.GetTaskByID(r.Context(), email, courseID, taskID)
	if err != nil {
		log.Printf("❌ Error al obtener tarea: %v", err)
		rt.respondError(w, apperrors.NewGoogleAPIError(err))
		return
	}

	rt.respondJSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"task":   task,
	})
}

func (rt *Router) handleUpdateTask(w http.ResponseWriter, r *http.Request, email string) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/tasks/")
	taskID := strings.Split(path, "/")[0]

	var req struct {
		CourseID    string  `json:"course_id"`
		Title       string  `json:"title"`
		Description string  `json:"description"`
		DueDate     string  `json:"due_date"`
		MaxPoints   float64 `json:"max_points"`
		State       string  `json:"state"`
		WorkType    string  `json:"work_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rt.respondError(w, apperrors.NewBadRequestWithDetails("Error al parsear body", err.Error()))
		return
	}

	if req.CourseID == "" {
		rt.respondError(w, apperrors.NewValidationError("course_id", "es requerido"))
		return
	}
	if req.Title == "" {
		rt.respondError(w, apperrors.NewValidationError("title", "es requerido"))
		return
	}
	if len(req.Title) > 250 {
		rt.respondError(w, apperrors.NewValidationError("title", "no puede exceder 250 caracteres"))
		return
	}
	if req.MaxPoints < 0 {
		rt.respondError(w, apperrors.NewValidationError("max_points", "no puede ser negativo"))
		return
	}

	var dueDate time.Time
	var err error
	if req.DueDate != "" {
		dueDate, err = time.Parse(time.RFC3339, req.DueDate)
		if err != nil {
			rt.respondError(w, apperrors.NewBadRequestWithDetails("Formato de fecha inválido", "Usa ISO 8601"))
			return
		}
		if dueDate.Before(time.Now()) {
			rt.respondError(w, apperrors.NewValidationError("due_date", "no puede ser en el pasado"))
			return
		}
	}

	log.Printf("📝 Actualizando tarea | Email: %s | ID: %s", email, taskID)

	task := domain.Task{
		ID:          taskID,
		CourseID:    req.CourseID,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     dueDate,
		MaxPoints:   req.MaxPoints,
		State:       req.State,
		WorkType:    req.WorkType,
	}

	if err := rt.googleAdapter.UpdateTask(r.Context(), email, task); err != nil {
		log.Printf("❌ Error al actualizar tarea: %v", err)
		rt.respondError(w, apperrors.NewGoogleAPIError(err))
		return
	}

	rt.respondJSON(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Tarea actualizada exitosamente",
	})
}

func (rt *Router) handleDeleteTask(w http.ResponseWriter, r *http.Request, email string) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/tasks/")
	taskID := strings.Split(path, "/")[0]

	courseID := r.URL.Query().Get("course_id")
	if courseID == "" {
		rt.respondError(w, apperrors.NewBadRequest("course_id es requerido (query param)"))
		return
	}

	log.Printf("🗑️ Eliminando tarea | Email: %s | ID: %s", email, taskID)

	if err := rt.googleAdapter.DeleteTask(r.Context(), email, courseID, taskID); err != nil {
		log.Printf("❌ Error al eliminar tarea: %v", err)
		rt.respondError(w, apperrors.NewGoogleAPIError(err))
		return
	}

	rt.respondJSON(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Tarea eliminada exitosamente",
	})
}

// ============================================
// HANDLERS DE CLASSROOM - RUTAS DINÁMICAS (cursos)
// ============================================

func (rt *Router) handleCourseRoutes(w http.ResponseWriter, r *http.Request, email string) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/courses/")
	parts := strings.Split(path, "/")

	if len(parts) < 1 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}

	courseID := parts[0]

	// GET /api/v1/courses/{id}/students
	if len(parts) == 2 && parts[1] == "students" && r.Method == http.MethodGet {
		log.Printf("📚 Obteniendo estudiantes | Email: %s | Course: %s", email, courseID)

		students, err := rt.googleAdapter.GetStudentsByCourse(r.Context(), email, courseID)
		if err != nil {
			log.Printf("❌ Error al obtener estudiantes para %s: %v", email, err)
			rt.respondError(w, apperrors.NewGoogleAPIError(err))
			return
		}

		rt.respondJSON(w, http.StatusOK, map[string]interface{}{
			"status":    "success",
			"course_id": courseID,
			"students":  students,
			"count":     len(students),
		})
		return
	}

	// POST /api/v1/courses/{id}/students
	if len(parts) == 2 && parts[1] == "students" && r.Method == http.MethodPost {
		var req struct {
			StudentID string `json:"student_id"`
			Email     string `json:"email"`
			Name      string `json:"name"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			rt.respondError(w, apperrors.NewBadRequestWithDetails("Error al parsear body", err.Error()))
			return
		}

		if req.StudentID == "" {
			rt.respondError(w, apperrors.NewValidationError("student_id", "es requerido"))
			return
		}
		if req.Email == "" {
			rt.respondError(w, apperrors.NewValidationError("email", "es requerido"))
			return
		}
		if !strings.Contains(req.Email, "@") {
			rt.respondError(w, apperrors.NewValidationError("email", "inválido"))
			return
		}

		log.Printf("📚 Agregando estudiante | Email: %s | Course: %s", email, courseID)

		student := domain.Student{
			ID:       req.StudentID,
			Email:    req.Email,
			Name:     req.Name,
			CourseID: courseID,
		}

		if err := rt.googleAdapter.SaveStudent(r.Context(), email, courseID, student); err != nil {
			log.Printf("❌ Error al agregar estudiante: %v", err)
			rt.respondError(w, apperrors.NewGoogleAPIError(err))
			return
		}

		rt.respondJSON(w, http.StatusOK, map[string]string{
			"status":  "success",
			"message": "Estudiante agregado exitosamente",
		})
		return
	}

	// POST /api/v1/courses/{id}/sync
	if len(parts) == 2 && parts[1] == "sync" && r.Method == http.MethodPost {
		log.Printf("🔄 Sincronizando curso | Email: %s | Course: %s", email, courseID)

		students, err := rt.googleAdapter.GetStudentsByCourse(r.Context(), email, courseID)
		if err != nil {
			log.Printf("❌ Error al sincronizar para %s: %v", email, err)
			rt.respondError(w, apperrors.NewGoogleAPIError(err))
			return
		}

		rt.respondJSON(w, http.StatusOK, map[string]interface{}{
			"status":         "success",
			"message":        "Curso sincronizado exitosamente",
			"course_id":      courseID,
			"students_count": len(students),
		})
		return
	}

	http.NotFound(w, r)
}

// ============================================
// HANDLERS DE CLASSROOM - RUTAS DINÁMICAS (tareas)
// ============================================

func (rt *Router) handleTaskRoutes(w http.ResponseWriter, r *http.Request, email string) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/tasks/")
	parts := strings.Split(path, "/")

	if len(parts) < 1 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}

	taskID := parts[0]

	// GET /api/v1/tasks/{taskId}/submissions
	if len(parts) == 2 && parts[1] == "submissions" && r.Method == http.MethodGet {
		courseID := r.URL.Query().Get("course_id")
		if courseID == "" {
			rt.respondError(w, apperrors.NewBadRequest("course_id es requerido"))
			return
		}

		log.Printf("📚 Obteniendo entregas | Email: %s | Tarea: %s", email, taskID)

		submissions, err := rt.googleAdapter.GetTaskSubmissions(r.Context(), email, courseID, taskID)
		if err != nil {
			log.Printf("❌ Error al obtener entregas: %v", err)
			rt.respondError(w, apperrors.NewGoogleAPIError(err))
			return
		}

		rt.respondJSON(w, http.StatusOK, map[string]interface{}{
			"status":      "success",
			"task_id":     taskID,
			"submissions": submissions,
			"count":       len(submissions),
		})
		return
	}

	// POST /api/v1/tasks/{taskId}/grade
	if len(parts) == 2 && parts[1] == "grade" && r.Method == http.MethodPost {
		var req struct {
			CourseID     string  `json:"course_id"`
			SubmissionID string  `json:"submission_id"`
			Grade        float64 `json:"grade"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			rt.respondError(w, apperrors.NewBadRequestWithDetails("Error al parsear body", err.Error()))
			return
		}

		if req.CourseID == "" {
			rt.respondError(w, apperrors.NewValidationError("course_id", "es requerido"))
			return
		}
		if req.SubmissionID == "" {
			rt.respondError(w, apperrors.NewValidationError("submission_id", "es requerido"))
			return
		}
		if req.Grade < 0 || req.Grade > 100 {
			rt.respondError(w, apperrors.NewValidationError("grade", "debe estar entre 0 y 100"))
			return
		}

		log.Printf("📝 Calificando tarea | Email: %s | Tarea: %s | Nota: %.2f", email, taskID, req.Grade)

		err := rt.googleAdapter.GradeTask(r.Context(), email, req.CourseID, taskID, req.SubmissionID, req.Grade)
		if err != nil {
			log.Printf("❌ Error al calificar: %v", err)
			rt.respondError(w, apperrors.NewGoogleAPIError(err))
			return
		}

		rt.respondJSON(w, http.StatusOK, map[string]string{
			"status":  "success",
			"message": "Calificación asignada exitosamente",
		})
		return
	}

	http.NotFound(w, r)
}

// ============================================
// HANDLERS DE WEBHOOK
// ============================================

func (rt *Router) handleWebhook(w http.ResponseWriter, r *http.Request) {
	var payload map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("❌ Error al parsear webhook: %v", err)
		rt.respondError(w, apperrors.NewBadRequestWithDetails("Error al parsear webhook", err.Error()))
		return
	}

	log.Printf("📨 Webhook recibido")

	event := domain.ClassroomEvent{
		ResourceID: getString(payload, "resourceId"),
		EventType:  getString(payload, "eventType"),
		ReceivedAt: time.Now(),
	}

	go func() {
		ctx := context.Background()
		if err := rt.googleAdapter.SaveEvent(ctx, "admin@example.com", event); err != nil {
			log.Printf("❌ Error al guardar evento en Sheets: %v", err)
		}
	}()

	rt.respondJSON(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Webhook procesado correctamente",
	})
}

func (rt *Router) handleWebhookVerification(w http.ResponseWriter, r *http.Request) {
	challenge := r.URL.Query().Get("challenge")
	if challenge != "" {
		log.Printf("🔐 Verificando webhook")
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(challenge))
		return
	}

	rt.respondError(w, apperrors.NewBadRequest("Challenge no encontrado"))
}

// ============================================
// HANDLERS DE HEALTH CHECK
// ============================================

func (rt *Router) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	rt.respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":    "ok",
		"service":   "orienmod",
		"version":   "1.0.0",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// ============================================
// MIDDLEWARE DE AUTENTICACIÓN
// ============================================

type HandlerFunc func(http.ResponseWriter, *http.Request, string)

func (rt *Router) requireAuth(w http.ResponseWriter, r *http.Request, email string, fn HandlerFunc) {
	if email == "" {
		log.Printf("⛔ Intento de acceso sin autenticación | Path: %s", r.URL.Path)
		rt.respondError(w, apperrors.NewUnauthorized("Usuario no autenticado. Envía X-User-Email"))
		return
	}

	if _, err := rt.tokenStore.GetToken(r.Context(), email); err != nil {
		log.Printf("⛔ Token no encontrado | Email: %s", email)
		rt.respondError(w, apperrors.NewUnauthorized("Token no válido o expirado"))
		return
	}

	fn(w, r, email)
}

// ============================================
// FUNCIONES AUXILIARES
// ============================================

func (rt *Router) getUserEmail(ctx context.Context, token *oauth2.Token) (string, error) {
	client := rt.oauthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var userInfo struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return "", err
	}

	if userInfo.Email == "" {
		return "", err
	}

	return userInfo.Email, nil
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

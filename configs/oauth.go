// configs/oauth.go
package configs

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// NewOAuthConfigFromFile - Carga configuración desde credentials.json
func NewOAuthConfigFromFile() (*oauth2.Config, error) {
	// Rutas donde buscar el archivo
	paths := []string{
		"configs/credentials_wed.json",
		"configs/credentials.json",
		"credentials.json",
		"./configs/credentials.json",
	}

	var data []byte
	var err error
	var loadedPath string

	for _, path := range paths {
		data, err = os.ReadFile(path)
		if err == nil {
			loadedPath = path
			break
		}
	}

	if err != nil {
		return nil, fmt.Errorf("no se encontró credentials.json: %v", err)
	}

	log.Printf("✅ Credenciales cargadas desde: %s", loadedPath)

	// Estructura del archivo de Google
	var config struct {
		Web struct {
			ClientID     string   `json:"client_id"`
			ClientSecret string   `json:"client_secret"`
			AuthURI      string   `json:"auth_uri"`
			TokenURI     string   `json:"token_uri"`
			RedirectURIs []string `json:"redirect_uris"`
		} `json:"web"`
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error al parsear JSON: %v", err)
	}

	// ✅ VERIFICACIONES CRÍTICAS
	if config.Web.ClientID == "" {
		return nil, fmt.Errorf("❌ client_id no encontrado")
	}
	if config.Web.ClientSecret == "" {
		return nil, fmt.Errorf("❌ client_secret no encontrado")
	}
	if len(config.Web.RedirectURIs) == 0 {
		return nil, fmt.Errorf("❌ redirect_uris no encontrado")
	}

	// ✅ SCOPES COMPLETOS - LECTURA Y ESCRITURA PARA CLASSROOM
	scopes := []string{
		// ============================================
		// 📖 LECTURA DE CLASSROOM
		// ============================================
		"https://www.googleapis.com/auth/classroom.courses.readonly",
		"https://www.googleapis.com/auth/classroom.rosters.readonly",
		"https://www.googleapis.com/auth/classroom.profile.emails",
		"https://www.googleapis.com/auth/classroom.coursework.students.readonly",
		"https://www.googleapis.com/auth/classroom.guardianlinks.students.readonly",

		// ============================================
		// ✍️ ESCRITURA EN CLASSROOM
		// ============================================
		// Para crear y gestionar cursos
		"https://www.googleapis.com/auth/classroom.courses",

		// Para gestionar estudiantes
		"https://www.googleapis.com/auth/classroom.rosters",

		// ✅ PARA CREAR TAREAS COMO PROFESOR (NECESARIO)
		"https://www.googleapis.com/auth/classroom.coursework.me",

		// ✅ PARA VER Y CALIFICAR TAREAS DE ESTUDIANTES (NECESARIO)
		"https://www.googleapis.com/auth/classroom.coursework.students",

		// Para gestionar tutores
		"https://www.googleapis.com/auth/classroom.guardianlinks.students",

		// ============================================
		// 📊 GOOGLE SHEETS
		// ============================================
		"https://www.googleapis.com/auth/spreadsheets",

		// ============================================
		// 👤 INFORMACIÓN DEL USUARIO
		// ============================================
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	}

	log.Printf("✅ Configuración OAuth2 cargada correctamente")
	log.Printf("   Client ID: %s...", config.Web.ClientID[:15])
	log.Printf("   Redirect URI: %s", config.Web.RedirectURIs[0])
	log.Printf("   Scopes: %d scopes configurados (incluye escritura)", len(scopes))

	// Mostrar scopes importantes
	log.Printf("   📝 Scopes de escritura activos:")
	log.Printf("      - classroom.courses (crear cursos)")
	log.Printf("      - classroom.rosters (gestionar estudiantes)")
	log.Printf("      - classroom.coursework.me (crear tareas) ✅")
	log.Printf("      - classroom.coursework.students (calificar tareas) ✅")

	return &oauth2.Config{
		ClientID:     config.Web.ClientID,
		ClientSecret: config.Web.ClientSecret,
		RedirectURL:  config.Web.RedirectURIs[0],
		Scopes:       scopes,
		Endpoint:     google.Endpoint,
	}, nil
}

// NewOAuthConfig - Fallback a variables de entorno
func NewOAuthConfig() *oauth2.Config {
	config, err := NewOAuthConfigFromFile()
	if err == nil {
		return config
	}

	log.Printf("⚠️ Error al cargar credentials.json: %v", err)
	log.Println("⚠️ Usando variables de entorno como fallback")

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")

	if clientID == "" || clientSecret == "" || redirectURL == "" {
		log.Println("❌ Variables de entorno incompletas")
		log.Printf("   GOOGLE_CLIENT_ID: %s", clientID)
		log.Printf("   GOOGLE_CLIENT_SECRET: %s", clientSecret)
		log.Printf("   GOOGLE_REDIRECT_URL: %s", redirectURL)
		return nil
	}

	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			// 📖 LECTURA
			"https://www.googleapis.com/auth/classroom.courses.readonly",
			"https://www.googleapis.com/auth/classroom.rosters.readonly",
			"https://www.googleapis.com/auth/classroom.profile.emails",
			"https://www.googleapis.com/auth/classroom.coursework.students.readonly",

			// ✍️ ESCRITURA
			"https://www.googleapis.com/auth/classroom.courses",
			"https://www.googleapis.com/auth/classroom.rosters",
			"https://www.googleapis.com/auth/classroom.coursework.me",       // ✅ Crear tareas
			"https://www.googleapis.com/auth/classroom.coursework.students", // ✅ Calificar tareas

			// 📊 SHEETS Y PERFIL
			"https://www.googleapis.com/auth/spreadsheets",
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

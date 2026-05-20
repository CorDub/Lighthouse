package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"Lighthouse/handlers"
	"Lighthouse/internal/database"
	"Lighthouse/middleware"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	const port = "8080"

	// ENV variables
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	environment := os.Getenv("ENV")
	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleSecret := os.Getenv("GOOGLE_SECRET")
	mailtrapUsername := os.Getenv("MAILTRAP_USERNAME")
	mailtrapPassword := os.Getenv("MAILTRAP_PASSWORD")
	sameSite := http.SameSiteStrictMode
	secure := true
	from := os.Getenv("FROM_TEST")
	baseURL := "http://localhost:8080"
	// To DO LATER ON: add From versions for staging / prod

	if environment == "dev" {
		sameSite = http.SameSiteLaxMode
		secure = false
	}

	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	// Connecting to DB
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(dbConn)

	// config
	apiCfg := handlers.ApiConfig{
		DB: dbQueries,
		JWT: jwtSecret,
		Env: environment,
		Secure: secure,
		SameSite: sameSite,
		GoogleClientID: googleClientId,
		GoogleSecret: googleSecret,
		MailtrapUsername: mailtrapUsername,
		MailtrapPassword: mailtrapPassword,
		From: from,
		BaseURL: baseURL,
	}

	//server setup
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./dist"))
	//fallback
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := "./dist" + r.URL.Path
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			http.ServeFile(w, r, "./dist/index.html")
			return
		}
		fileServer.ServeHTTP(w, r)
	})

	//routes
	mux.HandleFunc("POST /api/login", apiCfg.Login)
	mux.HandleFunc("POST /api/logout", apiCfg.Logout)
	mux.HandleFunc("POST /api/users", apiCfg.CreateUser)
	mux.HandleFunc("POST /api/refresh", apiCfg.Refresh)
	mux.HandleFunc("GET /api/google", apiCfg.GoogleLogin)
	mux.HandleFunc("GET /api/callback", apiCfg.Callback)
	mux.HandleFunc("GET /api/checkAuth", apiCfg.CheckAuth)
	mux.HandleFunc("POST /api/checkPassword", apiCfg.CheckPassword)
	mux.HandleFunc("POST /api/changePassword", apiCfg.ChangePassword)

	//protected
	mux.Handle("GET /api/users", middleware.ValJWT(apiCfg.JWT, http.HandlerFunc(apiCfg.GetUsers)))

	handler := middleware.Cors(apiCfg.Env, mux)

	//server start
	server := &http.Server{
		Addr: ":" + port,
		Handler: handler,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
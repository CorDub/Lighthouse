package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"Lighthouse/handlers"
	"Lighthouse/internal/database"
	"Lighthouse/middleware"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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
	youtubeKey := os.Getenv("YOUTUBE_API_KEY")
	youtubeSecret := os.Getenv("YOUTUBE_SECRET")
	youtubeClientId := os.Getenv("YOUTUBE_CLIENT_ID")
	youtubeOAuthRedirectURL := os.Getenv("YOUTUBE_OAUTH_REDIRECT_URL")
	frontendOrigin := os.Getenv("FRONTEND_ORIGIN")
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

	// Creating YouTube SDK Service
	youtubeService, err := youtube.NewService(context.Background(), option.WithAPIKey(youtubeKey))
	if err != nil {
		log.Fatalf("failed to create youtube service: %v", err)
	}

	// Creating YouTube OAuth config
	youtubeOAuthConfig := &oauth2.Config{
		ClientID: youtubeClientId,
		ClientSecret: youtubeSecret,
		RedirectURL: youtubeOAuthRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/youtube.readonly",
			"https://www.googleapis.com/auth/yt-analytics.readonly",
		},
		Endpoint: google.Endpoint,
	}

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
		YouTubeService: youtubeService,
		YouTubeOAuthConfig: youtubeOAuthConfig,
		FrontEndOrigin: frontendOrigin,
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
	mux.HandleFunc("POST /api/users/creatorInvite", apiCfg.CreateCreatorFromInvite)
	mux.HandleFunc("GET /api/youtube/oauth/callback", apiCfg.YouTubeOAuthCallback)

	//protected
	mux.Handle("GET /api/users", middleware.ValJWT(apiCfg.JWT, http.HandlerFunc(apiCfg.GetUsers)))
	mux.Handle("PATCH /api/users/{id}/language", middleware.ValJWT(apiCfg.JWT, http.HandlerFunc(apiCfg.UpdateUserLanguage)))
	mux.Handle("POST /api/reports", middleware.ValJWT(apiCfg.JWT, http.HandlerFunc(apiCfg.CreateReport)))
	mux.Handle("GET /api/users/creatorsAvailable", middleware.ValJWT(apiCfg.JWT, http.HandlerFunc(apiCfg.GetCreatorsAvailable)))
	mux.Handle("POST /api/invite", middleware.ValJWT(apiCfg.JWT, http.HandlerFunc(apiCfg.CreateMagicLinkInvite)))
	mux.Handle("POST /api/connectChannel/youtube", middleware.ValJWT(apiCfg.JWT, http.HandlerFunc(apiCfg.ConnectYouTubeChannel)))

	handler := middleware.Cors(apiCfg.Env, mux)

	//server start
	server := &http.Server{
		Addr: ":" + port,
		Handler: handler,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}

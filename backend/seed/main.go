package main

import (
	"Lighthouse/internal/auth"
	"Lighthouse/handlers"
	"Lighthouse/internal/database"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)


type tentativeUser struct{
	Email string
	Password string
	Language database.Language
	Role database.Role
}


func seedUser(
	apiCfg handlers.ApiConfig, 
	email string, 
	password string,
	language database.Language,
	role database.Role) (string, error) {
	hash, err := auth.HashPassword(password)
	if err != nil {
		log.Printf("Couldn't hash the password: %s", err)
		return "", err
	}

	params := database.CreateUserParams{
		Email: email,
		HashedPassword: sql.NullString{
			String: hash,
			Valid: true,
		},
		Language: language,
		Role: role,
	}

	_, err2 := apiCfg.DB.CreateUser(context.Background(), params)
	if err2 != nil {
		log.Printf("Couldn't create user: %s", err2)
		return "", err2
	}

	exitString := fmt.Sprintf("User with email %s has been created", email)
	return exitString, nil
}


func main() {
	err := godotenv.Load(".env")
	if err != nil {
			log.Println("Error loading .env:", err)
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(dbConn)

	apiCfg := handlers.ApiConfig{
		DB: dbQueries,
	}

	log.Println("Seeding commencing")

	user1 := tentativeUser{
		Email: "corentindubois22@gmail.com",
		Password: "CavemanSwing78!",
		Language: "en",
		Role: "visitor",
	}
	user2 := tentativeUser{
		Email: "pedro.mcadmin@gmail.com",
		Password: "ProperlyDone45!",
		Language: "en",
		Role: "visitor",
	}
	user3 := tentativeUser{
		Email: "jorge.subadmin@gmail.com",
		Password: "YeahOK26!",
		Language: "en",
		Role: "visitor",
	}
	user4 := tentativeUser{
		Email: "agency@agency.com",
		Password: "VivaLosTokens",
		Language: "es",
		Role: "agency",
	}

	userSlice := []tentativeUser{user1, user2, user3, user4}

	for _, user := range userSlice {
		exitStr, err := seedUser(apiCfg, user.Email, user.Password, user.Language, user.Role)
		if err != nil {
			log.Printf("Couldn't seed user with email %s", user.Email)
			return
		}
		log.Println(exitStr)
	}
	log.Println("Database successfully seeded")
}
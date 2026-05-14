package handlers

import (
	"Lighthouse/test"
	"context"
	"database/sql"
	"errors"
	"net/http"
	"testing"
)

func TestCreateUser(t *testing.T) {
	db := th.SetupTestDB(t)

	email1 := "corentindubois22@gmail.com"
	password1 := "CavemanSwing78!"

	tests := []struct {
		name string
		email string
		password string
		dbPresence bool
		wantErr bool
	}{
		{
			name: "correct creation",
			email: email1,
			password: password1,
			dbPresence: true,
			wantErr: false,
		},
		{
			name: "duplicate email",
			email: email1,
			password: password1,
			dbPresence: false,
			wantErr: true,
		},
		{
			name: "no email",
			email: "",
			password: password1,
			dbPresence: false,
			wantErr: true,
		},
		{
			name: "invalid email",
			email: "jeronipo",
			password: password1,
			dbPresence: false,
			wantErr: true,
		},
		{
			name: "email too long",
			email: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqr",
			password: password1,
			dbPresence: false,
			wantErr: true,
		},
		{
			name: "no password",
			email: email1,
			password: "",
			dbPresence: false,
			wantErr: true,
		},
		{
			name: "password too short",
			email: email1,
			password: "itiswh",
			dbPresence: false,
			wantErr: true,
		},
		{
			name: "password too long",
			email: email1,
			password: "itiswhitiswhitiswhitiswhitiswhitiswhitiswhitiswhitiswhitiswhitiswhitiswhitiswhitiswhitiswh",
			dbPresence: false,
			wantErr: true,
		},
		{
			name: "password with non-printable characters",
			email: email1,
			password: "password\u0000",
			dbPresence: false,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//setup new transaction for each test
			queries := th.SetupTestTransaction(t, db)
			apiCfg := ApiConfig{DB: queries}

			payload := map[string]string{
				"email": test.email,
				"password": test.password,
			}

			w, r := th.CreateTestWR(t, "POST", "/api/users", payload)

			if test.name == "duplicate email" {
				th.AddUser(t, queries, test.email, test.password)
			}

			apiCfg.CreateUser(w, r)
			res := w.Result()

			// check error presence
			if res.StatusCode != http.StatusCreated {
				if !test.wantErr {
					errRes := th.DecodeResponse[th.ErrorResponse](t,w)
					t.Fatalf("Unexpected error: %s", errRes.Error)
				}
			}

			//check presence in db
			_, errDbPresence := apiCfg.DB.GetUserByEmail(context.Background(), test.email)
			if test.dbPresence {
				if errors.Is(errDbPresence, sql.ErrNoRows) {
					t.Errorf("No user found in db")
					return
				}

				if errDbPresence != nil {
					t.Errorf("error getting user from the database: %s", errDbPresence)
					return
				}
			}
		})
	}
}


func TestGetUsers(t *testing.T) {
	db := th.SetupTestDB(t)

	type tentativeUser struct{
		Email string
		Password string
	}

	//prepare initial state
	user1 := tentativeUser{
		Email: "corentindubois22@gmail.com",
		Password: "CavemanSwing78!",
	}
	user2 := tentativeUser{
		Email: "pedro.mcadmin@gmail.com",
		Password: "ProperlyDone45!",
	}
	user3 := tentativeUser{
		Email: "jorge.subadmin@gmail.com",
		Password: "YeahOK26!",
	}

	userSlice := []tentativeUser{user1, user2, user3}
	
	// prepare test cases
	tests := []struct{
		name string
		expUsersNum int
		wantErr bool
		status int
	}{
		{
			name: "correct",
			expUsersNum: 3,
			wantErr: false,
			status: http.StatusOK,
		},
		{
			name: "wrong amount",
			expUsersNum: 4,
			wantErr: true,
			status: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// setup transaction per test
			queries := th.SetupTestTransaction(t, db)
			apiCfg := ApiConfig{DB: queries}

			// setup initial state
			for _, user := range userSlice {
				th.AddUser(t, queries, user.Email, user.Password)
			}

			// create dummy w r
			w, r := th.CreateTestWR(t, "GET", "/api/users", nil)

			// run
			apiCfg.GetUsers(w, r)
			res := w.Result()

			// checks
			if !test.wantErr {
				//status
				if test.status != res.StatusCode {
					errRes := th.DecodeResponse[th.ErrorResponse](t, w)
					t.Fatalf("Unexpected error: %s", errRes.Error)
				}

				// expNumUsers
				decodedRes := th.DecodeResponse[getUsersResponse](t, w)
				if len(decodedRes.Users) != test.expUsersNum {
					t.Fatalf("Missing or extra users, exp: %d, got: %d", test.expUsersNum, len(decodedRes.Users))
				}
			}
		})
	}
}
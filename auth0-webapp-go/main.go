package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/auth0/go-auth0/authentication"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

var oauth2Config oauth2.Config

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	oauth2Config = oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("CALLBACK_URL"),
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("https://%s/authorize", os.Getenv("DOMAIN")),
			TokenURL: fmt.Sprintf("https://%s/oauth/token", os.Getenv("DOMAIN")),
		},
	}
}

var stateStore = make(map[string]bool)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/callback", callbackHandler)
	mux.HandleFunc("/logout", logoutHandler)

	handler := corsMiddleware(mux)

	log.Fatal(http.ListenAndServe(":3000", handler))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		allowedOrigin := "http://localhost:5500"

		if origin == allowedOrigin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session")
	if err != nil {

		fmt.Fprint(w, `
		<style>
			body {
				margin: 0;
				height: 40vh;
				display: flex;
				justify-content: center;
				align-items: center;
				font-family: Arial, sans-serif;
				background-color: #f5f5f5;
			}
			.login-button {
				display: inline-block;
				padding: 12px 24px;
				background-color: #eb5424;
				color: white;
				font-weight: bold;
				text-decoration: none;
				border-radius: 6px;
				box-shadow: 0 4px 6px rgba(0,0,0,0.1);
				transition: background-color 0.3s ease;
			}
			.login-button:hover {
				background-color: #cf481f;
			}
		</style>
		<a href="/login" class="login-button">Login with Auth0</a>
	`)

		return
	}

	auth0Client, err := authentication.New(
		context.Background(),
		os.Getenv("DOMAIN"),
		authentication.WithClientID(os.Getenv("CLIENT_ID")),
		authentication.WithClientSecret(os.Getenv("CLIENT_SECRET")),
	)
	if err != nil {
		http.Error(w, "Auth0 client error", http.StatusInternalServerError)
		return
	}

	userInfo, err := auth0Client.UserInfo(context.Background(), session.Value)
	if err != nil {
		http.Error(w, "User info error", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8" />
				<meta name="viewport" content="width=device-width, initial-scale=1" />
				<title>Welcome</title>
				<style>
					body {
						font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
						background-color: #f7f9fc;
						color: #333;
						margin: 0;
						padding: 0;
					}
					.container {
						max-width: 400px;
						margin: 80px auto;
						background-color: white;
						padding: 30px;
						border-radius: 10px;
						box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
						text-align: center;
					}
					h1 {
						margin-bottom: 10px;
						color: #2c3e50;
					}
					p {
						margin: 10px 0;
						font-size: 18px;
					}
					.logout-button {
						display: inline-block;
						margin-top: 25px;
						padding: 12px 28px;
						background-color: #e74c3c;
						color: white;
						font-weight: bold;
						text-decoration: none;
						border-radius: 6px;
						font-size: 16px;
						box-shadow: 0 4px 6px rgba(0,0,0,0.15);
						transition: background-color 0.3s ease;
					}
					.logout-button:hover {
						background-color: #c0392b;
					}
				</style>
			</head>
			<body>
				<div class="container">
					<h1>Welcome, %s!</h1>
					<a class="logout-button" href="/logout">Logout</a>
				</div>
			</body>
			</html>
			`, userInfo.FamilyName)
}

func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	state := generateState()
	stateStore[state] = true

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
	})

	url := oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusFound)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth_state")
	if err != nil || !stateStore[cookie.Value] {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}
	delete(stateStore, cookie.Value)

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing code", http.StatusBadRequest)
		return
	}

	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Token exchange failed", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "session",
		Value: token.AccessToken,
		Path:  "/",
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	logoutURL := fmt.Sprintf("https://%s/v2/logout?client_id=%s&returnTo=%s",
		os.Getenv("DOMAIN"), os.Getenv("CLIENT_ID"), "http://localhost:3000")
	http.Redirect(w, r, logoutURL, http.StatusFound)
}

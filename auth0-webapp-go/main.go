package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/auth0/go-auth0/authentication"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

var oauth2Config oauth2.Config

var (
	clientID     string
	clientSecret string
	domain       string
	callbackURL  string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	clientID = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
	domain = os.Getenv("DOMAIN")
	callbackURL = os.Getenv("CALLBACK_URL")

	if clientID == "" || clientSecret == "" || domain == "" || callbackURL == "" {
		log.Fatal("Missing required environment variables")
	}

	// Configure OAuth2
	oauth2Config = oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  callbackURL,
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("https://%s/authorize", domain),
			TokenURL: fmt.Sprintf("https://%s/oauth/token", domain),
		},
	}
}

var stateStore = make(map[string]bool)

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/callback", callbackHandler)
	http.HandleFunc("/logout", logoutHandler)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session")
	if err != nil {
		fmt.Fprint(w, `<a href="/login">Login with Auth0</a>`)
		return
	}

	auth0Client, err := authentication.New(
		context.Background(),
		domain,
		authentication.WithClientID(clientID),
		authentication.WithClientSecret(clientSecret),
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

	userJSON, _ := json.MarshalIndent(userInfo, "", "  ")
	fmt.Fprintf(w, `<h1>Welcome</h1><pre>%s</pre><a href="/logout">Logout</a>`, userJSON)
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
		Secure:   false, // change to true when not localhost anymore
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
		domain, clientID, "http://localhost:3000")
	http.Redirect(w, r, logoutURL, http.StatusFound)
}

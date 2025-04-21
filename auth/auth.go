package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"urodstvo-launcher/minecraft"

	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type AuthService struct {
	clientId     string
	redirectURI  string
	state        string
	codeVerifier string

	app          *application.App
	authWindow   *application.WebviewWindow
}

func NewAuthService() *AuthService {
	if err := godotenv.Load(".env"); err != nil{
		fmt.Println("not founded .env")
	}

	clientId := os.Getenv("MICROSOFT_CLIENT_ID")
	redirectURI := os.Getenv("MICROSOFT_REDIRECT_URI")
	if clientId == "" || redirectURI == "" {
		fmt.Println("CLIENT_ID & REDIRECT_URI not found in .env")
	}

	return &AuthService{
		clientId:    clientId,
		redirectURI: redirectURI,
	}
}

func (a *AuthService) SetApp(app *application.App) {
	a.app = app

	a.authWindow = a.app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Hidden: true,
		Width: 600,
		Height: 800,
		Title: "Microsoft Login",
	})

	a.startRedirectServer()
}

func (a *AuthService) startRedirectServer() {
	http.HandleFunc("/auth-callback", func(w http.ResponseWriter, r *http.Request) {
		code, err := minecraft.ParseAuthCodeURL(r.URL.String(), &a.state)
		if err != nil {
			a.app.EmitEvent("auth:microsoft:failed", err.Error())
			a.authWindow.Hide()
			return
		}

		resp, err := minecraft.CompleteLogin(a.clientId, "", a.redirectURI, code, a.codeVerifier)
		if err != nil {
			a.app.EmitEvent("auth:microsoft:failed", err.Error())
			a.authWindow.Hide()
			return
		}
		a.app.EmitEvent("auth:microsoft:success", resp)
		a.authWindow.Hide()
	})

	go func() {
		err := http.ListenAndServe(":34115", nil)
		if err != nil {
			log.Println("redirect server error:", err)
		}
	}()
}

func (a *AuthService) AddMicrosoftAccount() error {
	loginURL, state, codeVerifier, err := minecraft.GetSecureLoginData(a.clientId, a.redirectURI, nil)
	if err != nil {
		return fmt.Errorf("failed to get login URL: %w", err)
	}

	a.state = state
	a.codeVerifier = codeVerifier

	a.authWindow.SetURL(loginURL)

	a.authWindow.Show()
	return nil
}

func (a *AuthService) AddFreeAccount(username string) {
	a.app.EmitEvent("auth:free:success", username)
}

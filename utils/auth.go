package utils

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

func NewAuth() (*auth.Client, error) {
	app, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsJSON([]byte(os.Getenv("FIREBASE_JSON"))))
	if err != nil {
		return nil, err
	}

	client, err := app.Auth(context.Background())

	return client, nil
}

func Authenticate(r *http.Request) (string, error) {
	client, err := NewAuth()
	if err != nil {
		return "", err
	}

	header := r.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("no authorization header")
	}
	tokList := strings.Split(header, "Bearer ")
	if len(tokList) != 2 {
		return "", errors.New("no bearer token")
	}
	idToken := tokList[1]

	token, err := client.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		return "", err
	}

	return token.UID, nil
}

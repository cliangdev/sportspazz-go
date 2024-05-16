package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SignInWithPasswordRequest struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
}

func NewSignInWithPasswordRequest(email, password string) SignInWithPasswordRequest {
	return SignInWithPasswordRequest{
		Email:             email,
		Password:          password,
		ReturnSecureToken: true,
	}
}

type SignInWithPasswordResponse struct {
	IdToken      string `json:"idToken"`
	Email        string `json:"email"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
	LocalId      string `json:"localId"`
}

type FirebaseClient struct {
	apiKey string
}

func NewFirebaseClient(apiKey string) *FirebaseClient {
	return &FirebaseClient{
		apiKey: apiKey,
	}
}

func (c *FirebaseClient) SignInWithEmailAndPassword(req SignInWithPasswordRequest) (*SignInWithPasswordResponse, error) {
	url := "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=" + c.apiKey
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to sign in, status code %d", resp.StatusCode)
	}

	var res SignInWithPasswordResponse
	json.NewDecoder(resp.Body).Decode(&res)

	return &res, nil
}

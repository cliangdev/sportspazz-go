package client

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	idTokenCertURL        = "https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com"
	signInWithPasswordURL = "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key="
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
	apiKey     string
	projectID  string
	publicKeys map[string]*rsa.PublicKey
	logger     *slog.Logger
}

func NewFirebaseClient(apiKey, projectID string, logger *slog.Logger) *FirebaseClient {
	return &FirebaseClient{
		apiKey:    apiKey,
		projectID: projectID,
		logger:    logger,
	}
}

func (c *FirebaseClient) SignInWithEmailAndPassword(req SignInWithPasswordRequest) (*SignInWithPasswordResponse, error) {
	url := signInWithPasswordURL + c.apiKey
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

// Cache or move it outside.
// Unable to use Firebase admin SDK due to private Cache-Control response header without a max-age value.
func (c *FirebaseClient) fetchPublicKeys() error {
	resp, err := http.Get(idTokenCertURL)
	if err != nil {
		return fmt.Errorf("error fetching public keys: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error fetching public keys: status %d", resp.StatusCode)
	}

	keys := make(map[string]string)
	err = json.NewDecoder(resp.Body).Decode(&keys)
	if err != nil {
		return fmt.Errorf("error decoding public keys: %v", err)
	}

	c.publicKeys = make(map[string]*rsa.PublicKey)
	for kid, keyStr := range keys {
		key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(keyStr))
		if err != nil {
			return fmt.Errorf("error parsing public key: %v", err)
		}
		c.publicKeys[kid] = key
	}
	return nil
}

func (c *FirebaseClient) VerifyIDToken(idToken string) (*jwt.Token, error) {
	c.fetchPublicKeys()

	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token")
	}

	headerSegment := parts[0]
	headerBytes, err := jwt.DecodeSegment(headerSegment)
	if err != nil {
		return nil, fmt.Errorf("error decoding header: %v", err)
	}

	var header struct {
		Kid string `json:"kid"`
	}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return nil, fmt.Errorf("error unmarshaling header: %v", err)
	}

	key, ok := c.publicKeys[header.Kid]
	if !ok {
		return nil, fmt.Errorf("public key not found")
	}

	token, err := jwt.Parse(idToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	if err := c.verifyClaims(claims); err != nil {
		return nil, err
	}

	return token, nil
}

func (c *FirebaseClient) verifyClaims(claims jwt.MapClaims) error {
	now := time.Now().Unix()
	if !claims.VerifyExpiresAt(now, true) {
		return fmt.Errorf("token has expired")
	}
	if !claims.VerifyIssuedAt(now, true) {
		return fmt.Errorf("token used before issued")
	}
	if !claims.VerifyAudience(c.projectID, true) {
		return fmt.Errorf("invalid audience")
	}
	if !claims.VerifyIssuer(fmt.Sprintf("https://securetoken.google.com/%s", c.projectID), true) {
		return fmt.Errorf("invalid issuer")
	}
	return nil
}

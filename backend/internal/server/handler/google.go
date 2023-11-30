package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"wordwiz/pkg/logger"
	"wordwiz/pkg/security"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	SecretKey                 = "MySuperSecretKey123qwerty"
	AccessTokenExpirationTime = 100 * 24 * time.Hour
	oauthGoogleUrlAPI         = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
)

var (
	googleOauthConfig *oauth2.Config
)

func init() {
	googleOauthConfig = &oauth2.Config{
		ClientID:     "515365069373-1v8kedreb6qjp7nbhnekqnhiv3bra0r5.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "http://localhost:8080/auth/google/callback", // Change this to your redirect URL
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}
}

func (h *Handler) GoogleLogin(c *gin.Context) {
	oauthState := generateStateOauthCookie(c.Writer)
	url := googleOauthConfig.AuthCodeURL(oauthState)
	http.Redirect(c.Writer, c.Request, url, http.StatusTemporaryRedirect)
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func (h *Handler) HandleGoogleCallback(c *gin.Context) {
	// Read oauthState from Cookie
	oauthState, err := c.Cookie("oauthstate")
	if err != nil {
		logger.Get().Err(err).Msg("Error retrieving oauth state from cookie")
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("error retrieving oauth state from cookie"))
		return
	}

	if c.Query("state") != oauthState {
		logger.Get().Err(err).Msg("Invalid OAuth Google state")
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("invalid OAuth Google state"))
		return
	}

	userInfo, err := GetGoogleUserInfo(c.Request.Context(), c.Query("code"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	firstName := userInfo["given_name"].(string)
	lastName := userInfo["family_name"].(string)
	email := userInfo["email"].(string)
	fullName := firstName + " " + lastName

	userID, err := h.stg.SignInUp(c.Request.Context(), email, fullName)
	if err != nil {
		logger.Get().Err(err).Msg("Failed to SignInUp")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tokenInfo := map[string]interface{}{
		"uid": userID,
		"ip":  c.ClientIP(),
	}

	// Generate token
	accessToken, err := security.GenerateJWT(tokenInfo, AccessTokenExpirationTime, SecretKey)
	if err != nil {
		logger.Get().Err(err).Msg("Failed to GenerateJWT")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "Authorization",
		Value:   fmt.Sprintf("Bearer %s", accessToken),
		Expires: time.Now().Add(time.Hour * 24), // Token expires in 24 hours
		Path:    "/",
	})

	c.JSON(http.StatusOK, gin.H{"msg": "Login successful! Welcome back to your account."})
}

// GetGoogleUserInfo retrieves user information using the Google token.
func GetGoogleUserInfo(ctx context.Context, code string) (userInfo map[string]interface{}, err error) {
	token, err := googleOauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}

	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}

	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authCookie, err := c.Request.Cookie("Authorization")
		if err != nil {
			logger.Get().Err(err).Msg("get cookie")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "get token from cookie"})
			return
		}

		if authCookie == nil || authCookie.Value == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return
		}

		tokenParts := strings.Split(authCookie.Value, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			return
		}

		tokenString := tokenParts[1]

		claims, err := security.ExtractClaims(tokenString, SecretKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		uid, ok := claims["uid"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "cannot parse 'uid' field"})
			return
		}

		c.Set("uid", uid)

		c.Next()
	}
}

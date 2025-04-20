package minecraft

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)


const (
	_authURL = "https://login.microsoftonline.com/consumers/oauth2/v2.0/authorize"
	_tokenURL = "https://login.microsoftonline.com/consumers/oauth2/v2.0/token"
	_scope   = "XboxLive.signin offline_access"
)

func generatePKCEData() (codeVerifier string, codeChallenge string, codeChallengeMethod string, err error) {
	randomBytes := make([]byte, 96) 
	if _, err = io.ReadFull(rand.Reader, randomBytes); err != nil {
		return
	}

	codeVerifier = base64.RawURLEncoding.EncodeToString(randomBytes)
	if len(codeVerifier) > 128 {
		codeVerifier = codeVerifier[:128]
	}

	hash := sha256.Sum256([]byte(codeVerifier))
	codeChallenge = base64.RawURLEncoding.EncodeToString(hash[:])
	codeChallengeMethod = "S256"
	return
}

func generateState() (string, error) {
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(randomBytes), nil
}

func GetLoginURL(clientID, redirectURI string) (string, error) {
	values := url.Values{}
	values.Set("client_id", clientID)
	values.Set("response_type", "code")
	values.Set("redirect_uri", redirectURI)
	values.Set("response_mode", "query")
	values.Set("scope", _scope)

	u, err := url.Parse(_authURL)
	if err != nil {
		return "", err
	}
	u.RawQuery = values.Encode()
	return u.String(), nil
}

func GetSecureLoginData(clientID, redirectURI string, stateOpt *string) (loginURL, state, codeVerifier string, err error) {
	codeVerifier, codeChallenge, codeChallengeMethod, err := generatePKCEData()
	if err != nil {
		return
	}

	var stateVal string
	if stateOpt != nil {
		stateVal = *stateOpt
	} else {
		stateVal, err = generateState()
		if err != nil {
			return
		}
	}

	values := url.Values{}
	values.Set("client_id", clientID)
	values.Set("response_type", "code")
	values.Set("redirect_uri", redirectURI)
	values.Set("response_mode", "query")
	values.Set("scope", _scope)
	values.Set("state", stateVal)
	values.Set("code_challenge", codeChallenge)
	values.Set("code_challenge_method", codeChallengeMethod)

	u, err := url.Parse(_authURL)
	if err != nil {
		return
	}
	u.RawQuery = values.Encode()
	return u.String(), stateVal, codeVerifier, nil
}

func getAuthorizationToken(clientID, redirectURI, authCode string, clientSecret, codeVerifier string) (*AuthorizationTokenResponse, error) {
	values := url.Values{}
	values.Set("client_id", clientID)
	values.Set("scope", _scope)
	values.Set("code", authCode)
	values.Set("redirect_uri", redirectURI)
	values.Set("grant_type", "authorization_code")

	
	if clientSecret != "" {
		values.Set("client_secret", clientSecret)
	}
	if codeVerifier != "" {
		values.Set("code_verifier", codeVerifier)
	}
	
	req, err := http.NewRequest("POST", _tokenURL, bytes.NewBufferString(values.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", getUserAgent())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result AuthorizationTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func refreshAuthorizationToken(clientID, clientSecret string, redirectURI *string, refreshToken string) (*AuthorizationTokenResponse, error) {
	form := url.Values{}
	form.Set("client_id", clientID)
	form.Set("scope", "XboxLive.signin offline_access")
	form.Set("refresh_token", refreshToken)
	form.Set("grant_type", "refresh_token")

	if clientSecret != "" {
		form.Set("client_secret", clientSecret)
	}

	_ = redirectURI

	req, err := http.NewRequest("POST", "https://login.live.com/oauth20_token.srf", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", getUserAgent())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result *AuthorizationTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func authenticateWithXBL(accessToken string) (*XBLResponse, error) {
	body := map[string]any{
		"Properties": map[string]any{
			"AuthMethod": "RPS",
			"SiteName":   "user.auth.xboxlive.com",
			"RpsTicket":  "d=" + accessToken,
		},
		"RelyingParty": "http://auth.xboxlive.com",
		"TokenType":    "JWT",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://user.auth.xboxlive.com/user/authenticate", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", getUserAgent())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result *XBLResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func authenticateWithXSTS(xblToken string) (*XSTSResponse, error) {
	body := map[string]any{
		"Properties": map[string]any{
			"SandboxId": "RETAIL",
			"UserTokens": []string{
				xblToken,
			},
		},
		"RelyingParty": "rp://api.minecraftservices.com/",
		"TokenType":    "JWT",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://xsts.auth.xboxlive.com/xsts/authorize", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", getUserAgent())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result XSTSResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func authenticateWithMinecraft(userhash, xstsToken string) (*MinecraftAuthenticateResponse, error) {
	body := map[string]string{
		"identityToken": fmt.Sprintf("XBL3.0 x=%s;%s", userhash, xstsToken),
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.minecraftservices.com/authentication/login_with_xbox", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", getUserAgent())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result *MinecraftAuthenticateResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetStoreInformation(accessToken string) (*MinecraftStoreResponse, error) {
	req, err := http.NewRequest("GET", "https://api.minecraftservices.com/entitlements/mcstore", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("User-Agent", getUserAgent())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result *MinecraftStoreResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getProfile(accessToken string) (*MinecraftProfileResponse, error) {
	req, err := http.NewRequest("GET", "https://api.minecraftservices.com/minecraft/profile", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("User-Agent", getUserAgent())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result *MinecraftProfileResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func CompleteLogin(clientID, clientSecret, redirectURI, authCode, codeVerifier string) (*CompleteLoginResponse, error) {
	tokenResp, err := getAuthorizationToken(clientID, redirectURI, authCode, clientSecret, codeVerifier)
	if err != nil {
		return nil, err
	}
	accessToken := tokenResp.AccessToken
	
	xblResp, err := authenticateWithXBL(accessToken)
	if err != nil {
		return nil, err
	}
	xblToken := xblResp.Token
	userhash := xblResp.DisplayClaims.Xui[0].Uhs
	
	xstsResp, err := authenticateWithXSTS(xblToken)
	if err != nil {
		return nil, err
	}
	xstsToken := xstsResp.Token
	
	mcResp, err := authenticateWithMinecraft(userhash, xstsToken)
	if err != nil {
		return nil, err
	}
	mcAccessToken := mcResp.AccessToken
	if mcAccessToken == "" {
		return nil, errors.New("AzureAppNotPermitted")
	}
	
	profile, err := getProfile(mcAccessToken)
	if err != nil {
		return nil, err
	}
	if profile.Error != "" && profile.Error == "NOT_FOUND" {
		return nil, errors.New("AccountNotOwnMinecraft")
	}
	var response *CompleteLoginResponse = &CompleteLoginResponse{
		MinecraftProfileResponse: *profile,
		AccessToken: mcResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
	}
	return response, nil
}

func CompleteRefresh(clientID, clientSecret, redirectURI, refreshToken string) (*CompleteLoginResponse, error) {
	tokenResp, err := refreshAuthorizationToken(clientID, clientSecret, &redirectURI, refreshToken)
	if err != nil {
		return nil, errors.New("InvalidRefreshToken")
	}
	
	accessToken := tokenResp.AccessToken

	xblResp, err := authenticateWithXBL(accessToken)
	if err != nil {
		return nil, err
	}
	xblToken := xblResp.Token
	userhash := xblResp.DisplayClaims.Xui[0].Uhs

	xstsResp, err := authenticateWithXSTS(xblToken)
	if err != nil {
		return nil, err
	}
	xstsToken := xstsResp.Token

	mcResp, err := authenticateWithMinecraft(userhash, xstsToken)
	if err != nil {
		return nil, err
	}
	mcAccessToken := mcResp.AccessToken

	profile, err := getProfile(mcAccessToken)
	if err != nil {
		return nil, err
	}
	if profile.Error != "" && profile.Error == "NOT_FOUND" {
		return nil, errors.New("AccountNotOwnMinecraft")
	}

	var response *CompleteLoginResponse = &CompleteLoginResponse{
		MinecraftProfileResponse: *profile,
		AccessToken: mcResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
	}
	return response, nil
}


func UrlContainsAuthCode(rawURL string) bool {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	q := parsed.Query()
	_, ok := q["code"]
	return ok
}


func GetAuthCodeFromURL(rawURL string) *string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil
	}
	q := parsed.Query()
	if val, ok := q["code"]; ok && len(val) > 0 {
		return &val[0]
	}
	return nil
}

func ParseAuthCodeURL(rawURL string, expectedState *string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	q := parsed.Query()

	if expectedState != nil {
		state := q.Get("state")
		if state != *expectedState {
			return "", errors.New("state mismatch")
		}
	}

	code := q.Get("code")
	if code == "" {
		return "", errors.New("authorization code not found in URL")
	}

	return code, nil
}

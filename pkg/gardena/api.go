package gardena

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const baseURL = "https://api.smart.gardena.dev/v1"
const husqvarnaTokenURL = "https://api.authentication.husqvarnagroup.dev/v1/oauth2/token"
const ApiHealthURL = "/health"
const LocationsURL = "/locations"

const clientIDFile = "client-id"
const clientSecretFile = "client-secret"

type API struct {
	baseURL        string
	authUrl        string
	httpClient     *http.Client
	secretFilePath string

	clientID     string
	clientSecret string
	accessToken  string
	userID       string
	tokenExpAt   time.Time
}

type authResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
	Provider    string `json:"provider"`
	UserID      string `json:"user_id"`
	TokenType   string `json:"token_type"`
}

type APIBuilder struct {
	api *API
}

// NewAPI returns an APIBuilder with the default base/authentication url
// as well as an httpClient preconfigured.
// The default urls can be overwritten by the corresponding 'with%' methods.
func NewAPI() *APIBuilder {
	return &APIBuilder{api: &API{
		baseURL:    baseURL,
		authUrl:    husqvarnaTokenURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}}
}

// WithBaseURL sets a given base url for the APIBuilder
func (b *APIBuilder) WithBaseURL(u string) *APIBuilder {
	b.api.baseURL = u
	return b
}

// WithAuthURL sets a given authentication url for the APIBuilder
func (b *APIBuilder) WithAuthURL(u string) *APIBuilder {
	b.api.authUrl = u
	return b
}

// WithSecretFilePath sets the path the required secret files
func (b *APIBuilder) WithSecretFilePath(p string) *APIBuilder {
	b.api.secretFilePath = p
	return b
}

// Initialize initializes the API from the Builder.
// The client-id and client-secret are read from the configured secret files. Also,
// the api authenticates, acquiring an access token.
func (b *APIBuilder) Initialize() (*API, error) {
	api := b.api
	if api.secretFilePath == "" {
		return nil, fmt.Errorf("secretpath can not be empty")
	}

	if !strings.HasSuffix(api.secretFilePath, "/") {
		api.secretFilePath = api.secretFilePath + "/"
	}
	clientID, err := readFromSecretFile(api.secretFilePath + clientIDFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read client-id from secret file, got err:\n %w", err)
	}
	api.clientID = clientID
	clientSecret, err := readFromSecretFile(api.secretFilePath + clientSecretFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read client-secret from secret file, got err:\n %w", err)
	}
	api.clientSecret = clientSecret

	if err := api.authenticate(); err != nil {
		return nil, fmt.Errorf("unable to authenticate, got err:\n %w", err)
	}
	log.Println("Successfully initialized gardena smart system api!")
	return api, nil
}

// authenticate requests an access token from the configured authentication endpoint and stores it in the API
func (api *API) authenticate() error {
	if api.clientID == "" || api.clientSecret == "" {
		return fmt.Errorf("api not initialized, client-id or client-secret was empty")
	}
	if api.accessToken == "" || !api.tokenExpAt.IsZero() && time.Now().After(api.tokenExpAt.Add(time.Duration(-3600))) {
		res, err := api.httpClient.PostForm(api.authUrl, url.Values{
			"grant_type":    {"client_credentials"},
			"client_id":     {api.clientID},
			"client_secret": {api.clientSecret},
		})
		if err != nil {
			return fmt.Errorf("unable to request access token, got err %w", err)
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("unable to read authentication response, got error: %w", err)
		}
		var auth authResponse
		if err := json.Unmarshal(body, &auth); err != nil {
			return fmt.Errorf("unable to parse authentication response to json, err: %w", err)
		}
		api.userID = auth.UserID
		api.accessToken = auth.TokenType + " " + auth.AccessToken
		api.tokenExpAt = time.Now().Add(time.Second * time.Duration(auth.ExpiresIn))
	}
	return nil
}

// GetAPIHealthURL returns the health url for the API with the configured base url
func (api *API) GetAPIHealthURL() string {
	return api.baseURL + ApiHealthURL
}

// GetBaseURL returns the configured base url
func (api *API) GetBaseURL() string {
	return api.baseURL
}

// query sets up an HTTP GET request against the configured base url + the given path, using the
// configured client id and access token. The response is returned as http.Response
func (api *API) query(path string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, api.baseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to setup request for endpoint %s%s, got err:\n %w", baseURL, path, err)
	}
	req.Header.Set("X-Api-Key", api.clientID)
	req.Header.Set("Authorization", api.accessToken)
	res, err := api.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to query endpoint %s%s, got err:\n %w", baseURL, path, err)
	}
	return res, nil
}

// GetLocations queries the locations of the LocationsURL and returns the result as json
func (api *API) GetLocations() (*Locations, error) {
	res, err := api.query(LocationsURL)
	if err != nil {
		return nil, fmt.Errorf("querying for locations failed, got err:\n %w", err)
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("unable to read response, got status code %d, response: %v", res.StatusCode, res)
	}

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response, got error:\n%w", err)
	}

	locations := Locations{}
	if err := json.Unmarshal(responseBody, &locations); err != nil {
		return nil, fmt.Errorf("unmarshal of locations response failed, got err:\n%w", err)
	}
	return &locations, nil
}

func (api *API) GetInitialStateFor(location Location) (*State, error) {
	res, err := api.query(LocationsURL + "/" + location.Id)
	if err != nil {
		return nil, fmt.Errorf("unable to query location %s, got err:\n%w", location.Id, err)
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("unable to read response, got status code %d, response: %v", res.StatusCode, res)
	}

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unmarshal of locations response failed, got err:\n%w", err)
	}

	state := State{}
	if err := json.Unmarshal(responseBody, &state); err != nil {
		return nil, fmt.Errorf("unable to unmarshal json state, got error: %v", res)
	}
	return &state, nil
}

// readFromSecretFile reads the file on the given path and returns the value as string
func readFromSecretFile(path string) (string, error) {
	val, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("unable to read secret file at '%s', got err:\n '%w'", path, err)
	}
	s := string(val)
	return strings.TrimSuffix(s, "\n"), nil
}

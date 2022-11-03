package gardena

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestAPIInitialize(t *testing.T) {
	expectedAccessToken := "abc123def456"
	expectedUserID := "some-user"
	expectedTokenType := "Bearer"
	expectedExpireIn := 86399
	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
				"access_token": "` + expectedAccessToken + `",
				"token_type": "` + expectedTokenType + `",
				"expires_in": ` + fmt.Sprint(expectedExpireIn) + `,
				"user_id": "` + expectedUserID + `"
			}`))
	}))
	defer authServer.Close()

	expectedClientID := "<some-client-id>"
	expectedClientSecret := "<some-client-secret>"
	secretFilePath := setupSecretFilesWithTmpDir(expectedClientID, expectedClientSecret)
	defer os.RemoveAll(secretFilePath)

	apiStub, err := NewAPI().
		WithAuthURL(authServer.URL).
		WithSecretFilePath(secretFilePath).
		Initialize()
	if err != nil {
		log.Fatalf("Unable to initialize api, got err %v", err)
	}
	if apiStub.clientID != expectedClientID {
		log.Fatalf("Wrong clientID, expected %s, got %s", expectedClientID, apiStub.clientID)
	}
	if apiStub.clientSecret != expectedClientSecret {
		log.Fatalf("Wrong clientSecret, expected %s, got %s", expectedClientSecret, apiStub.clientSecret)
	}
	if apiStub.accessToken != expectedTokenType+" "+expectedAccessToken {
		log.Fatalf("Wrong accessToken, expected %s, got %s", expectedAccessToken, apiStub.accessToken)
	}
	if apiStub.userID != expectedUserID {
		log.Fatalf("Wrong user id, expected %v, got %v", expectedUserID, apiStub.userID)
	}

	// Check that expire time is not 0 or after now+expire_in because of time passed between function call and test
	// should be sufficient and bypasses need for mocking/injecting time.Now()
	if apiStub.tokenExpAt.IsZero() || apiStub.tokenExpAt.After(time.Now().Add(time.Second*time.Duration(expectedExpireIn))) {
		log.Fatalf("Wrong tokenExpAt, got '%v', now is '%v'", apiStub.tokenExpAt, time.Now())
	}
}

func TestLoadLocationsFromApi(t *testing.T) {
	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"access_token": "abc123def456"}`))
	}))
	defer authServer.Close()

	locationsProvider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		locationsResponse, err := os.ReadFile("../../test/locations.json")
		if err != nil {
			log.Fatal(err)
		}
		w.Write(locationsResponse)
	}))
	defer locationsProvider.Close()

	secretFilePath := setupSecretFilesWithTmpDir("<some-client-id>", "<some-client-secret>")
	defer os.RemoveAll(secretFilePath)

	apiStub, err := NewAPI().
		WithBaseURL(locationsProvider.URL).
		WithAuthURL(authServer.URL).
		WithSecretFilePath(secretFilePath).
		Initialize()
	if err != nil {
		t.Fatalf("Unable to initialize api, got err: ")
	}

	allLocations, err := apiStub.GetLocations()
	if err != nil {
		t.Fatalf("Unable to query locations, got err: ")
	}
	if len(allLocations.Data) != 1 {
		t.Fatalf("List of locations should be 1, is '%v'", len(allLocations.Data))
	}

	expectedLocation := LocationFromApi{
		Id:   "123abc",
		Type: typeLocation,
	}
	expectedLocation.Attributes.Name = "My Garden"
	locationOne := allLocations.Data[0].LocationFromApi
	if !reflect.DeepEqual(locationOne, expectedLocation) {
		t.Fatalf("Expected '%v', got '%v'", expectedLocation, locationOne)
	}
}

func TestAPI_authenticateRefreshExpiredToken(t *testing.T) {
	expectedAccessToken := "abc123def456"
	expectedTokenType := "Bearer"
	expectedExpireIn := 86399
	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
				"access_token": "` + expectedAccessToken + `",
				"token_type": "` + expectedTokenType + `",
				"expires_in": ` + fmt.Sprint(expectedExpireIn) + `
			}`))
	}))
	defer authServer.Close()

	api := API{
		baseURL:    "www.baseurl.com",
		authUrl:    authServer.URL,
		httpClient: &http.Client{},

		clientID:     "abc",
		clientSecret: "def",
		accessToken:  expectedTokenType + " 987zyx",
		tokenExpAt:   time.Now().Add(time.Hour * -24),
	}
	if err := api.authenticate(); err != nil {
		log.Fatalf("Authentication failed with err:\n %v", err)
	}

	if api.accessToken != expectedTokenType+" "+expectedAccessToken {
		log.Fatalf("Access token wasn't refreshed, expected '%s', got '%v'", expectedTokenType+" "+expectedAccessToken, api.accessToken)
	}
}

// setupSecretFilesWithTmpDir creates a tmp dir and writes the provided information into the expected
// files. Cleanup with defer os.RemoveAll(secretFilePath)
func setupSecretFilesWithTmpDir(clientID, clientSecret string) string {
	expectedClientID := clientID
	secretFilePath, err := os.MkdirTemp("", "gxxs")
	if err != nil {
		log.Fatal(err)
	}
	clientIdFile := filepath.Join(secretFilePath, "client-id")
	if err := os.WriteFile(clientIdFile, []byte(expectedClientID), 0666); err != nil {
		log.Fatal(err)
	}

	expectedClientSecret := clientSecret
	clientSecretFile := filepath.Join(secretFilePath, "client-secret")
	if err := os.WriteFile(clientSecretFile, []byte(expectedClientSecret), 0666); err != nil {
		log.Fatal(err)
	}
	return secretFilePath
}

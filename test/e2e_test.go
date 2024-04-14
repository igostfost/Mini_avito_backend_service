package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/igostfost/avito_backend_trainee/pkg/types"
	"net/http"
	"testing"
)

func TestAuth(t *testing.T) {
	serverURL := "http://localhost:8000"

	// Тестовые эндпоинты для проверки
	endpoints := []struct {
		method   string
		path     string
		query    string
		body     []byte
		expected int
	}{
		{"POST", "/banner", "", []byte("fds"), http.StatusUnauthorized},
		{"GET", "/user_banner", "?tag_id=1&feature_id=1", nil, http.StatusUnauthorized},
		{"GET", "/banner", "", nil, http.StatusUnauthorized},
		{"DELETE", "/banner/1", "", nil, http.StatusUnauthorized},
		{"PATCH", "/banner/1", "", nil, http.StatusUnauthorized},
	}

	for _, endpoint := range endpoints {
		var resp *http.Response
		var err error

		switch endpoint.method {
		case "GET":
			resp, err = http.Get(serverURL + endpoint.path + endpoint.query)
		case "POST":
			resp, err = http.Post(serverURL+endpoint.path+endpoint.query, "application/json", bytes.NewBuffer(endpoint.body))
		case "DELETE":
			req, _ := http.NewRequest("DELETE", serverURL+endpoint.path, nil)
			resp, err = http.DefaultClient.Do(req)
		case "PATCH":
			req, _ := http.NewRequest("PATCH", serverURL+endpoint.path, bytes.NewBuffer(endpoint.body))
			resp, err = http.DefaultClient.Do(req)
		}

		if err != nil {
			t.Errorf("Failed to send %s request to %s: %v", endpoint.method, endpoint.path, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != endpoint.expected {
			t.Errorf("Unexpected status code from %s request to %s. Expected: %d, Got: %d", endpoint.method, endpoint.path, endpoint.expected, resp.StatusCode)
			continue
		}
		fmt.Printf("CHECK %s AUTH TO %s  - SUCCESS\n", endpoint.method, endpoint.path)
	}
}

func TestSignUpAdmin(t *testing.T) {
	serverURL := "http://localhost:8000/auth/sign-up/admin"

	inputData := map[string]string{
		"username": "testAdmin",
		"password": "testPassword",
	}
	jsonData, err := json.Marshal(inputData)
	if err != nil {
		t.Fatalf("Failed to marshal JSON data: %v", err)
	}

	resp, err := http.Post(serverURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}

	var responseData struct {
		ID      int `json:"id"`
		IsAdmin int `json:"is_admin"`
	}
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	fmt.Println("TestSignUpAdmin -  SUCCESS")
}

func TestPostBanner(t *testing.T) {

	// ------ AUTH START -------
	signInURL := "http://localhost:8000/auth/sign-in/admin"
	inputData := map[string]string{
		"username": "testAdmin",
		"password": "testPassword",
	}
	jsonData, err := json.Marshal(inputData)
	if err != nil {
		t.Fatalf("Failed to marshal JSON data: %v", err)
	}

	resp, err := http.Post(signInURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}

	var responseToken struct {
		Token string `json:"token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&responseToken)
	if err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}
	token := responseToken.Token
	// ------ AUTH END -------

	// ------ POST BANNER START -------
	inputBanner := &types.BannerRequest{
		TagIds:    []int{1},
		FeatureId: 1,
		Content: types.Content{
			Title: "Test Banner",
			Text:  "Test content",
			Url:   "http://example.com",
		},
		IsActive: true,
	}
	jsonData, err = json.Marshal(inputBanner)
	if err != nil {
		t.Fatalf("Failed to marshal JSON data: %v", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:8000/banner", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create POST request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}

	var responseData struct {
		BannerId int `json:"bannerId"`
	}
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}
	// ------ POST BANNER END -------

	fmt.Println("TestSignInAdmin - SUCCESS")
	fmt.Println("TestPostBanner - SUCCESS")

}

func TestPatchBanner(t *testing.T) {

	// ------ AUTH START -------
	signInURL := "http://localhost:8000/auth/sign-in/admin"
	inputData := map[string]string{
		"username": "testAdmin",
		"password": "testPassword",
	}
	jsonData, err := json.Marshal(inputData)
	if err != nil {
		t.Fatalf("Failed to marshal JSON data: %v", err)
	}

	resp, err := http.Post(signInURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}

	var responseToken struct {
		Token string `json:"token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&responseToken)
	if err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}
	token := responseToken.Token
	// ------ AUTH END -------

	// ------ PATCH BANNER START -------
	inputBanner := &types.BannerRequest{
		TagIds:    []int{1},
		FeatureId: 1,
		Content: types.Content{
			Title: "Test BannerUPDATE",
			Text:  "Test content",
			Url:   "http//example.com",
		},
		IsActive: true,
	}
	jsonData, err = json.Marshal(inputBanner)
	if err != nil {
		t.Fatalf("Failed to marshal JSON data: %v", err)
	}

	req, err := http.NewRequest("PATCH", "http://localhost:8000/banner/1", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create POST request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}

	type ResponseData struct {
		Description string `json:"description"`
	}

	var responseData string
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	if responseData != "OK" {
		t.Fatalf("Unexpected response: %s", responseData)
	}
	// ------ POST BANNER END -------

	fmt.Println("TestPatchBanner - SUCCESS")
}

func TestSignUpUser(t *testing.T) {
	serverURL := "http://localhost:8000/auth/sign-up"

	inputData := map[string]string{
		"username": "testUser",
		"password": "testPassword",
	}
	jsonData, err := json.Marshal(inputData)
	if err != nil {
		t.Fatalf("Failed to marshal JSON data: %v", err)
	}

	resp, err := http.Post(serverURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}

	var responseData struct {
		ID int `json:"id"`
	}
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	fmt.Println("TestSignUpUser -  SUCCESS")
}

func TestGetUserBanner(t *testing.T) {

	// ------ AUTH START -------
	signInuRL := "http://localhost:8000/auth/sign-in"
	inputData := map[string]string{
		"username": "testUser",
		"password": "testPassword",
	}
	jsonData, err := json.Marshal(inputData)
	if err != nil {
		t.Fatalf("Failed to marshal JSON data: %v", err)
	}

	resp, err := http.Post(signInuRL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}

	var responseData struct {
		Token string `json:"token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}
	token := responseData.Token
	// ------ AUTH END -------

	// ------ GET USER BANNER START -------
	serverURL := "http://localhost:8000"
	tagID := 1
	featureID := 1
	useLastRevision := "false"

	getUserBannerUrl := fmt.Sprintf("%s/user_banner?tag_id=%d&feature_id=%d&use_last_revision=%s", serverURL, tagID, featureID, useLastRevision)

	req, err := http.NewRequest("GET", getUserBannerUrl, nil)
	if err != nil {
		t.Fatalf("Failed to create GET request: %v", err)
	}

	authToken := token
	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}

	var banner types.Content
	err = json.NewDecoder(resp.Body).Decode(&banner)
	if err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	expectedBanner := types.Content{
		Title: "Test BannerUPDATE",
		Text:  "Test content",
		Url:   "http//example.com",
	}

	if banner != expectedBanner {
		t.Fatalf("Unexpected banner data: %+v", banner)
	}
	// ------ GET USER BANNER END -------

	fmt.Println("TestSignInUser - SUCCESS")
	fmt.Println("TestGetUserBanner - SUCCESS")
}

func TestDeleteBanner(t *testing.T) {

	// ------ AUTH START -------
	signInURL := "http://localhost:8000/auth/sign-in/admin"
	inputData := map[string]string{
		"username": "testAdmin",
		"password": "testPassword",
	}
	jsonData, err := json.Marshal(inputData)
	if err != nil {
		t.Fatalf("Failed to marshal JSON data: %v", err)
	}

	resp, err := http.Post(signInURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}

	var responseToken struct {
		Token string `json:"token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&responseToken)
	if err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}
	token := responseToken.Token
	// ------ AUTH END -------

	// ------ DELETE BANNER START -------
	req, err := http.NewRequest("DELETE", "http://localhost:8000/banner/1", nil)
	if err != nil {
		t.Fatalf("Failed to create DELETE request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send DELETE request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}
	// ------ DELETE BANNER END -------

	fmt.Println("TestPatchBanner - SUCCESS")
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	tenantID     = "YOUR_TENANT_ID"
	clientID     = "YOUR_CLIENT_ID"
	clientSecret = "YOUR_CLIENT_SECRET"
)

// アクセストークンを取得
func getAccessToken() (string, error) {
	url := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID)

	data := "grant_type=client_credentials" +
		"&client_id=" + clientID +
		"&client_secret=" + clientSecret +
		"&scope=https://graph.microsoft.com/.default"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	token, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("アクセストークン取得エラー: %v", result)
	}

	return token, nil
}

// Teams会議を作成
func createTeamsMeeting(accessToken string) (string, error) {
	url := "https://graph.microsoft.com/v1.0/me/onlineMeetings"

	meetingData := map[string]interface{}{
		"startDateTime": "2025-04-03T12:00:00Z",
		"endDateTime":   "2025-04-03T13:00:00Z",
		"subject":       "Test Meeting",
	}

	jsonData, _ := json.Marshal(meetingData)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if joinURL, ok := result["joinUrl"].(string); ok {
		return joinURL, nil
	}

	return "", fmt.Errorf("会議URLの取得に失敗しました: %v", result)
}

// メイン関数
func main() {
	token, err := getAccessToken()
	if err != nil {
		fmt.Println("アクセストークン取得エラー:", err)
		os.Exit(1)
	}

	joinURL, err := createTeamsMeeting(token)
	if err != nil {
		fmt.Println("会議作成エラー:", err)
		os.Exit(1)
	}

	fmt.Println("Teams会議URL:", joinURL)
}

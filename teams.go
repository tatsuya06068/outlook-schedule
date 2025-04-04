package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

// Microsoft Graph API エンドポイント
const (
	graphAPIURL       = "https://graph.microsoft.com/v1.0"
	tokenEndpoint     = "https://login.microsoftonline.com/%s/oauth2/v2.0/token"
	clientID          = "YOUR_CLIENT_ID"
	clientSecret      = "YOUR_CLIENT_SECRET"
	tenantID          = "YOUR_TENANT_ID"
	scope             = "https://graph.microsoft.com/.default" // クライアントクレデンシャルフローで必要なスコープ
)

// アクセストークンを取得する関数
func getAccessToken() (string, error) {
	client := resty.New()

	// POSTリクエストのボディを設定
	data := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"scope":         scope,
		"grant_type":    "client_credentials",
	}

	// リクエストを送信
	resp, err := client.R().
		SetFormData(data).
		Post(fmt.Sprintf(tokenEndpoint, tenantID))

	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return "", err
	}

	// アクセストークンを取得
	accessToken := result["access_token"].(string)
	return accessToken, nil
}

// Teams会議のURLを作成する関数
func createOnlineMeeting(accessToken string) error {
	// 会議データの作成
	meetingRequest := map[string]interface{}{
		"startDateTime": "2025-05-01T12:00:00Z",
		"endDateTime":   "2025-05-01T13:00:00Z",
		"subject":       "Sample Meeting",
	}

	// リクエストをJSONに変換
	meetingJSON, err := json.Marshal(meetingRequest)
	if err != nil {
		return fmt.Errorf("error marshalling meeting data: %v", err)
	}

	// 会議の作成リクエスト
	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+accessToken).
		SetHeader("Content-Type", "application/json").
		SetBody(meetingJSON).
		Post(fmt.Sprintf("%s/me/onlineMeetings", graphAPIURL))

	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}

	// レスポンスの解析
	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return fmt.Errorf("error unmarshalling response: %v", err)
	}

	// 会議の参加URLを取得
	joinURL := result["joinUrl"].(string)
	fmt.Println("Meeting URL: ", joinURL)

	return nil
}

func main() {
	// アクセストークンの取得
	accessToken, err := getAccessToken()
	if err != nil {
		log.Fatalf("Error getting access token: %v", err)
	}

	// Teams会議のURLを作成
	if err := createOnlineMeeting(accessToken); err != nil {
		log.Fatalf("Error creating online meeting: %v", err)
	}
}

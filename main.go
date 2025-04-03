package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Microsoft Graph API のエンドポイント
const GraphAPIBaseURL = "https://graph.microsoft.com/v1.0"

// アクセストークン（環境変数から取得）
var AccessToken = os.Getenv("MS_GRAPH_ACCESS_TOKEN")

// Outlook の予定を取得（直近の予定を取得）
func getEvents() error {
	url := GraphAPIBaseURL + "/me/events?$orderby=start/dateTime&$top=5"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("取得エラー: %s", string(body))
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)
	events := result["value"].([]interface{})

	fmt.Println("📅 取得した予定:")
	for _, e := range events {
		event := e.(map[string]interface{})
		fmt.Printf("- %s (%s - %s)\n", event["subject"], event["start"].(map[string]interface{})["dateTime"], event["end"].(map[string]interface{})["dateTime"])
	}

	return nil
}

// Outlook の予定を作成
func createEvent() (string, error) {
	url := GraphAPIBaseURL + "/me/events"

	event := map[string]interface{}{
		"subject":  "Goで作成した会議",
		"body": map[string]string{
			"contentType": "HTML",
			"content":     "この会議はGoコードから作成されました。",
		},
		"start": map[string]string{
			"dateTime": "2024-04-10T10:00:00",
			"timeZone": "Asia/Tokyo",
		},
		"end": map[string]string{
			"dateTime": "2024-04-10T11:00:00",
			"timeZone": "Asia/Tokyo",
		},
		"location": map[string]string{
			"displayName": "オンライン会議",
		},
	}

	eventData, _ := json.Marshal(event)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(eventData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+AccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("作成エラー: %s", string(body))
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)
	eventID := result["id"].(string)

	fmt.Println("✅ 予定作成成功！ID:", eventID)
	return eventID, nil
}

// Outlook の予定を更新
func updateEvent(eventID string) error {
	url := GraphAPIBaseURL + "/me/events/" + eventID

	updateData := map[string]interface{}{
		"subject": "【更新済み】Goで作成した会議",
	}
	updateJSON, _ := json.Marshal(updateData)

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(updateJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+AccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("更新エラー: %s", string(body))
	}

	fmt.Println("✅ 予定更新成功！")
	return nil
}

// Outlook の予定を削除
func deleteEvent(eventID string) error {
	url := GraphAPIBaseURL + "/me/events/" + eventID

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("削除エラー: %s", string(body))
	}

	fmt.Println("✅ 予定削除成功！")
	return nil
}

func main() {
	// ① 予定の取得
	err := getEvents()
	if err != nil {
		fmt.Println("予定の取得に失敗:", err)
		return
	}

	// ② 予定の登録
	eventID, err := createEvent()
	if err != nil {
		fmt.Println("予定の作成に失敗:", err)
		return
	}

	// ③ 予定の更新
	err = updateEvent(eventID)
	if err != nil {
		fmt.Println("予定の更新に失敗:", err)
		return
	}

	// ④ 予定の削除
	err = deleteEvent(eventID)
	if err != nil {
		fmt.Println("予定の削除に失敗:", err)
		return
	}
}

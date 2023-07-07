package main

import (
	"encoding/json"
	"fmt"
	"time"
)

func GetLastPushEvents(user string) ([]GithubEvent, error) {
	headers := map[string]string{
		"Accept": "application/vnd.github.v3+json",
	}
	url := fmt.Sprintf("https://api.github.com/users/%s/events/public", user)
	res, err := SendRequest("GET", url, nil, headers)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var jsonData []GithubEvent
	err = json.NewDecoder(res.Body).Decode(&jsonData)
	if err != nil {
		return nil, err
	}

	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		return nil, err
	}

	now := time.Now().In(loc)
	yesterday := now.AddDate(0, 0, -1)
	today := now.AddDate(0, 0, 0)

	var pushEvents []GithubEvent
	for _, e := range jsonData {
		createdAt, err := time.Parse(time.RFC3339, e.CreatedAt)
		if err != nil {
			return nil, err
		}
		if e.Type == "PushEvent" && createdAt.After(yesterday) && createdAt.Before(today) {
			pushEvents = append(pushEvents, e)
		}
	}
	return pushEvents, nil
}

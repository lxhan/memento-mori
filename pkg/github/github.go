package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Actor struct {
	Id           int    `json:"id"`
	Login        string `json:"login"`
	DisplayLogin string `json:"display_login"`
	GravatarId   string `json:"gravatar_id"`
	Url          string `json:"url"`
	AvatarUrl    string `json:"avatart_url"`
}

type Repo struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Payload struct {
	PushId       int    `json:"push_id"`
	Size         int    `json:"size"`
	DistinctSize int    `json:"distinct_size"`
	Ref          string `json:"ref"`
	Head         string `json:"head"`
	Before       string `json:"before"`
}

type Author struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Commit struct {
	Sha      string `json:"sha"`
	Author   Author `json:"author"`
	Message  string `json:"message"`
	Distinct bool   `json:"distinct"`
	Url      string `json:"url"`
}

type GithubEvent struct {
	Id        string   `json:"id"`
	Type      string   `json:"type"`
	Actor     Actor    `json:"actor"`
	Repo      Repo     `json:"repo"`
	Commits   []Commit `json:"commits"`
	Public    bool     `json:"public"`
	CreatedAt string   `json:"created_at"`
}

func GetTodayPushEvents(user string, limit int) ([]GithubEvent, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/users/%s/events/public?per_page=%d", user, limit), nil)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	req.Header.Add("Accept", "application/vnd.github.v3+json")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	defer res.Body.Close()

	var jsonData []GithubEvent
	decoder := json.NewDecoder(res.Body)

	err = decoder.Decode(&jsonData)

	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	return jsonData, nil
}
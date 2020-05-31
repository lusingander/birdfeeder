package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type postDetail struct {
	Number         int       `json:"number"`
	Name           string    `json:"name"`
	FullName       string    `json:"full_name"`
	Wip            bool      `json:"wip"`
	BodyMd         string    `json:"body_md"`
	CreatedAt      time.Time `json:"created_at"`
	Message        string    `json:"message"`
	URL            string    `json:"url"`
	UpdatedAt      time.Time `json:"updated_at"`
	Tags           []string  `json:"tags"`
	Category       string    `json:"category"`
	RevisionNumber int       `json:"revision_number"`
	CreatedBy      struct {
		Name       string `json:"name"`
		ScreenName string `json:"screen_name"`
		Icon       string `json:"icon"`
	} `json:"created_by"`
	UpdatedBy struct {
		Name       string `json:"name"`
		ScreenName string `json:"screen_name"`
		Icon       string `json:"icon"`
	} `json:"updated_by"`
}

func (p *postDetail) print() {
	fmt.Println(p.Number, p.Name, p.CreatedAt)
}

type postsResponse struct {
	Posts      []*postDetail `json:"posts"`
	PrevPage   int           `json:"prev_page"`
	NextPage   int           `json:"next_page"`
	TotalCount int           `json:"total_count"`
	Page       int           `json:"page"`
	PerPage    int           `json:"per_page"`
	MaxPerPage int           `json:"max_per_page"`
}

func (r *postsResponse) print() {
	fmt.Println(r.PrevPage, r.NextPage, r.TotalCount, r.Page, r.PerPage, r.MaxPerPage)
	for _, p := range r.Posts {
		p.print()
	}
}

func buildGetPostsURL(page int, team, token string) string {
	perPage := 50
	base := fmt.Sprintf("https://api.esa.io/v1/teams/%s/posts", team)
	query := fmt.Sprintf("page=%d&per_page=%d&sort=created&order=asc&access_token=%s", page, perPage, token)
	return base + "?" + query
}

func fetchPosts(page int, team, token string) (*postsResponse, error) {
	url := buildGetPostsURL(page, team, token)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var posts postsResponse
	err = json.Unmarshal(bytes, &posts)
	if err != nil {
		return nil, err
	}

	return &posts, nil
}

func fetchAllPosts(cfg *config) ([]*postDetail, error) {
	posts := make([]*postDetail, 0)
	page := 1
	for {
		postsRes, err := fetchPosts(page, cfg.Team, cfg.Token)
		if err != nil {
			return nil, err
		}
		posts = append(posts, postsRes.Posts...)
		if postsRes.NextPage == 0 {
			break
		}
		page = postsRes.NextPage
	}
	return posts, nil
}

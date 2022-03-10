package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const baseUrl = "https://node-hnapi.herokuapp.com"
const postResource = "item"
const march22PostId = "30515750"

func main() {
	resp, err := http.Get(baseUrl + "/" + postResource + "/" + march22PostId)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

type PostResourceResponse struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Points   int    `json:"points"`
	User     string `json:"user"`
	Time     int    `json:"time"`
	TimeAgo  string `json:"time_ago"`
	Type     string `json:"type"`
	Content  string `json:"content"`
	URL      string `json:"url"`
	Comments []struct {
		ID       int           `json:"id"`
		Level    int           `json:"level"`
		User     string        `json:"user,omitempty"`
		Time     int           `json:"time"`
		TimeAgo  string        `json:"time_ago"`
		Content  string        `json:"content,omitempty"`
		Comments []interface{} `json:"comments"`
		Deleted  bool          `json:"deleted,omitempty"`
		Dead     bool          `json:"dead,omitempty"`
	} `json:"comments"`
	CommentsCount int `json:"comments_count"`
}

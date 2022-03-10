package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

const baseUrl = "https://node-hnapi.herokuapp.com"
const postResource = "item"
const march22PostId = "30515750"

var disqualifyRegEx = regexp.MustCompile(`(?i)(nft | chain | web3 | bitcoin | decentralize | democratize | onsite)`)

type Post struct {
	ID            int           `json:"id"`
	Title         string        `json:"title"`
	Points        int           `json:"points"`
	User          string        `json:"user"`
	Time          int           `json:"time"`
	TimeAgo       string        `json:"time_ago"`
	Type          string        `json:"type"`
	Content       string        `json:"content"`
	URL           string        `json:"url"`
	Comments      []PostComment `json:"comments"`
	CommentsCount int           `json:"comments_count"`
}

type PostComment struct {
	ID       int           `json:"id"`
	Level    int           `json:"level"`
	User     string        `json:"user,omitempty"`
	Time     int           `json:"time"`
	TimeAgo  string        `json:"time_ago"`
	Content  string        `json:"content,omitempty"`
	Comments []interface{} `json:"comments"`
	Deleted  bool          `json:"deleted,omitempty"`
	Dead     bool          `json:"dead,omitempty"`
}

func main() {
	resp, err := http.Get(baseUrl + "/" + postResource + "/" + march22PostId)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var result Post

	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	filteredSubmissions := removeInvalidEntries(result.Comments)

	fmt.Println(filteredSubmissions)

}

func removeInvalidEntries(comments []PostComment) []PostComment {
	filtered := make([]PostComment, 0)

	for _, comment := range comments {
		hasMatch, err := regexp.MatchString(`(?i)(nft | chain | web3 | bitcoin | decentralize | democratize | onsite)`, comment.Content)

		if err != nil {
			panic(err)
		}

		if !hasMatch {
			filtered = append(filtered, comment)
		}
	}
	return filtered
}

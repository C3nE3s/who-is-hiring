package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

//can get post ID through user call in future to automate rather than hardcode values

const baseUrl = "https://node-hnapi.herokuapp.com"
const postResource = "item"
const userResource = "user"
const march22PostId = "30515750"

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

type Listing struct {
	Title       string
	Description string
	Links       string
	Time        int
	Score       int
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

	listings := transformAndRank(result.Comments)
	writeToCSV(listings)
}

// this is On^2 (nested loops)
// TODO: make this more efficient
func transformAndRank(comments []PostComment) []Listing {
	final := make([]Listing, 0)

	for _, comment := range comments {
		temp := transformRawListing(comment.Content)
		temp.Time = comment.Time
		temp.Score = getListingRelevanceRank(comment.Content)
		final = append(final, temp)
	}

	return final
}

func getListingRelevanceRank(listing string) int {

	regexMap := map[string]*regexp.Regexp{
		//any 6 digit number starting w/ 1 or 2; 1 or 2 w/ 2 digits after, preceded by us $ sign
		"compensation":     regexp.MustCompile(`(?mi)([12]\d{2},\d{3})|([$][1,2]\d{2})`),
		"typescript":       regexp.MustCompile(`(?mi)typescript`),
		"remote":           regexp.MustCompile(`(?im)remote`),
		"location":         regexp.MustCompile(`(?m)([nN]orth(?:[[:blank:]])[aA](?:merica))|(US)|([sS](?:tates))`),
		"desired_tech":     regexp.MustCompile(`(?mi)(react)|(nextjs)|(gatsby)|(nuxt)|(svelte)`),
		"front_end":        regexp.MustCompile(`(?im)front(?:[\se-])`),
		"full_stack":       regexp.MustCompile(`(?mi)(full-stack)|(fullstack)|(full stack)`),
		"acceptable_tech":  regexp.MustCompile(`(?mi)(angular)|(vue)`),
		"luxury":           regexp.MustCompile(`(?mi)(unlimited)|(patern*)|(parental)|(4[/s-]day)|(family[/s-]friendly)`),
		"seniority":        regexp.MustCompile(`(?mi)(senior)|(mid)`),
		"precise_location": regexp.MustCompile(`(?m)([G][Aa])|([Aa]tlanta)`),
		"undesired_tech":   regexp.MustCompile(`(?mi)(ember)|(jquery)|(angularjs)|(wordpress)`),
		//if the listing mentions onsite, I punish more for a prefix or post fix 'and', 'or' hints at a hybrid model
		"onsite_soft": regexp.MustCompile(`(?mi)((?:(or)) onsite)|(onsite (?:or))`),
		"onsite_hard": regexp.MustCompile(`(?mi)((?:(and)) onsite)|(onsite (?:and))`),
		"disqualify":  regexp.MustCompile(`(?mi)(nft | blockchain | web3 | bitcoin | decentralize | democratize )`),
	}

	scoreMap := map[string]int{
		"compensation":     5,
		"typescript":       5,
		"remote":           5,
		"location":         5,
		"desired_tech":     4,
		"front_end":        4,
		"full_stack":       3,
		"acceptable_tech":  3,
		"luxury":           3,
		"seniority":        2,
		"precise_location": 1,
		"undesired_tech":   -1,
		"onsite_soft":      -3,
		"onsite_hard":      -5,
		"disqualify":       -45,
	}

	listingScore := 0

	for key, value := range regexMap {
		if value.MatchString(listing) {
			listingScore += scoreMap[key]
		}
	}
	return listingScore

}

func transformRawListing(listing string) Listing {
	splitLisiting := Listing{}
	reader := strings.NewReader(listing)
	tokenizer := html.NewTokenizer(reader)

	prevStartToken := tokenizer.Token()

	for {
		tt := tokenizer.Next()
		switch {
		case tt == html.ErrorToken:
			if tokenizer.Err() == io.EOF {
				return splitLisiting
			} else {
				fmt.Println(html.ErrorToken.String())
			}
			break
		case tt == html.StartTagToken:
			prevStartToken = tokenizer.Token()
		case tt == html.TextToken:
			if prevStartToken.Data != "a" {
				// first non <a> node will always be title
				if reflect.DeepEqual(splitLisiting, StructuredListing{}) {
					splitLisiting.Title = transformTokenToText(tokenizer.Text())
				} else {
					splitLisiting.Description += transformTokenToText(tokenizer.Text()) + " "
				}
			} else {
				for _, a := range prevStartToken.Attr {
					if a.Key == "href" {
						splitLisiting.Links += a.Val + " "
						break
					}
				}

			}
		}
	}
}

func transformTokenToText(token []byte) string {
	return strings.TrimSpace(html.UnescapeString(string(token)))
}

func writeToCSV(listings []Listing) {
	return
}

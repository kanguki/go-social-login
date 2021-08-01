package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

type Facebook struct {
	config FacebookConfiguration
}

type FacebookConfiguration struct {
	state  string
	oauth2 *oauth2.Config
}

func NewFacebook(state string, oauth2 *oauth2.Config) Facebook {
	return Facebook{
		config: FacebookConfiguration{
			state:  state,
			oauth2: oauth2,
		},
	}
}

func (f Facebook) HandleLogin(w http.ResponseWriter, r *http.Request) {
	url := f.config.oauth2.AuthCodeURL(f.config.state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
func (f Facebook) HandleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != f.config.state {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", f.config.state, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")

	token, err := f.config.oauth2.Exchange(r.Context(), code)
	if err != nil {
		fmt.Printf("exchange token with code failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	accessToken := url.QueryEscape(token.AccessToken)
	request := "https://graph.facebook.com/me" +
		"?fields=name,birthday,email,hometown,picture,age_range,gender,link,quotes,feed,photos,friends" +
		"&access_token=" + accessToken
	log.Printf("request: %v", request)
	resp, err := http.Get(request)
	if err != nil {
		fmt.Printf("Get: %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ReadAll: %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	var fbUserData FbUserData
	err = json.Unmarshal(response, &fbUserData)
	if err != nil {
		fmt.Printf("Error login facebook: Unmarshal error: %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	log.Printf("parseResponseBody: %+v\naccessToken: %s", fbUserData, accessToken)

	http.Redirect(w, r, "/authenticated", http.StatusTemporaryRedirect)

}

type FbUserData struct {
	Id       string //used to query users' data. refer to graph api reference for more details
	Name     string
	Birthday string `format:"dd/mm/yyyy"`
	Email    string
	Picture  struct {
		Data struct {
			Url string
		}
	}
	AgeRange    FbAgeRange `json:"age_range"`
	Gender      string
	LinkProfile string `json:"link"`
	Quotes      string
	Feed        struct {
		Data   []FbFeedDatum
		Paging FbFeedPaging
	}
	Friends struct {
		Data    []FbFriendsDatum
		Summary FbFriendsSummary
	}
}

type FbAgeRange struct {
	Max int
	Min int
}
type FbFeedDatum struct {
	CreatedTime string `json:"created_time" format:"yyyy-MM-ddTHH:mm:ss+timeZone"`
	Id          string `format:"feedId_randomBigInt"`
}
type FbFeedPaging struct {
	PreviousFeed string `json:"previous"`
	NextFeed     string `json:"next"`
}
type FbFriendsDatum struct {
}
type FbFriendsSummary struct {
	TotalCount int `json:"total_count"`
}

package handlers

import (
	// "encoding/json"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

type Google struct {
	config GoogleConfiguration
}

type GoogleConfiguration struct {
	state  string
	oauth2 *oauth2.Config
}

func NewGoogle(state string, oauth2 *oauth2.Config) Google {
	return Google{
		config: GoogleConfiguration{
			state:  state,
			oauth2: oauth2,
		},
	}
}

func (g Google) HandleLogin(w http.ResponseWriter, r *http.Request) {
	url := g.config.oauth2.AuthCodeURL(g.config.state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
func (g Google) HandleCallback(w http.ResponseWriter, r *http.Request) {

	state := r.FormValue("state")
	if state != g.config.state {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", g.config.state, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")

	token, err := g.config.oauth2.Exchange(r.Context(), code)
	if err != nil {
		fmt.Printf("exchange token with code failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	accessToken := url.QueryEscape(token.AccessToken)
	request := "https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken
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
	var ggUserData GoogleUserData
	err = json.Unmarshal(response, &ggUserData)
	if err != nil {
		fmt.Printf("Error login facebook: Unmarshal error: %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	log.Printf("parseResponseBody: %+v\naccessToken: %s", ggUserData, accessToken)

	http.Redirect(w, r, "/authenticated", http.StatusTemporaryRedirect)

}

type GoogleUserData struct {
	Id string
	Email string
	VerifiedEmail bool `json:"verified_email"`
	DisplayName string `json:"name"`
	GivenName string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Avatar string `json:"picture"`
	Locale string
}
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
	"mo.io/goLogin/internal/handlers"
)

var protocol, port, serverAddress, ssl_cert, ssl_key = getServerInfo()

//facebook
var (
	fbAppId       = os.Getenv("FB_APP_ID")
	fbAppSecret   = os.Getenv("FB_APP_SECRET")
	fbRedirectUri = os.Getenv("FB_REDIRECT_URI")
	fbAuthState   = os.Getenv("FB_AUTH_STATE")
	fbOauth       = &oauth2.Config{
		ClientID:     fbAppId,
		ClientSecret: fbAppSecret,
		RedirectURL:  serverAddress + fbRedirectUri,
		Scopes: []string{"public_profile", "email", "user_birthday",
			"user_friends", "user_gender", "user_link", "user_photos", "user_hometown"},
		Endpoint: facebook.Endpoint,
	}
	f = handlers.NewFacebook(fbAuthState, fbOauth)
)

//google
var (
	ggAppId       = os.Getenv("GG_APP_ID")
	ggAppSecret   = os.Getenv("GG_APP_SECRET")
	ggRedirectUri = os.Getenv("GG_REDIRECT_URI")
	ggAuthState   = os.Getenv("GG_AUTH_STATE")
	ggOauth       = &oauth2.Config{
		ClientID:     ggAppId,
		ClientSecret: ggAppSecret,
		RedirectURL:  serverAddress + ggRedirectUri,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/plus.me"},
		Endpoint: google.Endpoint,
	}
	g = handlers.NewGoogle(ggAuthState, ggOauth)
)

func main() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/authenticated", handleAuthenticated)

	//fb
	http.HandleFunc("/login/facebook", f.HandleLogin)
	http.HandleFunc(fbRedirectUri, f.HandleCallback)

	//gg
	http.HandleFunc("/login/google", g.HandleLogin)
	http.HandleFunc(ggRedirectUri, g.HandleCallback)

	fmt.Print("Started running on " + serverAddress + "\n")
	if protocol == "https" {
		log.Fatal(http.ListenAndServeTLS(port, ssl_cert, ssl_key, nil))
	}
	log.Fatal(http.ListenAndServe(port, nil))
}

func getServerInfo() (protocol, port, serverAddress, ssl_cert, ssl_key string) {
	host := Getenv("HOST", "localhost")
	ssl_cert = os.Getenv("SSL_CERT")
	ssl_key = os.Getenv("SSL_KEY")
	if ssl_cert != "" && ssl_key != "" {
		protocol = "https"
	}
	port = os.Getenv("PORT")
	if port == "" {
		log.Fatal("Missing port")
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	serverAddress = protocol + "://" + host + port
	return protocol, port, serverAddress, ssl_cert, ssl_key
}

func Getenv(envKey, defaultVal string) string {
	key := os.Getenv(envKey)
	if key != "" {
		return key
	}
	return defaultVal
}

const htmlIndex = `<html><body>
Logged in with <a href="/login">facebook</a>
</body></html>
`

func handleMain(w http.ResponseWriter, r *http.Request) {
	index, err := ioutil.ReadFile("/home/mo/goLoginFb/asset/index.html")
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(htmlIndex))
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(index)
}

func handleAuthenticated(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`Tada, you are authenticated`))
}

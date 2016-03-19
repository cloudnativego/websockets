package server

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/astaxie/beego/session"
)

var chatTemplate *template.Template

type chatUser struct {
	Email      string
	PictureURL string
	PubKey     string
	SubKey     string
}

func chatHandler(sessionManager *session.Manager, messagingConfig *pubsubConfig) http.HandlerFunc {
	chatTemplate = template.Must(template.ParseFiles("assets/templates/chat.html"))

	return func(w http.ResponseWriter, r *http.Request) {

		session, _ := sessionManager.SessionStart(w, r)
		defer session.SessionRelease(w)

		// Getting the profile from the session
		profile, _ := session.Get("profile").(map[string]interface{})
		fmt.Printf("%+v", profile)
		fmt.Printf("email: %s, picture: %s\n", profile["email"].(string), profile["picture"].(string))
		data := chatUser{Email: profile["email"].(string), PictureURL: profile["picture"].(string), PubKey: messagingConfig.PublishKey, SubKey: messagingConfig.SubscribeKey}

		chatTemplate.Execute(w, data)
	}
}

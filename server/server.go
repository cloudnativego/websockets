package server

import (
	"net/http"

	"github.com/astaxie/beego/session"
	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/cloudnativego/cf-tools"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

type authConfig struct {
	ClientID     string
	ClientSecret string
	Domain       string
	CallbackURL  string
}

type pubsubConfig struct {
	PublishKey   string
	SubscribeKey string
}

// NewServer configures and returns a Server.
func NewServer(appEnv *cfenv.App) *negroni.Negroni {

	authClientID, _ := cftools.GetVCAPServiceProperty("authzero", "id", appEnv)
	authSecret, _ := cftools.GetVCAPServiceProperty("authzero", "secret", appEnv)
	authDomain, _ := cftools.GetVCAPServiceProperty("authzero", "domain", appEnv)
	authCallback, _ := cftools.GetVCAPServiceProperty("authzero", "callback", appEnv)

	subKey, _ := cftools.GetVCAPServiceProperty("pubnub", "subkey", appEnv)
	pubKey, _ := cftools.GetVCAPServiceProperty("pubnub", "pubkey", appEnv)

	messagingConfig := &pubsubConfig{
		PublishKey:   pubKey,
		SubscribeKey: subKey,
	}

	config := &authConfig{
		ClientID:     authClientID,
		ClientSecret: authSecret,
		Domain:       authDomain,
		CallbackURL:  authCallback,
	}

	sessionManager, _ := session.NewManager("memory", `{"cookieName":"gosessionid","gclifetime":3600}`)
	go sessionManager.GC()

	n := negroni.Classic()
	mx := mux.NewRouter()

	initRoutes(mx, sessionManager, config, messagingConfig)

	n.UseHandler(mx)
	return n
}

func initRoutes(mx *mux.Router, sessionManager *session.Manager, config *authConfig, messagingConfig *pubsubConfig) {
	mx.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("./assets/images/"))))
	mx.HandleFunc("/", homeHandler(config))
	mx.HandleFunc("/callback", callbackHandler(sessionManager, config))
	mx.Handle("/user", negroni.New(
		negroni.HandlerFunc(isAuthenticated(sessionManager)),
		negroni.Wrap(http.HandlerFunc(userHandler(sessionManager))),
	))
	mx.Handle("/chat", negroni.New(
		negroni.HandlerFunc(isAuthenticated(sessionManager)),
		negroni.Wrap(http.HandlerFunc(chatHandler(sessionManager, messagingConfig))),
	))
	mx.HandleFunc("/broadcast", broadcastHandler(messagingConfig)).Methods("POST")
	mx.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public/"))))
}

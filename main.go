package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/auth0/go-jwt-middleware"
	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/odise/moxy"
)

func readConfig(path string) (*Config, error) {

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	config, err := ParseConfig(string(buf))
	if err != nil {
		return nil, err
	}

	return config, nil
}

// SetRequestID sets request ID as response header
func SetRequestID(request *http.Request, response *http.Response, id string) {
	response.Header.Set("X-RevProx-Id", id)
}

var port string
var configFile string

func init() {
	const (
		defaultport       = ":3009"
		defaultConfigFile = "./config"
	)
	flag.StringVar(&port, "port", defaultport, "Port the server is listening at")
	flag.StringVar(&port, "p", defaultport, "Port the server is listening at (shorthand)")
	flag.StringVar(&configFile, "config", defaultConfigFile, "Config file")
	flag.StringVar(&configFile, "c", defaultConfigFile, "Config file (shorthand)")

}

func main() {

	flag.Parse()

	cfg, err := readConfig(configFile)
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()

	for _, p := range cfg.Proxies {
		createProxies(&p, router)
	}

	app := negroni.New()
	app.UseHandler(router)
	if port[0] != ':' {
		port = ":" + port
	}
	app.Run(port)
}

var myHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user")
	log.Printf("revprox: This is an authenticated request")
	log.Printf("revprox: Claim content:\n")
	for k, v := range user.(*jwt.Token).Claims {
		log.Printf("revprox: %s :\t%#v\n", k, v)
	}
})

func addHost(r *mux.Router, hostname string) *mux.Router {
	s := r.Host(hostname).Subrouter()
	return s
}

func addSecureRoute(router *mux.Router, proxy *moxy.ReverseProxy, route string, jwtMiddleware *jwtmiddleware.JWTMiddleware) {
	router.Handle(route, negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(myHandler),
		negroni.HandlerFunc(proxy.HandlerWithNext),
	))
}

func addPublicRoute(router *mux.Router, proxy *moxy.ReverseProxy, route string) {
	router.Handle(route, negroni.New(
		negroni.HandlerFunc(proxy.HandlerWithNext),
	))
}

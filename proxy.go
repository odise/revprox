package main

import (
	"log"

	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/odise/moxy"
)

func createProxies(config *ProxyConfig, router *mux.Router) {

	filters := []moxy.FilterFunc{SetRequestID}
	proxy := moxy.NewReverseProxy(config.Targets, filters)
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Secret), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
		ErrorHandler:  login,
	})

	host := addHost(router, config.Name)
	for _, p := range config.PublicPath {
		log.Printf("revprox: setting public route %s to %s", p, config.Name)
		addPublicRoute(host, proxy, p)
	}

	for _, p := range config.AuthPath {
		log.Printf("revprox: setting secure route %s to %s", p, config.Name)
		addSecureRoute(host, proxy, p, jwtMiddleware)
	}

}

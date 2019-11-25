package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/zhughes3/website/middleware"
	"github.com/zhughes3/website/router"
	"github.com/zhughes3/website/user"
	"net/http"
	"os"
)

func NewRouter() http.Handler {
	appRouter := mux.NewRouter()


	for _, r := range getRoutes() {
		var handler http.Handler
		handler = r.HandlerFunc
		if r.Protected {
			handler = middleware.JWTMiddleware(r.HandlerFunc)
		}
		appRouter.Handle(r.Pattern, handler).Methods(r.Method)
	}

	initFileServer(appRouter)
	loggedRouter := handlers.LoggingHandler(os.Stdout, appRouter)
	return loggedRouter
}

func getRoutes() []router.Route {
	var routes []router.Route
	routes = append(routes, user.Routes...)
	return routes
}

func initFileServer(r *mux.Router) {
	fs := http.FileServer(http.Dir("public/"))
	r.PathPrefix("/").Handler(http.StripPrefix("/", fs)).Methods("GET")
}

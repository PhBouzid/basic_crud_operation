package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type App struct{
	Router *mux.Router

}



func(a *App) Initialize(){

}

func commonMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func jwtVerify(next http.Handler) http.Handler{
	return nil
}

func initProjectRoutes(route *mux.Router){
	route.HandleFunc("/project", nil).Methods("GET","OPTIONS")
	route.HandleFunc("/project/{id}",nil).Methods("GET","OPTIONS")
	route.HandleFunc("/project",nil).Methods("POST","OPTIONS")
	route.HandleFunc("/project",nil).Methods("PUT","OPTIONS")
}

func initBacklogRoutes(route *mux.Router){
	route.HandleFunc("/project", nil).Methods("GET","OPTIONS")
	route.HandleFunc("/project/{id}",nil).Methods("GET","OPTIONS")
	route.HandleFunc("/project",nil).Methods("POST","OPTIONS")
	route.HandleFunc("/project",nil).Methods("PUT","OPTIONS")
}

func initUserStoryRoutes(route *mux.Router){

}

func initTaskRoutes(route *mux.Router){

}

func(a *App) run(addr string){
	a.Router.Use(commonMiddleware)
	subroute := a.Router.PathPrefix("/Scrum/API/v1").Subrouter()
	initProjectRoutes(subroute)
	initBacklogRoutes(subroute)
	initUserStoryRoutes(subroute)
	initTaskRoutes(subroute)
	http.Handle("/", a.Router)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(8000), nil))
}
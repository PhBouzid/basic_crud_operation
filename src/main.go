package main


import (
	"fmt"
	"github.com/gorilla/mux"
)

type Config struct{
	ConnectString string
	DatabaseName string

}

func main(){
	fmt.Println("Bookref Service started on port: " + " version: " )
	r := mux.NewRouter()
	r.HandleFunc("/doc/get",nil)
}



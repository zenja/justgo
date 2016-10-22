package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/zenja/justgo/controller"
	"github.com/zenja/justgo/utils"
)

func main() {
	// Load DB file
	boltFilename := flag.String("boltdb", "", "db file for boltdb")
	flag.Parse()
	if *boltFilename == "" {
		log.Fatalln("please specify boltdb file")
	}
	if err := utils.OpenDB(*boltFilename); err != nil {
		log.Fatalf("init DB failed: %s\n", err)
	}

	// TODO Add shutdown hook to close DB

	http.HandleFunc("/tutorial/new/", controller.AddTutorial)
	http.HandleFunc("/tutorial/save/", controller.SaveTutorial)
	http.HandleFunc("/tutorial/all/", controller.ListTutorials)
	http.HandleFunc("/tutorial/edit/", controller.EditTutorial)
	http.HandleFunc("/tutorial/delete/", controller.DeleteTutorial)
	http.ListenAndServe(":8080", nil)
}

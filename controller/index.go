package controller

import (
	"github.com/zenja/justgo/model"
	"github.com/zenja/justgo/template"
	"log"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	err := template.All.ExecuteTemplate(w, "test.html", &model.ExtendedTutorial{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

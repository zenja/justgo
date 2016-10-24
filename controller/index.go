package controller

import (
	"github.com/zenja/justgo/model"
	"net/http"
	"github.com/zenja/justgo/template"
	"log"
)

func Index(w http.ResponseWriter, r *http.Request) {
	err := template.All.ExecuteTemplate(w, "test.html", &model.Tutorial{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

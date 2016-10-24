package controller

import (
	"fmt"
	"github.com/zenja/justgo/model"
	"github.com/zenja/justgo/utils"
	"log"
	"net/http"
	"regexp"
	"github.com/zenja/justgo/template"
)

var tutorialValidPath = regexp.MustCompile("^/tutorial/(.*)/(.+)$")

func ListTutorials(w http.ResponseWriter, r *http.Request) {
	keys, err := utils.FetchAllKeys()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = template.All.ExecuteTemplate(w, "list-tutorial.html", keys)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func DeleteTutorial(w http.ResponseWriter, r *http.Request) {
	// Get key from /tutorial/delete/<key>
	m := tutorialValidPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}
	key := m[2]

	if err := utils.RemoveTutorial(key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/tutorial/all/", http.StatusFound)
}

func AddTutorial(w http.ResponseWriter, r *http.Request) {
	err := template.All.ExecuteTemplate(w, "add-tutorial.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func SaveTutorial(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	title := r.FormValue("title")
	description := r.FormValue("description")
	code := r.FormValue("code")
	expStdout := r.FormValue("expected_stdout")
	t := &model.Tutorial{
		Key:            key,
		Title:          title,
		Description:    description,
		Code:           code,
		ExpectedStdout: expStdout,
	}
	if err := utils.AddTutorial(t); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error in SaveTutorial: %s\n", err)
		return
	}
	http.Redirect(w, r, "/tutorial/all/", http.StatusFound)
}

func EditTutorial(w http.ResponseWriter, r *http.Request) {
	// Get key from /tutorial/edit/<key>
	m := tutorialValidPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}
	key := m[2]

	// Fetch tutorial from DB
	t, err := utils.FetchTutorial(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error in EditTutorial: %s\n", err)
		return
	}
	if t == nil {
		// Tutorial not found
		fmt.Fprintf(w, "Key %s not found.\n", key)
		http.NotFound(w, r)
		return
	}

	err = template.All.ExecuteTemplate(w, "edit-tutorial.html", t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error in EditTutorial: %s\n", err)
		return
	}
}

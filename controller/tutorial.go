package controller

import (
	"fmt"
	"github.com/zenja/justgo/model"
	"github.com/zenja/justgo/template"
	"github.com/zenja/justgo/utils"
	"log"
	"net/http"
	"regexp"
	"sort"
)

var tutorialValidPath = regexp.MustCompile("^/tutorial/(.*)/(.+)$")

func ListTutorials(w http.ResponseWriter, r *http.Request) {
	keys, err := utils.FetchAllKeys()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sort.Strings(keys)
	err = template.All.ExecuteTemplate(w, "list-tutorial.html", keys)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type tutorialPointerSlice []*model.Tutorial

func (s tutorialPointerSlice) Len() int {
	return len(s)
}

func (s tutorialPointerSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s tutorialPointerSlice) Less(i, j int) bool {
	return s[i].Key < s[j].Key
}

func TutorialsOverview(w http.ResponseWriter, r *http.Request) {
	tts, err := utils.FetchAllTutorials()
	if err != nil {
		http.Error(w, "failed to fetch all tutorials: "+err.Error(), http.StatusInternalServerError)
		return
	}
	sort.Sort(tutorialPointerSlice(tts))
	err = template.All.ExecuteTemplate(w, "tutorials-overview.html", tts)
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
	newKey := r.FormValue("new_key")
	originKey := r.FormValue("origin_key")
	title := r.FormValue("title")
	description := r.FormValue("description")
	code := r.FormValue("code")
	expStdout := r.FormValue("expected_stdout")
	if newKey != originKey && originKey != "" {
		if err := utils.RemoveTutorial(originKey); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error in SaveTutorial: failed to delete origin tutorial %s: %s\n", originKey, err)
			return
		}
	}
	t := &model.Tutorial{
		Key:            newKey,
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

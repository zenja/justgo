package controller

import (
	"fmt"
	"github.com/zenja/justgo/model"
	"github.com/zenja/justgo/template"
	"github.com/zenja/justgo/utils"
	"log"
	"net/http"
	"regexp"
)

var testerValidPath = regexp.MustCompile("^/test/(.+)$")

func Test(w http.ResponseWriter, r *http.Request) {
	// Get key from /test/<key>
	m := testerValidPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		log.Println("not matched") // fixme
		return
	}
	key := m[1]

	// Fetch tutorial from DB
	t, err := utils.FetchTutorial(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error when fetching key [%s]: %s\n", key, err)
		return
	}
	if t == nil {
		// Tutorial not found
		fmt.Fprintf(w, "Key %s not found.\n", key)
		http.NotFound(w, r)
		return
	}

	// Fetch the previous & next key
	prevKey, nextKey, _ := utils.FetchPreNextKey(key)
	et := &model.ExtendedTutorial{
		Tutorial: *t,
		PrevKey:  prevKey,
		NextKey:  nextKey,
	}

	err = template.All.ExecuteTemplate(w, "test.html", et)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error in executing template test.html: %s\n", err)
		return
	}
}

func Play(w http.ResponseWriter, r *http.Request) {
	err := template.All.ExecuteTemplate(w, "test.html", &model.ExtendedTutorial{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

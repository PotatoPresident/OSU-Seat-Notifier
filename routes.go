package main

import (
	"html/template"
	"net/http"
	"strings"
)

func index(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	courseCodes := strings.Split(query.Get("courseCodes"), ",")
	for _, code := range courseCodes {
		if code != "" {
			trackCourse(code)
		}
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	err := tmpl.Execute(w, trackedCourses)
	if err != nil {
		panic(err)
	}
}

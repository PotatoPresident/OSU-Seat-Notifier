package main

import (
	"embed"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

var (
	//go:embed all:templates/*
	templateFS embed.FS

	//go:embed css/output.css
	css embed.FS

	trackedCourses = map[string]CourseSearchResponse{}
)

func main() {
	go update()

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(ChangeMethod)

	router.Handle("/css/output.css", http.FileServer(http.FS(css)))

	router.HandleFunc("/", index)
	router.HandleFunc("/index.html", index)

	log.Println("Listening on http://localhost:8000")
	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Fatalln(err)
	}
}

func trackCourse(courseCode string) {
	if _, exists := trackedCourses[courseCode]; !exists {
		trackedCourses[courseCode] = CourseSearchResponse{}
	}
}

func ChangeMethod(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			switch method := r.PostFormValue("_method"); method {
			case http.MethodPut:
				fallthrough
			case http.MethodPatch:
				fallthrough
			case http.MethodDelete:
				r.Method = method
			default:
			}
		}
		next.ServeHTTP(w, r)
	})
}

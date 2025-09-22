package main

import (
	"html/template"
	"net/http"
)

type PageData struct {
	Message string
}

func main() {
	fs := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.Execute(w, PageData{})

	})
	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			prenom := r.FormValue("prenom")
			message := "Hello " + prenom + "!"
			tmpl :=
				template.Must(template.ParseFiles("templates/index.html"))
			tmpl.Execute(w, PageData{Message: message})

		}

	})
	http.ListenAndServe(":8080", nil)
}

package main

import (
	"html/template"
	"net/http"
	"strconv"
)

type PageData struct {
	Message string
}

var a [7][6]int
var turn bool

func main() {
	turn = false
	for _, val := range a {
		for i := range val {
			val[i] = 0
		}
	}
	fs := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.Execute(w, PageData{})

	})
	http.HandleFunc("/submit", submit)
	http.HandleFunc("/turn", DoTurn)
	http.ListenAndServe(":8080", nil)
}

func submit(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		prenom := r.FormValue("prenom")
		message := "Hello " + prenom + "!"
		tmpl :=
			template.Must(template.ParseFiles("templates/index.html"))
		tmpl.Execute(w, PageData{Message: message})

	}

}
func DoTurn(w http.ResponseWriter, r *http.Request) {
	x, err := strconv.Atoi(r.FormValue("x"))
	if err == nil {
		placepiece(x)
		return
	}
	print("error ntm touche pas a mon code connard de tes morts")

}
func placepiece(x int) { //false = red, true = yellow
	for i := 0; i <= 7; i++ {

		if a[x][i] == 0 {
			if turn {
				a[x][i] = 1
				turn = !turn
				return
			} else {
				a[x][i] = 2
				turn = !turn
				return
			}
		}
	}
	//pas pu placer le pion :(
}

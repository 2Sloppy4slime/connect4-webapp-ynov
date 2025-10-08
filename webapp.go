package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

var lastrow int
var lastcol int

type PageData struct {
	Message string
}

type Color struct {
	Color string
}

var a [7][6]int
var turn bool //false = yellow, true = red

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
	http.HandleFunc("/color", givecolor)
	http.HandleFunc("/pos", getpos)
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
func getpos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, strconv.Itoa(lastcol)+" "+strconv.Itoa(lastrow))
}
func givecolor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	color := "yellow"
	if turn {
		color = "red"
	} else {
		color = "yellow"
	}
	fmt.Fprint(w, color)
}

func DoTurn(w http.ResponseWriter, r *http.Request) {
	x, err := strconv.Atoi(r.FormValue("x"))
	if err == nil {
		placepiece(x)
		if horizontalcheck(x, turn) || verticalcheck(x, turn) || diagcheck(x, x, turn) {
			print("y'a qqun qui a gagfnéé")
		}
	}
	print("error ntm touche pas a mon code connard de tes morts")

}
func placepiece(x int) {
	for i := 0; i <= 6; i++ {

		if a[x][i] == 0 {
			if turn {
				lastcol = x
				lastrow = i
				a[x][i] = 1
				turn = !turn
				return
			} else {
				lastcol = x
				lastrow = i
				a[x][i] = 2
				turn = !turn
				return
			}
		}
	}
	//pas pu placer le pion :(
}

func horizontalcheck(y int, col bool) bool {
	acc := 0
	red := false
	for _, v := range a {
		if acc == 4 {
			if red == col {
				return true
			} else {
				return false
			}
		}
		if v[y] == 0 {
			acc = 0
		}
		if acc == 0 { //début d'acc
			switch v[y] {
			case 2:
				acc++
				red = true

			case 1:
				acc++
				red = false
			}
		} else if (red && v[y] == 1) || (!red && v[y] == 2) { //reset d'acc
			acc = 1
			red = !red
		} else if (!red && v[y] == 1) || (red && v[y] == 2) {
			acc++
		}
	}
	return false
}

func verticalcheck(x int, col bool) bool {
	acc := 0
	red := false
	for i := 0; i <= 5; i++ {
		if acc == 4 {
			if red == col {
				return true
			} else {
				return false
			}
		}
		if a[x][i] == 0 {
			acc = 0
		}
		if acc == 0 { //début d'acc
			switch a[x][i] {
			case 2:
				acc++
				red = true

			case 1:
				acc++
				red = false
			}
		} else if (red && a[x][i] == 1) || (!red && a[x][i] == 2) { //reset d'acc
			acc = 1
			red = !red
		} else if (!red && a[x][i] == 1) || (red && a[x][i] == 2) {
			acc++
		}
	}
	return false
}

func diagcheck(x, y int, col bool) bool {
	//chercher le coin  haut gauche de la diag
	//descendre en passcheck
	//cherche rle coin haut droit de la diag
	//descendre en passcheck
	return false
}

package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

var lastrow int
var lastcol int
var checkcursorx = 0
var checkcursory = 0
var istetris = true
var redpoints = 0
var yellowpoints = 0

type PageData struct {
	Message   string
	BoardHTML template.HTML
	TurnName  string //red or yellow
	TurnClass string //p1 or p2
}

type Color struct {
	Color string
}

type WinData struct {
	Winner string
}

var lastWinner string
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

	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/submit", submit)
	http.HandleFunc("/turn", DoTurn)
	http.HandleFunc("/color", givecolor)
	http.HandleFunc("/pos", getpos)
	http.HandleFunc("/reset", reset)
	http.HandleFunc("/win", winPage)
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
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	x, err := strconv.Atoi(r.FormValue("x"))
	if err == nil {
		placepiece(x)
		if istetris {
			for range 4 {
				tetrisendturn()
			}

			if redpoints >= 5 {
				lastWinner = "Rouge"
				http.Redirect(w, r, "/win", http.StatusSeeOther)
				turn = !turn
				return
			} else if yellowpoints >= 5 {
				lastWinner = "Jaune"
				http.Redirect(w, r, "/win", http.StatusSeeOther)
				turn = !turn
				return
			} else {
				won := horizontalcheck(lastrow, !turn) || verticalcheck(x, !turn) || diagcheck(x, lastrow, !turn) || diagcheck2(x, lastrow, !turn)
				if won {
					if a[lastcol][lastrow] == 1 {
						lastWinner = "Rouge"
					} else {
						lastWinner = "Jaune"
					}
					http.Redirect(w, r, "/win", http.StatusSeeOther)
					turn = !turn
					return
				}
			}
		}
	}
	turn = !turn
	print("error ntm touche pas a mon code connard de tes morts\n")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func placepiece(x int) {
	for i := 0; i < 6; i++ {

		if a[x][i] == 0 {
			lastcol = x
			lastrow = i
			if turn {
				a[x][i] = 1
			} else {
				a[x][i] = 2
			}

			printBoard()

			return
		}
	}
}

func horizontalcheck(y int, col bool) bool {
	acc := 0
	red := false
	checkcursory = y
	for i, v := range a {
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
		checkcursorx = i
	}
	return false
}

func verticalcheck(x int, col bool) bool {
	acc := 0
	red := false
	checkcursorx = x
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
		checkcursory = i
	}
	return false
}
func diagcheck2(x, y int, col bool) bool { // dans ce sens / de haut en bas
	check, check2 := x, y
	for check < 6 && check2 < 5 {
		check++
		check2++
	}
	checkcursorx = x
	checkcursory = y
	acc, red := 0, false
	for check >= 0 && check < 7 && check2 >= 0 && check2 < 6 {
		if acc == 4 {
			return red == col
		}
		v := a[check][check2]
		if v == 0 {
			acc = 0
		} else if acc == 0 {
			if v == 2 {
				acc++
				red = true
			} else if v == 1 {
				acc++
				red = false
			}
		} else if (red && v == 1) || (!red && v == 2) {
			acc = 1
			red = !red
		} else if (!red && v == 1) || (red && v == 2) {
			acc++
		}
		checkcursorx = check
		checkcursory = check2
		check--
		check2--
	}
	if acc == 4 {
		return red == col
	}

	return false
}

func diagcheck(x, y int, col bool) bool {
	check, check2 := x, y
	checkcursorx = x
	checkcursory = y
	for check > 0 && check2 < 5 {
		check--
		check2++
	}
	acc, red := 0, false
	for check >= 0 && check < 7 && check2 >= 0 && check2 < 6 {

		if acc == 4 {
			return red == col
		}
		v := a[check][check2]
		if v == 0 {
			acc = 0
		} else if acc == 0 {
			if v == 2 {
				acc++
				red = true
			} else if v == 1 {
				acc++
				red = false
			}
		} else if (red && v == 1) || (!red && v == 2) {
			acc = 1
			red = !red
		} else if (!red && v == 1) || (red && v == 2) {
			acc++
		}
		checkcursorx = check
		checkcursory = check2
		check++
		check2--
	}
	if acc == 4 {
		return red == col
	}
	return false
}

func renderBoard() template.HTML {
	var sb strings.Builder

	for row := 5; row >= 0; row-- {
		for col := 0; col < 7; col++ {
			v := a[col][row]
			cls := "cell"
			if v == 1 {
				cls += " p1"
			}
			if v == 2 {
				cls += " p2"
			}
			fmt.Fprintf(&sb, `<div class="%s" data-col="%d" data-row="%d"></div>`, cls, col, row)
		}
	}
	return template.HTML(sb.String())
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	turnName := "Joueur Jaune"
	turnClass := "p2"
	if turn { //red's turn
		turnName = "Joueur Rouge"
		turnClass = "p1"
	}

	data := PageData{
		Message:   "",
		BoardHTML: renderBoard(),
		TurnName:  turnName,
		TurnClass: turnClass,
	}
	_ = tmpl.Execute(w, data)
}

func reset(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		for x := 0; x < 7; x++ {
			for y := 0; y < 6; y++ {
				a[x][y] = 0
			}
		}
		turn = false
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func winPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/win.html"))
	_ = tmpl.Execute(w, WinData{Winner: lastWinner})
}

func printBoard() {
	for row := 5; row >= 0; row-- {
		for col := 0; col < 7; col++ {
			fmt.Printf("%d ", a[col][row])
		}
		fmt.Println()
	}
	fmt.Println()
}

func gravity_fix() { //appeler cette fonction 5 fois au minimum si on veux que tout se passe bien
	tamp := 0
	for x := 0; x <= 6; x++ {
		for i := 5; i > 0; i-- {
			//on descend de 1 le pion ducoup
			if a[x][i-1] == 0 {
				tamp = a[x][i]
				a[x][i] = 0
				a[x][i-1] = tamp
			}
		}
	}
}

func tetrisendturn() {
	if horizontalcheck(lastrow, !turn) {
		//enlever derriere le curseur
		for i := 0; i < 4; i++ {
			a[checkcursorx][checkcursory] = 0
			checkcursorx--
		}
		//gravitas appellée 5 fois pour eviter pire edge case
		for i := 0; i < 5; i++ {
			gravity_fix()
		}
		if turn {
			redpoints++
		} else {
			yellowpoints++
		}
		print("redpoints : " + strconv.Itoa(redpoints))
		print("yellowpoints : " + strconv.Itoa(yellowpoints))
		return
	}
	if verticalcheck(lastcol, !turn) {
		//enlever derriere le curseur
		for i := 0; i < 4; i++ {
			a[checkcursorx][checkcursory] = 0
			checkcursory--
		}
		//gravitas appellée 5 fois pour eviter pire edge case
		for i := 0; i < 5; i++ {
			gravity_fix()
		}
		if turn {
			redpoints++
		} else {
			yellowpoints++
		}
		print("redpoints : " + strconv.Itoa(redpoints))
		print("yellowpoints : " + strconv.Itoa(yellowpoints))
		return
	}
	if diagcheck(lastcol, lastrow, !turn) {
		for i := 0; i < 4; i++ {
			a[checkcursorx][checkcursory] = 0
			checkcursorx--
			checkcursory++
		}
		//gravitas appellée 5 fois pour eviter pire edge case
		for i := 0; i < 5; i++ {
			gravity_fix()
		}
		if turn {
			redpoints++
		} else {
			yellowpoints++
		}
		print("redpoints : " + strconv.Itoa(redpoints))
		print("yellowpoints : " + strconv.Itoa(yellowpoints))
		return
	}
	if diagcheck2(lastcol, lastrow, !turn) {
		for i := 0; i < 4; i++ {
			a[checkcursorx][checkcursory] = 0
			checkcursorx++
			checkcursory++
		}

		//gravitas appellée 5 fois pour eviter pire edge case
		for i := 0; i < 5; i++ {
			gravity_fix()
		}
		if turn {
			redpoints++
		} else {
			yellowpoints++
		}
		print("redpoints : " + strconv.Itoa(redpoints))
		print("yellowpoints : " + strconv.Itoa(yellowpoints))
		return
	}
}

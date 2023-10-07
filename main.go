package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
    "regexp"
    "github.com/gorilla/mux"
	"html/template"

	"github.com/gorcon/rcon"
)

var MC_PASS = os.Getenv("MC_PASS")
var MC_ADDRPORT = os.Getenv("MC_ADDR")

var rconClient *rcon.Conn

func main() {
	initRcon()
	fmt.Println("Connected to server")
    //shout(rconClient, "DNAK")
    r := mux.NewRouter()
    r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
    r.HandleFunc("/shout/{message}", serveShout)
    r.HandleFunc("/shout", serveShout)
    r.HandleFunc("/", serveHome)
    http.Handle("/", r)
    http.ListenAndServe(":7069", nil)

}

func serveHome(w http.ResponseWriter, r *http.Request) {
    conn := rconClient
	players := listPlayers(conn)
    //players = "There are 2 of a max of 20 players online: User1, User2"
    re := regexp.MustCompile(`There are (\d+) of a max of \d+ players online: (.+)`)
    matches := re.FindStringSubmatch(players)
    userList := matches[2]
    users := []string{}
    usernames := regexp.MustCompile(`,\s*`).Split(userList, -1)
    for _, username := range(usernames) {
        users = append(users, username)
    }
    tmpl, err := template.ParseFiles("templates/base.html", "templates/list.html")
    if err != nil {
        log.Fatal(err)
    }

    err = tmpl.ExecuteTemplate(w, "base.html", users)
    return
}

func listPlayers(conn *rcon.Conn) string {
	response, err := conn.Execute("/list")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response)
	return response
}

func initRcon() {
	rconClient, _ = rcon.Dial(MC_ADDRPORT, MC_PASS)
}

func serveShout(w http.ResponseWriter, r *http.Request){
    fmt.Println(r.Method)
    if r.Method == "GET"{
    tmpl, err := template.ParseFiles("templates/base.html", "templates/shout.html")
    if err != nil {
        log.Fatal(err)
    }

    err = tmpl.ExecuteTemplate(w, "base.html", "")
    return
    }
    if r.Method == "POST" {
        value := r.PostFormValue("shout")
        fmt.Println(value)
        conn := rconClient
        _, err := conn.Execute("/say " + value)
        if err != nil {
            log.Fatal(err)
        }
        http.Redirect(w, r, "/shout", 301)
    }
}

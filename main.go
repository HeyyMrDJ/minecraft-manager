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
    r.HandleFunc("/unban/{message}", serveUnban)
    r.HandleFunc("/unban", serveUnban)
    r.HandleFunc("/ban/{message}", serveBan)
    r.HandleFunc("/ban", serveBan)
    r.HandleFunc("/shout/{message}", serveShout)
    r.HandleFunc("/shout", serveShout)
    r.HandleFunc("/", serveHome)
    http.Handle("/", r)
    http.ListenAndServe(":7069", nil)

}

func serveHome(w http.ResponseWriter, r *http.Request) {
    users := getPlayers()
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

func serveBan(w http.ResponseWriter, r *http.Request){
    fmt.Println(r.Method)
    if r.Method == "GET"{
    tmpl, err := template.ParseFiles("templates/base.html", "templates/ban.html")
    if err != nil {
        log.Fatal(err)
    }
    players := getPlayers()

    err = tmpl.ExecuteTemplate(w, "base.html", players)
    return
    }
    if r.Method == "POST" {
        value := r.PostFormValue("ban")
        fmt.Println(value)
        conn := rconClient
        _, err := conn.Execute("/ban " + value)
        if err != nil {
            log.Fatal(err)
        }
        http.Redirect(w, r, "/ban", 301)
    }
}

func serveUnban(w http.ResponseWriter, r *http.Request){
    fmt.Println(r.Method)
    if r.Method == "GET"{
    tmpl, err := template.ParseFiles("templates/base.html", "templates/unban.html")
    if err != nil {
        log.Fatal(err)
    }
    players := getBanned()

    err = tmpl.ExecuteTemplate(w, "base.html", players)
    return
    }
    if r.Method == "POST" {
        value := r.PostFormValue("unban")
        fmt.Println(value)
        conn := rconClient
        _, err := conn.Execute("/pardon " + value)
        if err != nil {
            log.Fatal(err)
        }
        http.Redirect(w, r, "/unban", 301)
    }
}

func getPlayers() []string{
    conn := rconClient
	players := listPlayers(conn)
    log.Print("Players Length:",len(players))
    users := []string{}
    if len(players) >= 44{
        re := regexp.MustCompile(`There are (\d+) of a max of \d+ players online: (.+)`)
        matches := re.FindStringSubmatch(players)
        userList := matches[2]
        usernames := regexp.MustCompile(`,\s*`).Split(userList, -1)
        for _, username := range(usernames) {
            users = append(users, username)
        }
    } else {
        //users = append(users, "")
        fmt.Println(1)
    }
    return users
}

func getBanned() []string{
    conn := rconClient
    banned, err := conn.Execute("/banlist")
    if err != nil {
        log.Fatal(err)
    }
	// Define a regular expression pattern to match usernames.
	re := regexp.MustCompile(`[:.](\w+) was`)

	// Find all matches of the pattern in the text.
	matches := re.FindAllStringSubmatch(banned, -1)

	// Extract and print the usernames.
	var usernames []string
	for _, match := range matches {
		if len(match) > 1 {
			usernames = append(usernames, match[1])
		}
	}

	// Print the matched usernames.
	if len(usernames) == 0 {
		fmt.Println("No users found.")
	} else {
		fmt.Println("Found users:")
		for _, username := range usernames {
			fmt.Println(username)
		}
	}
    return usernames
}

// Helper function to check if a string contains only digits.
func containsOnlyDigits(s string) bool {
	for _, char := range s {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

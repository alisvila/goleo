package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type user struct {
	Username string `json:"name"`
	Password string
	first    string
	last     string
}

var dbUser = map[string]user{}
var dbSeason = make(map[string]string)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", index)
	r.HandleFunc("/sign-up", signup)
	http.Handle("/", r)
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	u, err := GetUser(w, r)
	if err != nil {
		http.Error(w, "user dosent exist", http.StatusForbidden)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(u)
}

func signup(w http.ResponseWriter, r *http.Request) {
	// if alreadeyLoggedIn(r) {
	// 	http.Redirect(w, r, "/", http.StatusSeeOther)
	// 	return
	// }

	if r.Method == http.MethodPost {
		un := r.FormValue("username")
		p := r.FormValue("password")
		f := r.FormValue("firstname")
		l := r.FormValue("lastname")

		if _, ok := dbUser[un]; ok {
			http.Error(w, " username already exist", http.StatusForbidden)
		}

		sID := uuid.New()
		c := &http.Cookie{
			Name:  "season",
			Value: sID.String(),
		}

		http.SetCookie(w, c)
		dbSeason[c.Value] = un

		u := user{un, p, f, l}
		dbUser[un] = u

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) (user, error) {
	c, err := r.Cookie("season")
	if err != nil {
		fmt.Println(c)
		http.Redirect(w, r, "/signup", http.StatusSeeOther)
		return user{Username: "ali", Password: "123", first: "a", last: "ff"}, errors.New("user name dosent exist")
	} else {
		var u user
		if un, ok := dbSeason[c.Value]; ok {
			u = dbUser[un]
		}

		return u, nil
	}

}

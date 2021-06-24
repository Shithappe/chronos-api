package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"encoding/json"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type User struct {
	Id    string `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

var Users []User

func home_page(w http.ResponseWriter, r *http.Request) {
	Users = []User{}
	res, _ := db.Query("select * from `Users`")
	for res.Next() {
		var temp User
		res.Scan(&temp.Id, &temp.Login, &temp.Email, &temp.Name)
		Users = append(Users, temp)
	}
	fmt.Println(Users)
	json.NewEncoder(w).Encode(Users)
}

func register(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var user User
	json.Unmarshal(reqBody, &user)
	fmt.Println(user.Login, user.Email, user.Name)
	db.Exec("insert into chronos.Users (login, email, name) values (?, ?, ?)",
		user.Login, user.Email, user.Name)
	json.NewEncoder(w).Encode(user)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", home_page)
	myRouter.HandleFunc("/register", register).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", myRouter))
}

var db = InitDB()

func InitDB() *sql.DB {
	var err error
	db, err := sql.Open("mysql", "MAX:MAX@tcp(127.0.0.1:3306)/chronos")
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
	return db
}

func main() {

	// res, err := db.Query("select * from `Users`")
	// for res.Next() {
	// 	var temp User
	// 	res.Scan(&temp.Id, &temp.Login, &temp.Email, &temp.Name)
	// 	Users = append(Users, temp)
	// }

	handleRequests()
}

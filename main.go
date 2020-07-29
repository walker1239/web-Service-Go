package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"database/sql"
	"fmt"
	"encoding/json"
	"io/ioutil"
	_ "github.com/mattn/go-sqlite3"
)

type Note struct{
	Id int				`json:"id"`
	Title string		`json:"title"`	
	Description string	`json:"description"`
}


func GetConnection() *sql.DB {
	database, errordb := sql.Open("sqlite3", "./database.db")
	if errordb != nil {
		//return []Note{}, err
		log.Println(errordb.Error())
	}
	
	statement, errorstmt := database.Prepare("CREATE TABLE IF NOT EXISTS notes (id INTEGER PRIMARY KEY AUTOINCREMENT, title VARCHAR(100) NULL,description VARCHAR(500) NULL);")
	if errorstmt != nil {
		//return []Note{}, err
		log.Println(errorstmt.Error())
	}
	statement.Exec()
	return database
}

func holaName(w http.ResponseWriter, r *http.Request){
	name := mux.Vars(r)["name"]
	fmt.Fprint(w, "hola "+name)
}

func getNotes(w http.ResponseWriter, r *http.Request){
	db := GetConnection()
	var notes []Note
	rows, _ := db.Query("SELECT id, title, description FROM notes")
    var note Note
    for rows.Next() {
		rows.Scan(&note.Id, &note.Title, &note.Description)
		notes=append(notes,note)
    }
	json.NewEncoder(w).Encode(notes)
}

func getNote(w http.ResponseWriter, r *http.Request){
	db := GetConnection()

	id_String := mux.Vars(r)["id"]
	rows, _ := db.Query("SELECT id, title, description FROM notes WHERE id = "+id_String)
    var note Note
    for rows.Next() {
		rows.Scan(&note.Id, &note.Title, &note.Description)
    }
	json.NewEncoder(w).Encode(note)	
}

func createNote(w http.ResponseWriter, r *http.Request){
	db := GetConnection()
	var temp Note
	body,_:=ioutil.ReadAll(r.Body)
	json.Unmarshal(body,&temp)
	statement, err := db.Prepare("INSERT INTO notes (title, description) VALUES (?, ?)")
	statement.Exec(temp.Title, temp.Description)
	if err != nil {
		message1 := map[string]string{
			"code":"error",
			"message":err.Error(),
		}
		json.NewEncoder(w).Encode(message1)
		return
	}
	message1 := map[string]string{
		"code":"correcto",
		"message":"Nota registrada correctamente",
	}
	json.NewEncoder(w).Encode(message1)
}

func updateNote(w http.ResponseWriter, r *http.Request){
	db := GetConnection()
	var temp Note
	id_String := mux.Vars(r)["id"]
	body,_:=ioutil.ReadAll(r.Body)
	json.Unmarshal(body,&temp)
	statement, err := db.Prepare("UPDATE notes SET title = ?, description= ? WHERE id = ?")
	statement.Exec(temp.Title, temp.Description, id_String)
	if err != nil {
		message1 := map[string]string{
			"code":"error",
			"message":err.Error(),
		}
		json.NewEncoder(w).Encode(message1)
		return
	}
	message1 := map[string]string{
		"code":"correcto",
		"message":"Nota actualizada correctamente",
	}
	json.NewEncoder(w).Encode(message1)
}

func deleteNote(w http.ResponseWriter, r *http.Request){
	db := GetConnection()

	id_String := mux.Vars(r)["id"]

	statement, err := db.Prepare("DELETE FROM notes WHERE id = ?")
	statement.Exec(id_String)
	if err != nil {
		message1 := map[string]string{
			"code":"error",
			"message":err.Error(),
		}
		json.NewEncoder(w).Encode(message1)
		return
	}
	message1 := map[string]string{
		"code":"correcto",
		"message":"Nota eliminada correctamente",
	}
	json.NewEncoder(w).Encode(message1)
}

func main(){

	r := mux.NewRouter()
	r.HandleFunc("/note", createNote).Methods("POST")
	r.HandleFunc("/notes", getNotes).Methods("GET")
	r.HandleFunc("/notes/{id}", getNote).Methods("GET")
	r.HandleFunc("/notes/{id}", updateNote).Methods("PUT")
	r.HandleFunc("/notes/{id}", deleteNote).Methods("DELETE")

	log.Print("Corriendo en el puerto 8085")
	err := http.ListenAndServe(":8085",r)
	if err != nil{
		log.Fatal("error: ",err)
	}
}
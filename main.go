package main

import (
	"log"
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"strconv"
)

type Note struct{
	Id int				`json:"id"`
	Title string		`json:"title"`	
	Description string	`json:"description"`
	EmailUser string 		
}

type User struct{
	Email string		`json:"email"`		
	Password string		`json:"password"`	
	Token string			
}

var notes []Note
var users []User

var jwtKey = []byte("13Akfq195g")

func holaName(w http.ResponseWriter, r *http.Request){
	name := mux.Vars(r)["name"]
	fmt.Fprint(w, "hola "+name)
}

func getNotes(w http.ResponseWriter, r *http.Request){
	json.NewEncoder(w).Encode(notes)
}

func getUsers(w http.ResponseWriter, r *http.Request){
	json.NewEncoder(w).Encode(users)
}

func getNote(w http.ResponseWriter, r *http.Request){
	id_String := mux.Vars(r)["id"]
	id,_ := strconv.Atoi(id_String)
	for i := range notes{
		if notes[i].Id == id{
			json.NewEncoder(w).Encode(notes[i])
		}
	}
	
}

func getNotesByUser(w http.ResponseWriter, r *http.Request){
	var notesUser []Note

	tokenHeader := r.Header.Get("Authorization")
	for i := range users{
		if users[i].Token == tokenHeader{
			for j := range notes{
				if notes[j].EmailUser == users[i].Email {
					notesUser=append(notesUser,notes[j])
				}
			}
			json.NewEncoder(w).Encode(notesUser)
			return
		}
	}
	message1 := map[string]string{
		"code":"error",
		"message":"Token no reconocido, vuelva a iniciar sesión",
	}
	json.NewEncoder(w).Encode(message1)
}

func createNote(w http.ResponseWriter, r *http.Request){
	var temp Note
	tokenHeader := r.Header.Get("Authorization")
	body,_:=ioutil.ReadAll(r.Body)
	json.Unmarshal(body,&temp)
	for i := range users{
		if users[i].Token == tokenHeader{
			temp.EmailUser=users[i].Email
			temp.Id=len(notes)+1
			notes = append(notes,temp)
			message1 := map[string]string{
				"code":"correcto",
				"message":"Nota registrada correctamente",
			}
			json.NewEncoder(w).Encode(message1)
			return
		}
	}
	message1 := map[string]string{
		"code":"error",
		"message":"Token no reconocido, vuelva a iniciar sesión",
	}
	json.NewEncoder(w).Encode(message1)
}

func updateNote(w http.ResponseWriter, r *http.Request){
	var updateNote Note

	id_String := mux.Vars(r)["id"]
	id,_ := strconv.Atoi(id_String)

	tokenHeader := r.Header.Get("Authorization")
	body,_:=ioutil.ReadAll(r.Body)
	json.Unmarshal(body,&updateNote)
	for i := range users{
		if users[i].Token == tokenHeader{
			for j := range notes{
				if notes[j].Id == id && notes[j].EmailUser == users[i].Email {
					notes[j].Title=updateNote.Title
					notes[j].Description=updateNote.Description
					message1 := map[string]string{
						"code":"correcto",
						"message":"Nota actualizada correctamente",
					}
					json.NewEncoder(w).Encode(message1)
					return
				}
			}
			message1 := map[string]string{
				"code":"error",
				"message":"Nota no encontrada",
			}
			json.NewEncoder(w).Encode(message1)
			return
		}
	}
	message1 := map[string]string{
		"code":"error",
		"message":"Token no reconocido, vuelva a iniciar sesión",
	}
	json.NewEncoder(w).Encode(message1)
}

func remove(notesdel []Note, i int) []Note {
    notesdel[len(notesdel)-1], notesdel[i] = notesdel[i], notesdel[len(notesdel)-1]
    return notesdel[:len(notesdel)-1]
}

func deleteNote(w http.ResponseWriter, r *http.Request){
	id_String := mux.Vars(r)["id"]
	id,_ := strconv.Atoi(id_String)

	tokenHeader := r.Header.Get("Authorization")

	for i := range users{
		if users[i].Token == tokenHeader{
			for j := range notes{
				if notes[j].Id == id && notes[j].EmailUser == users[i].Email {
					notes=remove(notes,j);
					message1 := map[string]string{
						"code":"correcto",
						"message":"Nota eliminada correctamente",
					}
					json.NewEncoder(w).Encode(message1)
				}
			}
			message1 := map[string]string{
				"code":"error",
				"message":"Nota no encontrada",
			}
			json.NewEncoder(w).Encode(message1)
			return
		}
	}
	message1 := map[string]string{
		"code":"error",
		"message":"Token no reconocido, vuelva a iniciar sesión",
	}
	json.NewEncoder(w).Encode(message1)
}

func createUser(w http.ResponseWriter, r *http.Request){
	var account User
	body,_:=ioutil.ReadAll(r.Body)
	err:=json.Unmarshal(body,&account)
	if err != nil {
		log.Println(err.Error())
	}
	for i := range users{	
		if users[i].Email == account.Email{
			message1 := map[string]string{
				"code":"error",
				"message":"Correo ya existente",
			}
			json.NewEncoder(w).Encode(message1)
			return
		}
	}
	account.Token="none"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)
	users = append(users,account)
	message1 := map[string]string{
		"code":"correcto",
		"message":"Usuario registrado con exito",
	}
	json.NewEncoder(w).Encode(message1)
}

func loginUser(w http.ResponseWriter, r *http.Request){
	var account User

	body,_:=ioutil.ReadAll(r.Body)
	json.Unmarshal(body,&account)
	//notes = append(notes,temp)
	for i := range users{
		
		if users[i].Email == account.Email{
			err := bcrypt.CompareHashAndPassword([]byte(users[i].Password), []byte(account.Password))
			if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { 
				message := map[string]string{
					"code":"error",
					"message":"Contraseña o correo invalidos",
				}
				json.NewEncoder(w).Encode(message)
				return
			}
			userClaims := jwt.MapClaims{}
			userClaims["email"] = users[i].Email
			userClaims["password"] = users[i].Password
			at := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
  			token, err := at.SignedString(jwtKey)
			users[i].Token = token
			message := map[string]string{
				"code":"correcto",
				"message":"Login correcto",
				"Api": users[i].Token,
			}
			json.NewEncoder(w).Encode(message)
			return
		}
	}
	message1 := map[string]string{
		"code":"error",
		"message":"Ocurrio un error",
	}
	json.NewEncoder(w).Encode(message1)
}

func main(){

	notes=append(notes,Note{Id:1,Title:"titulo1",Description:"description1"})
	notes=append(notes,Note{Id:2,Title:"titulo2",Description:"description2"})

	r := mux.NewRouter()
	r.HandleFunc("/note", createNote).Methods("POST")
	r.HandleFunc("/notes", getNotes).Methods("GET")
	r.HandleFunc("/user/notes", getNotesByUser).Methods("GET")
	r.HandleFunc("/notes/{id}", getNote).Methods("GET")
	r.HandleFunc("/notes/{id}", updateNote).Methods("PUT")
	r.HandleFunc("/notes/{id}", deleteNote).Methods("DELETE")
	r.HandleFunc("/user", createUser).Methods("POST")
	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/login", loginUser).Methods("POST")

    //r.HandleFunc("/{name}", holaName).Methods("GET")

	log.Print("Corriendo en el puerto 8085")
	err := http.ListenAndServe(":8085",r)
	if err != nil{
		log.Fatal("error: ",err)
	}
}
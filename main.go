package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"main/database"
)

var uDB *database.UserDataBase

func main() {
	var err error
	uDB, err = database.NewUserDataBase("users.db")
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer uDB.Close()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/addUser", addUserHandler)
	http.HandleFunc("/authAdmin", authAdminHandler)
	http.HandleFunc("/editUser", editUserHandler)
	http.HandleFunc("/listUsers", listUsersHandler)
	http.Handle("/src/", http.StripPrefix("/src/", http.FileServer(http.Dir("web/src/"))))

	port := ":8080"
	fmt.Printf("Server listening on port %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/public/index.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := uDB.AddNewUser(user.Username, user.Password); err != nil {
		http.Error(w, "Error adding user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User added successfully")
}

func authAdminHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if !uDB.AuthenticateAdmin(user.Username, user.Password) {
		http.Error(w, "Invalid admin credentials", http.StatusUnauthorized)
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Admin logged in successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func editUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var user struct {
		UserID      string `json:"user_id"`
		NewUsername string `json:"new_username"`
		NewPassword string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := uDB.EditUser(user.UserID, user.NewUsername, user.NewPassword); err != nil {
		http.Error(w, "Error editing user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User edited successfully")
}

func listUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	users, err := uDB.ListUsers()
	if err != nil {
		http.Error(w, "Error listing users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Error encoding JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
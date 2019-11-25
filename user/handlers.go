package user

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/zhughes3/website/db"
	httpHelper "github.com/zhughes3/website/http"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	var users []User
	db.DB.Find(&users)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func ShowHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	params := mux.Vars(r)
	db.DB.First(&user, params["userId"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	type createuser struct {
		Email string `json:"email"`
		Name string `json:"name"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	var userToCreate createuser
	err := decoder.Decode(&userToCreate)
	if err != nil {
		httpHelper.NewErrorResponse(w, http.StatusInternalServerError, "Error: " + err.Error())
		return
	}
	user := User{
		Email: userToCreate.Email,
		Name: userToCreate.Name,
		Hash: getHashFromPassword(userToCreate.Password),
	}
	error := db.DB.Create(&user).Error
	if error != nil {
		httpHelper.NewErrorResponse(w, http.StatusUnauthorized, "Error: " + err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&user)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	type logininput struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	var userToLogin logininput
	err := decoder.Decode(&userToLogin)
	if err != nil {
		httpHelper.NewErrorResponse(w, http.StatusInternalServerError, "Error: " + err.Error())
		return
	}
	var user User
	db.DB.Where("email = ?", userToLogin.Email).Find(&user)
	w.Header().Set("Content-Type", "application/json")
	if user.checkPassword(userToLogin.Password) {
		token, err := user.generateJWT()
		if err != nil {
			httpHelper.NewErrorResponse(w, http.StatusUnauthorized, "Error: " + err.Error())
			return
		}
		json.NewEncoder(w).Encode(&token)
	} else {
		httpHelper.NewErrorResponse(w, http.StatusUnauthorized, "Password incorrect")
		return
	}
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var user User
	var users []User

	db.DB.First(&user, params["userId"])
	db.DB.Delete(&user)

	db.DB.Find(&users)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	type updateuser struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	var userToUpdate updateuser
	err := decoder.Decode(&userToUpdate)
	if err != nil {
		httpHelper.NewErrorResponse(w, http.StatusInternalServerError, "Error: " + err.Error())
		return
	}
	params := mux.Vars(r)
	var user User
	reqUserId := r.Header.Get("userId")

	w.Header().Set("Content-Type", "application/json")
	if params["userId"] != reqUserId {
		httpHelper.NewErrorResponse(w, http.StatusUnauthorized, "Not allowed to edit other users")
		return
	}
	db.DB.First(&user, params["userId"])
	db.DB.Model(&user).Update("name", userToUpdate.Name)
	json.NewEncoder(w).Encode(&user)
}

func getHashFromPassword(password string) string{
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 4)
	return string(hash)
}
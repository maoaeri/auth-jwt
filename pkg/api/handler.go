package api

import (
	"errors"
	"fmt"
	"html/template"
	"myapp/pkg/helper"
	"myapp/pkg/model"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type data struct {
	Response string
}

func Register(w http.ResponseWriter, r *http.Request) {
	connection := GetDB()
	defer CloseDB(connection)

	responseMap := map[string]interface{}{}

	if r.FormValue("signup") != "" && r.FormValue("email") != "" {
		user := model.User{
			Email:    r.FormValue("email"),
			Name:     r.FormValue("name"),
			Password: r.FormValue("password"),
		}

		var dbuser model.User
		err := connection.Where("email = ?", user.Email).First(&dbuser).Error
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			res := "Email already in use"
			responseMap["Response"] = res
			responseMap["success"] = false
		} else {
			user.Password = helper.GenerateHash(user.Password)
			result := connection.Create(&user)
			if result.Error != nil {
				fmt.Print("Error in creating user.")
				fmt.Print(result.Error)
			}
			responseMap["Response"] = "Successfully registered"
			responseMap["success"] = true
			if err := connection.Model(&user).Update("role", "user"); err != nil {
				fmt.Print("Error in setting role for user.")
				fmt.Print(err)
			}
		}
	}
	tmpl := template.Must(template.ParseFiles("../pkg/template/register.html"))
	if err := tmpl.Execute(w, responseMap); err != nil {
		fmt.Print(err)
	}

}

func LogIn(w http.ResponseWriter, r *http.Request) {

	connection := GetDB()
	defer CloseDB(connection)

	responseMaps := map[string]interface{}{}

	email, _ := helper.ExtractEmail(r)
	if email != "" {
		responseMaps["Email"] = email
	}

	if r.FormValue("login") != "" && r.FormValue("email") != "" {
		auth := model.Authentication{
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		var authuser model.User
		connection.Where("email = ?", auth.Email).First(&authuser)
		if authuser.Email == "" {
			responseMaps["Response"] = "Email or password is incorrect"
		} else {
			check := helper.CheckPasswordHash(auth.Password, authuser.Password)

			if !check {
				responseMaps["Response"] = "Email or password is incorrect"
			} else {
				validToken := helper.GenerateTokenPairString(authuser.Email, authuser.Role)

				//save refresh token to database
				if err := connection.Model(&authuser).Update("refresh_token", validToken["refreshtoken"]); err != nil {
					fmt.Print("Error in saving refresh token.")
					fmt.Print(err)
				}

				expirationTime := time.Now().Add(1 * time.Hour) // cookie expired after 1 hour

				cookie := &http.Cookie{
					Name:     "token",
					Value:    validToken["token"],
					Expires:  expirationTime,
					HttpOnly: true,
					Secure:   true,
				}
				http.SetCookie(w, cookie) //set cookies
				http.Redirect(w, r, "/index", http.StatusFound)
			}

		}

	}
	tmpl := template.Must(template.ParseFiles("../pkg/template/login.html"))
	tmpl.Execute(w, responseMaps)

}

func LogOut(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:   "token",
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

func Index(w http.ResponseWriter, r *http.Request) {
	connection := GetDB()
	defer CloseDB(connection)

	role := r.Header.Get("Role")
	if role == "admin" {
		tmpl := template.Must(template.ParseFiles("../pkg/template/AdminIndex.html"))

		users := GetAllUsers()

		tmpl.Execute(w, users)
	} else if role == "user" {
		tmpl := template.Must(template.ParseFiles("../pkg/template/UserIndex.html"))
		userEmail, err := helper.ExtractEmail(r)
		if err != nil {
			helper.RespondError(w, http.StatusInternalServerError, "Error in getting user email")
		}
		user := GetUser(userEmail)
		tmpl.Execute(w, user)
	}
}

//func RefreshTokenHandler()

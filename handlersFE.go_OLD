package main

import (
	"log"
	"net/http"
	"text/template"
	"time"
)

// Serve the homepage (this serves the HTML page)
func serveHome(w http.ResponseWriter, r *http.Request) {

	// Server frontend-directory (vuejs project)
	//fs := http.FileServer(http.Dir("./frontend/dist/"))
	//http.StripPrefix("/", fs).ServeHTTP(w, r)

	// Server frontend-directory (golang template project)
	s := &Server{}
	s.serveHTML(w, r)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("bearer_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// Cookie not found, handle the error accordingly
			//http.Error(w, "Unauthorized", http.StatusUnauthorized)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		// Some other error occurred
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Access the token stored in the cookie
	bearerToken := cookie.Value

	user, errUser := validateBearerToken(bearerToken)
	if errUser != nil {
		log.Printf("error validating baerer-token - redirecting to login: %s", errUser)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	log.Printf("User %s validated using baerer-token.", user)
	log.Printf("cookie: %s", cookie.Value)

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	user := User{}
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		user, err = dbValidateUser(username, password)
		if err != nil {
			http.Error(w, "Error - Failed to login user: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// ToDo - funktioniert noch nicht wenn falscher User/passwort
		if user.Id == nil {
			http.Error(w, "Error - invalid login: ", http.StatusInternalServerError)
		}

		log.Printf("User logged in: %s", user.Username)
		// Create JWT token
		token, err := GenerateJWT(user.Username.(string))
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		// update token in DB

		errUpdToken := dbUpdateToken(user.Username.(string), token)
		if errUpdToken != nil {
			http.Error(w, "Failed to update token in db", http.StatusInternalServerError)
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "bearer_token",
			Value:    token,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			Path:     "/",
		})

		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	tmpl.Execute(w, nil)

}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	// ToDo : Prüfung auf gülltigen cookie-token

	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Example data to pass to the HTML template
	data := struct {
		Title string
		Body  string
	}{
		Title: "Welcome to My Web Server",
		Body:  "This is a simple HTML page served by Go.",
	}

	// Render the template with data
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
	}

}

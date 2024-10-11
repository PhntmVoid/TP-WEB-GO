package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

type Etudiant struct {
	Nom    string
	Prenom string
	Age    int
	Sexe   string
}

type Promo struct {
	Nom         string
	Filiere     string
	Niveau      string
	NbEtudiants int
	LsEtudiants []Etudiant
}

type Change struct {
	Pair    bool
	Counter int
}

type User struct {
	Nom             string
	Prenom          string
	DateDeNaissance string
	Sexe            string
}

var (
	viewCount   int
	currentUser *User
)

func main() {
	temp, err := template.ParseGlob("./templates/*.html")
	if err != nil {
		fmt.Printf("ERREUR => %s\n", err.Error())
		os.Exit(2)
	}

	promo := Promo{
		Nom:         "B1 Informatique",
		Filiere:     "Informatique",
		Niveau:      "Bachelor 1",
		NbEtudiants: 3,
		LsEtudiants: []Etudiant{
			{Nom: "Dupont", Prenom: "Alice", Age: 20, Sexe: "F"},
			{Nom: "Marley", Prenom: "Bob", Age: 21, Sexe: "M"},
			{Nom: "Durand", Prenom: "Agathe", Age: 19, Sexe: "F"},
		},
	}

	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/promo", func(w http.ResponseWriter, r *http.Request) {
		promoHandler(w, r, temp, promo)
	})

	http.HandleFunc("/change", func(w http.ResponseWriter, r *http.Request) {
		changeHandler(w, r, temp)
	})

	http.HandleFunc("/user/form", func(w http.ResponseWriter, r *http.Request) {
		userFormHandler(w, r, temp)
	})

	http.HandleFunc("/user/treatment", func(w http.ResponseWriter, r *http.Request) {
		userTreatmentHandler(w, r, temp)
	})

	http.HandleFunc("/user/display", func(w http.ResponseWriter, r *http.Request) {
		userDisplayHandler(w, r, temp)
	})

	http.HandleFunc("/user/error", func(w http.ResponseWriter, r *http.Request) {
		userErrorHandler(w, r, temp)
	})

	fmt.Println("Serveur démarré sur http://localhost:8000")
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Printf("Erreur du serveur: %s\n", err.Error())
		os.Exit(1)
	}
}

func promoHandler(w http.ResponseWriter, r *http.Request, temp *template.Template, promo Promo) {
	err := temp.ExecuteTemplate(w, "promo.html", promo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func changeHandler(w http.ResponseWriter, r *http.Request, temp *template.Template) {
	viewCount++
	changeData := Change{
		Pair:    viewCount%2 == 0,
		Counter: viewCount,
	}

	err := temp.ExecuteTemplate(w, "change.html", changeData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func userFormHandler(w http.ResponseWriter, r *http.Request, temp *template.Template) {
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/user/form", http.StatusSeeOther)
		return
	}

	err := temp.ExecuteTemplate(w, "form.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func userTreatmentHandler(w http.ResponseWriter, r *http.Request, temp *template.Template) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/user/form", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Redirect(w, r, "/user/error", http.StatusSeeOther)
		return
	}

	user := User{
		Nom:             strings.TrimSpace(r.FormValue("nom")),
		Prenom:          strings.TrimSpace(r.FormValue("prenom")),
		DateDeNaissance: strings.TrimSpace(r.FormValue("date_naissance")),
		Sexe:            strings.TrimSpace(r.FormValue("sexe")),
	}

	if valid, _ := validateUser(user); !valid {
		http.Redirect(w, r, "/user/error", http.StatusSeeOther)
		return
	}

	currentUser = &user

	http.Redirect(w, r, "/user/display", http.StatusSeeOther)
}

func userDisplayHandler(w http.ResponseWriter, r *http.Request, temp *template.Template) {
	if currentUser == nil {
		err := temp.ExecuteTemplate(w, "display.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err := temp.ExecuteTemplate(w, "display.html", currentUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func userErrorHandler(w http.ResponseWriter, r *http.Request, temp *template.Template) {
	err := temp.ExecuteTemplate(w, "error.html", "Données invalides. Veuillez réessayer.")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func validateUser(user User) (bool, string) {
	if !isValidName(user.Nom) {
		return false, "Le nom doit être composé de lettres et avoir une taille entre 1 et 32 caractères."
	}

	if !isValidName(user.Prenom) {
		return false, "Le prénom doit être composé de lettres et avoir une taille entre 1 et 32 caractères."
	}

	if user.Sexe != "masculin" && user.Sexe != "féminin" && user.Sexe != "autre" {
		return false, "Le sexe doit être 'masculin', 'féminin' ou 'autre'."
	}

	if user.DateDeNaissance == "" {
		return false, "La date de naissance est requise."
	}

	return true, ""
}

func isValidName(s string) bool {
	if len(s) < 1 || len(s) > 32 {
		return false
	}
	for _, char := range s {
		if !isLetter(char) {
			return false
		}
	}
	return true
}

func isLetter(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= 'À' && c <= 'ÿ')
}

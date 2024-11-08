package controllers

import (
	"net/http"

//"github.com/gin-gonic/gin"
	"github.com/IamMaheshGurung/pagination/initializers"
	"github.com/IamMaheshGurung/pagination/models"
    "html/template"
)

func PeopleIndexGET(w http.ResponseWriter, r * http.Request) {
	// Get the people
	var people []models.Person
	initializers.DB.Find(&people)


    tmpl, err := template.New("index").ParseFiles("index.tmpl")
    if err != nil {
        http.Error(w, "Error parsing the template", http.StatusInternalServerError)
        return
    }

	// Render the page
    err = tmpl.Execute(w, map[string]interface{}{
        "people": people,
    })
    if err != nil {
        http.Error(w, "Error in rendering template", http.StatusInternalServerError)
    }
}



package controllers

import (
    "net/http"
    "log"
    "github.com/IamMaheshGurung/pagination/initializers"
    "github.com/IamMaheshGurung/pagination/models"
    "html/template"
    "os"
  

)
func checkFileExists(filePath string) {
    _, err := os.Stat(filePath)
    if err != nil {
        if os.IsNotExist(err) {
            log.Printf("File does not exist: %s", filePath)
        } else {
            log.Printf("Error accessing file: %s", filePath)
        }
    } else {
        log.Printf("File exists: %s", filePath)
    }
}

func PeopleIndexGET(w http.ResponseWriter, r *http.Request) {
       checkFileExists("templates/people/index.tmpl")
    
    // Get the people from the database
    var people []models.Person
    if err := initializers.DB.Find(&people).Error; err != nil {
        log.Printf("Error fetching people from the database: %v", err)
        http.Error(w, "Error fetching people data", http.StatusInternalServerError)
        return
    }

    // Parse all templates (top, index, bottom)
    tmpl := template.Must(template.New("index").ParseFiles(
               "templates/people/index.tmpl",
            ))
    /*if tmpl.Error != nil {
        log.Printf("Error parsing templates: %v", err)
        http.Error(w, "Error parsing the template", http.StatusInternalServerError)
        return
    }*/

    // Log the parsed templates
    log.Println("Successfully loaded templates")

    // Render the page with the people data
    err := tmpl.Execute(w, map[string]interface{}{
        "people": people,
    })
    if err != nil {
        log.Printf("Error rendering template: %v", err)
        http.Error(w, "Error rendering template", http.StatusInternalServerError)
    }
}


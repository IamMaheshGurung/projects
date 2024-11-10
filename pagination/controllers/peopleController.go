

package controllers

import (
    "net/http"
    "log"
    "github.com/IamMaheshGurung/pagination/initializers"
    "github.com/IamMaheshGurung/pagination/models"
    "html/template"
    //"fmt"
    "strconv"
    "github.com/gorilla/mux"
    "math"
  

)

//Pagination Data


type PaginationData struct {
    NextPage int
    PreviousPage int
    CurrentPage int
    TotalPages int
}



func PeopleIndexGET(w http.ResponseWriter, r *http.Request) {

       
    //Get Page Number
    vars := mux.Vars(r)
    pageStr := vars["page"]
   
log.Printf("Received page query: '%s'", pageStr)
    perPage:=20
    // Default to page 1 if "page" is missing or invalid
    page, err := strconv.Atoi(pageStr)
    if err != nil || page < 1 {
        // Default to page 1 if invalid or missing
        page = 1
        log.Printf("Invalid or missing page number '%s', defaulting to page 1", pageStr)
    }
    
    //Calculating totalPages
    var totalRows int64
    initializers.DB.Model(&models.Person{}).Count(&totalRows)
    totalPages := math.Ceil(float64(totalRows / int64(perPage))) 
    log.Printf("Total Pages available are %d", int(totalPages))

    //calculating offset
    offset := (page-1) * perPage

    
    // Get the people from the database
    var people []models.Person
    if err := initializers.DB.Limit(perPage).Offset(offset).Find(&people).Error; err != nil {

        log.Printf("Error fetching people from the database: %v", err)
        http.Error(w, "Error fetching people data", http.StatusInternalServerError)
        return
    }

    // Parse all templates (top, index, bottom)
    tmpl := template.Must(template.New("index").ParseFiles(
               "templates/people/top.tmpl","templates/people/index.tmpl","templates/people/bottom.tmpl",
            ))
    /*if tmpl.Error != nil {
        log.Printf("Error parsing templates: %v", err)
        http.Error(w, "Error parsing the template", http.StatusInternalServerError)
        return
    }*/

    // Log the parsed templates
    log.Println("Successfully loaded templates")

    // Render the page with the people data
    err = tmpl.Execute(w, map[string]interface{}{
        "people": people,
        "pagination":PaginationData{
            NextPage : page + 1,
            PreviousPage: page -1,
            CurrentPage: page,
            TotalPages : int(totalPages),
        },
        
    })
    if err != nil {
        log.Printf("Error rendering template: %v", err)
        http.Error(w, "Error rendering template", http.StatusInternalServerError)
    }
}


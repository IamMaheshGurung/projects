package controllers



import (
    "log"
    "net/http"
    "html/template"


)


func RenderError(w http.ResponseWriter, errorMessage string){
    tmpl, err := template.ParseFiles("templates/error.html")
    if err != nil {
        log.Printf("Unable to render the error %s", err)
        return 
    }
   err = tmpl.Execute(w, errorMessage)
    if err != nil {
        log.Printf("Unable to execute the html")
        return 
    }
}


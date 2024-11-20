package controllers



import(
"net/http"
"html/template"
)



func HomePageDisplay(w http.ResponseWriter, r * http.Request){
    tmpl, err := template.ParseFiles("templates/index.html")
    if err != nil {
        http.Error(w, "Unable to load the homepage", http.StatusInternalServerError)
        return 
    }

    tmpl.Execute(w, nil)

}


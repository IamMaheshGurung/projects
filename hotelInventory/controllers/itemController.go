package controllers




import(
    "net/http"
    "html/template"
    "github.com/IamMaheshGurung/projects/hotelInventory/initializers"

    "github.com/IamMaheshGurung/projects/hotelInventory/models"

)



//Functions which will render inventory list

func ShowInventory(w http.ResponseWriter, r * http.Request){
    var items []models.Item


    result := initializers.DB.Find(&items)
    if result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }



//for HTML page rendering with templates
    tmpl, err := template.ParseFiles("templates/index.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return 
    }
    tmpl.Execute(w, struct {
        Items []models.Item
    }{Items: items})
}


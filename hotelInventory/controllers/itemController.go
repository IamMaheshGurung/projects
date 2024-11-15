package controllers




import(
    "net/http"
    "html/template"
    "github.com/IamMaheshGurung/projects/hotelInventory/initializers"

    "github.com/IamMaheshGurung/projects/hotelInventory/models"
    "strconv"

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


func CreateInventory(w http.ResponseWriter, r * http.Request) {
    if r.Method == http.MethodPost {
        name := r.FormValue("name")
        quantity := r.FormValue("quantity")

        Intqty, err  := strconv.Atoi(quantity)
        if err != nil {
             http.Error(w, "Unable to convert the string to int", http.StatusInternalServerError)
         }

        item := models.Item{Name :name, Quantity:Intqty }
        result := initializers.DB.Create(&item)
        if result.Error != nil {
            http.Error(w, "Item not found" ,http.StatusNotFound)
            return 
        }
        
        


      http.Redirect(w, r, "/", http.StatusFound)
  }
}

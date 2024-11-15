package controllers




import(
    "net/http"
    "html/template"
    "github.com/IamMaheshGurung/projects/hotelInventory/initializers"

    "github.com/IamMaheshGurung/projects/hotelInventory/models"
    "strconv"
    "log"

)



//Functions which will render inventory list

func ShowInventory(w http.ResponseWriter, r *http.Request) {
    var items []models.Item

    // Query database for items
    result := initializers.DB.Find(&items)
    if result.Error != nil {
        log.Printf("Error fetching items: %v", result.Error)
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }

    log.Printf("Items found: %d", len(items))

    // Parse and render template
    tmpl, err := template.ParseFiles("templates/index.html")
    if err != nil {
        log.Printf("Unable to parse template: %v", err)
        http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    err = tmpl.Execute(w, struct {
        Items []models.Item
    }{Items: items})

    if err != nil {
        log.Printf("Error executing template: %v", err)
        http.Error(w, "Execution error: "+err.Error(), http.StatusInternalServerError)
        return
    }
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

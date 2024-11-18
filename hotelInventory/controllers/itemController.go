package controllers




import(
    "net/http"
    "gorm.io/gorm"
    "html/template"
    "github.com/IamMaheshGurung/projects/hotelInventory/initializers"

    "github.com/IamMaheshGurung/projects/hotelInventory/models"
    "strconv"
    "github.com/gorilla/mux"

    "log"

)



//Functions which will render inventory list

func ShowInventory(w http.ResponseWriter, r *http.Request) {
    var items []models.Item
    var totalGuest []models.GuestLog

    // Query database for items
    result := initializers.DB.Find(&items)
    if result.Error != nil {
        log.Printf("Error fetching items: %v", result.Error)
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }


    allGuest := initializers.DB.Find(&totalGuest)
    if allGuest.Error != nil {
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
        Items    []models.Item
        TotalGuest []models.GuestLog
    }{
        Items:     items,
        TotalGuest: totalGuest,
    })
    if err != nil {
        log.Printf("Error executing template: %v", err)
        http.Error(w, "Execution error: "+err.Error(), http.StatusInternalServerError)
        return
    }
}


func CreateInventory(w http.ResponseWriter, r *http.Request) {
    // Show the form to create a new item when the method is GET
    if r.Method == http.MethodGet {
        tmpl, err := template.ParseFiles("templates/create.html")
        if err != nil {
            log.Printf("Unable to parse template: %v", err)
            http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
            return
        }
        err = tmpl.Execute(w, nil) // Render the form with no data (empty form)
        if err != nil {
            log.Printf("Error executing template: %v", err)
            http.Error(w, "Execution error: "+err.Error(), http.StatusInternalServerError)
            return
        }
    } else if r.Method == http.MethodPost {
        // Handle form submission
        name := r.FormValue("name")
        quantity := r.FormValue("quantity")

        Intqty, err := strconv.Atoi(quantity)
        if err != nil {
            log.Printf("Error converting quantity to int: %v", err)
            http.Error(w, "Unable to convert quantity", http.StatusInternalServerError)
            return
        
        }



        //checking if the item is already in the list

        var checkItem models.Item
        
        result := initializers.DB.Where("name=?", name).First(&checkItem)

        if result.Error != nil  && result.Error != gorm.ErrRecordNotFound{
            log.Printf("Error querying item not found %v", result.Error)
            http.Error(w, "Error checking item", http.StatusInternalServerError)
            return
        }

        if result.RowsAffected > 0 {
            checkItem.Quantity += Intqty
            result := initializers.DB.Save(&checkItem)
            if result.Error != nil {
                log.Printf("Error updating item: %v", result.Error)
                http.Error(w, "Error updating item", http.StatusInternalServerError)
                return
            }

             } else {

    
        item := models.Item{Name: name, Quantity: Intqty}
        result := initializers.DB.Create(&item)
        if result.Error != nil {
            log.Printf("Error creating item: %v", result.Error)
            http.Error(w, "Item creation failed", http.StatusInternalServerError)
            return
        }
    }

        // Redirect to the inventory list after creating the item
        http.Redirect(w, r, "/", http.StatusFound)
    }
}


func EditInventory(w http.ResponseWriter, r * http.Request) {
    v := mux.Vars(r)
    itemID := v["id"]
    
    if itemID == "" {
        http.Error(w, "Item ID is required", http.StatusBadRequest)
        return
    }

    var item models.Item
    result := initializers.DB.First(&item, itemID)
    
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            http.Error(w, "Item not found", http.StatusNotFound)
        } else {
            log.Printf("Error fetching item: %v", result.Error)
            http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        }
        return
    }


    if r.Method == http.MethodGet{
        tmpl, err := template.ParseFiles("templates/edit.html")
        if err != nil {
            log.Printf("Unable to parse the edit.html file:%s", err)
            http.Error(w, "Unable to parse the file", http.StatusInternalServerError)
            return 
        }

        err = tmpl.Execute(w,item )
        if err != nil {
            log.Printf("Unable to execute the form %s", err)
            http.Error(w, "Unable to execute the file", http.StatusInternalServerError)
            return 
        }

    } else if r.Method == http.MethodPost {
        name := r.FormValue("name")
        quantity := r.FormValue("quantity")


        intQty, err := strconv.Atoi(quantity)
            if err != nil {
                log.Printf("Unable to convert the quantiy string into inteeger, %s", err)
                return 
            }

            item.Name = name
            item.Quantity = intQty
            result := initializers.DB.Save(&item)
            if result.Error != nil {
            log.Printf("Error updating item: %v", result.Error)
            http.Error(w, "Error updating item", http.StatusInternalServerError)
            return
        }

        // Redirect back to the inventory list
        http.Redirect(w, r, "/", http.StatusFound)
    }
}




// DeleteItem handles the DELETE request to remove an item
func ShowDeletePage(w http.ResponseWriter, r *http.Request) {
    // Extract item ID from URL
    vars := mux.Vars(r)
    itemID := vars["id"]

    if itemID == "" {
        http.Error(w, "Item ID is required", http.StatusBadRequest)
        return
    }

    var item models.Item
    result := initializers.DB.First(&item, itemID)
    
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            http.Error(w, "Item not found", http.StatusNotFound)
        } else {
            log.Printf("Error fetching item: %v", result.Error)
            http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        }
        return
    }

    // Render the confirmation template with the item details
    tmpl, err := template.ParseFiles("templates/delete.html")
    if err != nil {
        log.Printf("Error parsing template: %v", err)
        http.Error(w, "Error parsing template", http.StatusInternalServerError)
        return
    }

    err = tmpl.Execute(w, struct {
        ID string
    }{
        ID: itemID,
    })
    if err != nil {
        log.Printf("Error executing template: %v", err)
        http.Error(w, "Execution error: "+err.Error(), http.StatusInternalServerError)
    }
}

// DeleteItem handles the deletion of an item
func DeleteItem(w http.ResponseWriter, r *http.Request) {
    // Extract item ID from URL
    vars := mux.Vars(r)
    itemID := vars["id"]

    if itemID == "" {
        http.Error(w, "Item ID is required", http.StatusBadRequest)
        return
    }

    // Find the item by ID
    var item models.Item
    result := initializers.DB.First(&item, itemID)

    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            http.Error(w, "Item not found", http.StatusNotFound)
        } else {
            log.Printf("Error fetching item: %v", result.Error)
            http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        }
        return
    }

    // If the user confirmed deletion, delete the item
    if r.Method == http.MethodPost {
        confirm := r.FormValue("confirm")
        if confirm == "yes" {
            result = initializers.DB.Delete(&item)
            if result.Error != nil {
                log.Printf("Error deleting item: %v", result.Error)
                http.Error(w, "Error deleting item", http.StatusInternalServerError)
                return
            }
        http.Redirect(w, r, "/", http.StatusFound)

        }

        // Redirect to the inventory list after deletion or cancellation
        http.Redirect(w, r, "/", http.StatusFound)
    }
}




func Createuser(w http.ResponseWriter, r *http.Request) {
    // Show the form to create a new user when the method is GET
    if r.Method == http.MethodGet {
        tmpl, err := template.ParseFiles("templates/user.html")
        if err != nil {
            log.Printf("Unable to parse template: %v", err)
            http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
            return
        }
        err = tmpl.Execute(w, nil) // Render the form with no data (empty form)
        if err != nil {
            log.Printf("Error executing template: %v", err)
            http.Error(w, "Execution error: "+err.Error(), http.StatusInternalServerError)
            return
        }
    } else if r.Method == http.MethodPost {
        // Handle form submission
        email := r.FormValue("email")
        password := r.FormValue("password")

    
        user := models.User{Email: email, Password: password}
        result := initializers.DB.Create(&user)
        if result.Error != nil {
            log.Printf("Error creating item: %v", result.Error)
            http.Error(w, "User creation failed", http.StatusInternalServerError)
            return
        }
    }

        // Redirect to the inventory list after creating the item
        http.Redirect(w, r, "/", http.StatusFound)
    }






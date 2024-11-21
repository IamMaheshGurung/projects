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



// ShowInventory is the handler for displaying inventory items and guest logs
func ShowInventory(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r)  // Assuming GetUser is a function that retrieves the user from the request context
	if user == nil {
		log.Println("User not found or unauthorized in the inventory")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Query for the inventory items for the user
	var items []models.Item
	result := initializers.DB.Where("user_id = ?", user.ID).Find(&items)
	if result.Error != nil {
		log.Printf("Error fetching items for user %d: %v", user.ID, result.Error)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// If no items are found, log the message
	if len(items) == 0 {
		log.Printf("No inventory items found for user %d", user.ID)
	}

	// Query for the guest logs (total guests for the user)
	var totalGuest []models.GuestLog
	result = initializers.DB.Where("user_id = ?", user.ID).Find(&totalGuest)
	if result.Error != nil {
		log.Printf("Error fetching guest logs for user %d: %v", user.ID, result.Error)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// If no guest logs are found, log the message
	if len(totalGuest) == 0 {
		log.Printf("No guest logs found for user %d", user.ID)
	}

	// Parse the HTML template
	tmpl, err := template.ParseFiles("templates/inventory.html")
	if err != nil {
		log.Printf("Unable to parse template: %v", err)
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template with the fetched data
	err = tmpl.Execute(w, struct {
		Items      []models.Item
		TotalGuest []models.GuestLog
	}{
		Items:      items,      // Ensure items is not nil
		TotalGuest: totalGuest, // Ensure totalGuest is not nil
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
        http.Redirect(w, r, "/inventory", http.StatusFound)

        }
    }
}





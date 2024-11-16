package controllers




import(
    "net/http"
    "gorm.io/gorm"
    "html/template"
    "github.com/IamMaheshGurung/projects/hotelInventory/initializers"

    "github.com/IamMaheshGurung/projects/hotelInventory/models"
    "strconv"
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

// AutoUpdateInventory function: Automatically updates inventory based on guest count
func AutoUpdateInventory(w http.ResponseWriter, r *http.Request) {
    // Define items per guest (for simplicity, assume 1 sheet and 1 pillow cover per guest)
    itemsPerGuest := map[string]int{
        "sheet":        1,
        "pillow_cover": 1,
    }

    // Retrieve the total number of guests (from the `GuestLog` table)
    var totalGuest []models.GuestLog
    result := initializers.DB.Find(&totalGuest)
    if result.Error != nil {
        log.Printf("Error fetching guest data: %v", result.Error)
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }

    // Total number of guests
    totalGuestsCount := 0
    for _, guest := range totalGuest {
        totalGuestsCount += guest.TotalGuest
    }

    // Log the total number of guests
    log.Printf("Total guests: %d", totalGuestsCount)

    // Loop through each inventory item and adjust based on guest count
    for itemName, quantityPerGuest := range itemsPerGuest {
        // Calculate the total number of items to decrement based on guest count
        totalItemsToAdjust := totalGuestsCount * quantityPerGuest

        // Fetch the current inventory for the item
        var inventoryItem models.Item
        err := initializers.DB.Where("name = ?", itemName).First(&inventoryItem).Error
        if err != nil {
            log.Printf("Error fetching inventory for %s: %v", itemName, err)
            continue
        }

        // Decrease the inventory by the required amount
        if inventoryItem.Quantity >= totalItemsToAdjust {
            // If there are enough items, decrement them
            newQuantity := inventoryItem.Quantity - totalItemsToAdjust
            initializers.DB.Model(&inventoryItem).Update("quantity", newQuantity)
            log.Printf("Decreased inventory for %s by %d", itemName, totalItemsToAdjust)
        } else {
            // If not enough inventory, set the quantity to 0 (or handle based on business logic)
            initializers.DB.Model(&inventoryItem).Update("quantity", 0)
            log.Printf("Not enough %s in inventory. Setting quantity to 0.", itemName)
        }
    }

    // Respond with a success message
    w.Write([]byte("Inventory updated based on guest count"))
}





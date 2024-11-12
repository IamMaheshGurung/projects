package controllers




import(
    "encoding/json"
    "net/http"
    "time"
    "github.com/IamMaheshGurung/projects/hotelInventory/models"
    "github.com/IamMaheshGurung/projects/hotelInventory/initializers"
)


func AccountSet(w http.ResponseWriter, r * http.Request){
    var items =[]string{"bedsheet", "pillow cover", "blanket cover"}

    for _, itemType := range items {
        var inventory models.Inventory
        if err := initializers.DB.Where("item_type = ?", itemType).First(&inventory).Error; err != nil {
            http.Error(w, itemType + "not found", http.StatusNotFound)
            return 
        }

        if inventory.Quantity < 1 {
            http.Error(w, "not enough" + itemType, http.StatusBadRequest)
            return 
        }

        inventory.Quantity -= 1

        initializers.DB.Save(&inventory)
    }

    var logy models.UsageLog

    initializers.DB.First(&logy)
    logy.TotalGuest += 1
    logy.TimeStamp = time.Now().Format("2006-01-02 15:04:05")
    initializers.DB.Save(&logy)

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode("Stock updated and customer count updated")
}

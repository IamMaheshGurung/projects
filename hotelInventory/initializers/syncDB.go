package initializers

import(
    "log"
    "github.com/IamMaheshGurung/projects/hotelInventory/models"
)

func SyncDB(){
    err := DB.AutoMigrate(&models.UsageLog{}, &models.Inventory{})
    if err != nil {
        log.Fatal("Unable to automigrate in database:" , err)
    }

}



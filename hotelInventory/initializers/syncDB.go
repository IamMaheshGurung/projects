package initializers

import(
    "log"
    "github.com/IamMaheshGurung/projects/hotelInventory/models"
)

func SyncDB(){
    err := DB.AutoMigrate(&models.Item{}, &models.GuestLog{}, &models.User{})
    if err != nil {
        log.Fatal("Unable to automigrate in database:" , err)
    }

}



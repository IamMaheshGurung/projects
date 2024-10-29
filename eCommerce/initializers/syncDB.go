package initializers




import(
    /*"gorm.io/gorm"
    "gorm.io/driver/postgres"*/
    "github.com/IamMaheshGurung/eCommerce/models"
    "log"


)



func SyncDB() {
    err := DB.AutoMigrate(&models.User{})
    if err != nil {
        log.Fatal("failed to migrate date:", err)
    }
    log.Println("Database has been connected sucessfully")
}


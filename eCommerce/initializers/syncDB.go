package initializers



import(
    "github.com/IamMaheshGurung/ecommerce/models"
    "log"


)



func SyncDB() {
    err := DB.AutoMigrate(&models.User{})

    if err != nil {
        log.Fatalf("Failed to migrate the data : %s",err)
    }
    log.Println("Database has been connected successfully")

}


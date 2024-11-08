package initializers



import (
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
    "fmt"
    "os"
    "github.com/IamMaheshGurung/pagination/models"


)



var DB *gorm.DB



func SyncDB(){
    DB.AutoMigrate(&models.Person{})

}


func ConnectToDB(){
    var err error 

    dsn := os.Getenv("DSN")
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

    if err != nil {
        fmt.Println("Failed to connect to db")
    }
}

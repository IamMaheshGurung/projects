package initializers



import (
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
    "fmt"


)



var DB *gorm.DB



func SyncDB(){
    DB.AutoMigrate(&models.Person{})

}


func ConnectToDB(){
    var err error 


    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

    if err != nil {
        fmt.Println("Failed to connect to db")
    }
}

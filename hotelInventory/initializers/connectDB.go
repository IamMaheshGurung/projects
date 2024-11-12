package initializers


import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "log"
    "os"
)



var DB *gorm.DB

func ConnectDB(){
    dsn := os.Getenv("CONSTR")
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to the database: ",err)
    }
    DB = db
}

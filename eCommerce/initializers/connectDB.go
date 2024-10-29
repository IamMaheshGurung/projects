package initializers




import (
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
    "os"
    "log"
)

var DB *gorm.DB

func InitDB(){
    var err error

    dsn := os.Getenv("CONSTR")
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect DB")
    }
    log.Println("Connection to database has been successful")
}



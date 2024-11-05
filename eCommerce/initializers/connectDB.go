package initializers


import 
(
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
    "os"
    "log"
)


var DB *gorm.DB



func InitDB(){
    var err error

    dsn := os.Getenv("CONNSTR")

    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

    if err != nil {
        log.Fatalf("Unble to Connect the DB:%s", err)
    }

    log.Println("Has been connected to database")
}


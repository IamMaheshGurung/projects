package models



import(
    "gorm.io/gorm"
    //"gorm.io/driver/postgres"
)




type User struct {
    gorm.Model
    Email string `gorm:"unique"`
    Password string
}


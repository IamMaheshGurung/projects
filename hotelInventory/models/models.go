package models



import(
    "gorm.io/gorm"
)



type Item struct {
    gorm.Model
    Name     string `json:"name"`
    Quantity int    `json:"quantity"`
    UserID   uint   `json:"user_id"`
}



type GuestLog struct {
    gorm.Model
    UserID uint    `gorm:"not null"`     
    User   User    `gorm:"foreignKey:UserID"` 
}

type User struct {
    gorm.Model
    Email string `json:"email"`
    Password string `json:"password"`
    Item []Item  `gorm:"foreignkey:UserID"`
}

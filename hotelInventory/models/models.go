package models



import(
    "gorm.io/gorm"
)


type Item struct {
    gorm.Model
    Name string `json:"name"`
    Quantity int `json:"quantity"`
}


type GuestLog struct {
    gorm.Model
    TimeStamp string `json:"timesstamp"`
    TotalGuest int `json:"total_guest"`
}



package models



import(
    "gorm.io/gorm"
)


type Inventory struct {
    gorm.Model
    ItemType string `json:"item_type"`
    Quantity int `json:"quantity"`
}


type UsageLog struct {
    gorm.Model
    TimeStamp string `json:"timesstamp"`
    TotalGuest int `json:"total_guest"`
}



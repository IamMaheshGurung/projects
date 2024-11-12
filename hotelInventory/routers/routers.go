package routers





import(
    "github.com/IamMaheshGurung/projects/hotelInventory/controllers"
    "github.com/gorilla/mux"

)


func SetupRouters(router *mux.Router){
    router.HandleFunc("/inventory/setstock", controllers.AccountSet).Methods("POST")
}


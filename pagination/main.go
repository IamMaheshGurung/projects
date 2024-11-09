package main 



import(
    "net/http"
    "time"
    "log"
    "fmt"
    "github.com/IamMaheshGurung/pagination/initializers"
    
    "github.com/IamMaheshGurung/pagination/controllers"
    "github.com/gorilla/mux"
)




func init(){

    //initializers file


    initializers.LoadEnvVariables()
    initializers.ConnectToDB()
    initializers.SyncDB()
    initializers.CreatePeople()

}




func main(){
   router := mux.NewRouter()

    // Set up routes
    router.HandleFunc("/page", controllers.PeopleIndexGET).Methods("GET")
    router.HandleFunc("/page/{page:[0-9]+}", controllers.PeopleIndexGET).Methods("GET")

    server := http.Server{
        Addr : ":8080",
        Handler : router,
        IdleTimeout: 120 * time.Second,
    }
           fmt.Printf("Server is running at local host %s", server.Addr)
        err := server.ListenAndServe()
        if err != nil {
            log.Printf("Unable to connect to the server")
            
        }


}

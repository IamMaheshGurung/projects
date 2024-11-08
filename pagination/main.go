package main 



import(
    "net/http"
    "time"
    "github.com/IamMaheshGurung/pagination/initializers"
    
    "github.com/IamMaheshGurung/pagination/controllers"
)




func init(){

    //initializers file


    initializers.LoadEnvVariables()
    initializers.ConnectToDB()
    initializers.SyncDB()
    initializers.CreatePeople()

}




func main(){
    l := http.NewServeMux()
    l.HandleFunc("/", controllers.PeopleIndexGET)



    server := http.Server{
        Addr : ":8080",
        Handler : l,
        IdleTimeout: 120 * time.Second,
    }

    server.ListenAndServe()

}

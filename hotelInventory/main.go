package main



import (
    "net/http"
    "github.com/gorilla/mux"
    "time"
    "os"
    "os/signal"
    "context"
    "log"
    "github.com/IamMaheshGurung/projects/hotelInventory/initializers"
    


)




func init(){
initializers.LoadEnvVariables()

initializers.ConnectDB()
initializers.SyncDB()
}




func main(){

    router := mux.NewRouter()

    port:= os.Getenv("PORT")

    server := http.Server{
        Addr : ":"+ port,
        Handler : router,
        IdleTimeout: 120 * time.Second,
    }


    go func(){
         log.Printf("Server is running at %s", port)
        err := server.ListenAndServe()
        if err != nil {
            log.Fatalf("Unable to connect to the server%s", err)
            return 
        }
        
       
    }()


    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Kill, os.Interrupt)
    <- stop


    ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
    defer cancel()

    log.Println("GRACEFULL SHUTDOWN HAS BEEN REQUESTED WHAT TO DO")
    err := server.Shutdown(ctx)
    if err != nil{
        log.Printf("Unable to shutdown the server :%s", err)
        return 
    }
    log.Println("Server has been shutdown gracefuly")


}

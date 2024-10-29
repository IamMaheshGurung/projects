package main 



import (
    "net/http"
    "log"
    "time"
    "github.com/IamMaheshGurung/eCommerce/initializers"
    "github.com/IamMaheshGurung/eCommerce/controller"
    "github.com/IamMaheshGurung/eCommerce/middleware"
    "os"
    "os/signal"
    "fmt"
    "context"
)



func init(){
    initializers.LoadEnvVariables()
    initializers.InitDB()
    initializers.SyncDB()

}



func main(){

    http.HandleFunc("/signup", controllers.Signup)
    http.HandleFunc("/login", controllers.Login)
    http.Handle("/validate", middleware.RequireAuth(http.HandlerFunc(controllers.Validate)))



    port := ":" + os.Getenv("PORT")
    if port == ":" + ""{
        port = ":8080"
    }

    server := http.Server{
        Addr :  port,
        IdleTimeout: 120 * time.Second,
    }

    go func(){
        fmt.Printf("Server is running at port %s", port)
        err := server.ListenAndServe()
        if err != nil {
            log.Fatalf("Sorry unable to serve %s", err)
        }
    }()


    stop := make(chan os.Signal, 1)

    signal.Notify(stop, os.Kill, os.Interrupt)

    <- stop

    ctx, cancel := context.WithTimeout(context.Background(), 25 * time.Second)
    defer cancel()



    fmt.Println("Graceful shutdown has been requested")

    err := server.Shutdown(ctx)
    if err != nil{
        log.Fatalf("Unable to shutdown the server %s", err)
    }

    fmt.Println("server has been shutdown gracefully")
}








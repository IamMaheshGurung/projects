package main 



import(
"net/http"
"fmt"
"log"
"github.com/IamMaheshGurung/ecommerce/initializers"
"github.com/IamMaheshGurung/ecommerce/controllers"
"time"
"os"
"os/signal"
"context")




func init(){


    initializers.LoadEnvVariables()


}






func main(){
    handler := http.NewServeMux()
    handler.HandleFunc("/", controllers.Roothandler)
    port := os.Getenv("MYPORT")
    if port == "" {
        port = ":8080"
    }


    server := http.Server{
        Addr: port,
        Handler: handler,
        IdleTimeout: 120 *  time.Second,
    }

    go func(){
        fmt.Printf("Server is running at port %s", port)
        err := server.ListenAndServe()
        if err != nil {
            log.Fatalf("Sorry Unable to serve %s", err)
        }
    }()


    stop := make(chan os.Signal, 1)

    signal.Notify(stop, os.Kill, os.Interrupt)

    <-stop


    ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)

    defer cancel()

    fmt.Println("Request for graceful shutdown has been done")
    err := server.Shutdown(ctx)
    if err != nil {
        log.Fatalf("Unable to shutdown the server %s", err)
    }
    fmt.Println("Server has been shutdown gracefully")






}

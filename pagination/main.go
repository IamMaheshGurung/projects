package main 



import(
    "net/http"
    "time"
)




func init(){

    //initializers file




}




func main(){
    l := http.NewServeMux()




    server := http.Server{
        Addr : ":8080",
        Handler : l,
        IdleTimeout: 120 * time.Second,
    }

    server.ListenAndServe()

}

 package controllers



import(
    "fmt"
    "net/http"
    "log"
    "github.com/IamMaheshGurung/eCommerce/models"
    "encoding/json"
    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v5"
    "github.com/IamMaheshGurung/eCommerce/initializers"
    "github.com/IamMaheshGurung/eCommerce/middleware"
    "time"
    "os"
)
func RootHandler (w http.ResponseWriter, r *http.Request){
    switch r.URL.Path {
    case "/validate":
        if r.Method == http.MethodGet {
            middleware.RequireAuth(http.HandlerFunc(Validate)).ServeHTTP(w, r)
        } else {
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        }

    case "/login":
        if r.Method == http.MethodPost{
             Login(w,r)
         } else {
             http.Error(w, "Method Not Allowed", http.StatusBadRequest)
         }

    case "/signup":
        if r.Method == http.MethodPost{
             Signup(w,r)
         } else {
             http.Error(w, "Method Not Allowed", http.StatusBadRequest)
         }


    

    default:
        fmt.Println("Invalid Request")
        return
    }
}



func Signup(w http.ResponseWriter, r *http.Request){
    defer r.Body.Close()


    //Retriving the info from the request by user

    var body struct{
        Email string    
        Password string 
    }

    err := json.NewDecoder(r.Body).Decode(&body)

    if err != nil {
        log.Printf("Failed to decode the info from the body")
    }

    //Encrypting or hashing the password

    hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
    if err != nil {
        log.Printf("Unable to hash the password")
    }
    

    user := models.User{Email: body.Email, Password: string(hash)}
    result := initializers.DB.Create(&user)
    if result.Error != nil{
        log.Printf("Failed to create the user here")
    }
    w.WriteHeader(http.StatusOK)
}





func Login(w http.ResponseWriter, r *http.Request){

    var body struct {
        Email string    
        Password string 
    }

    err := json.NewDecoder(r.Body).Decode(&body)
    if err != nil {
        log.Printf("Failed to decode the body")
    }

    var user models.User

    initializers.DB.First(&user, "email=?", body.Email)
    if user.ID==0{
        http.Error(w, "Incorrect email", http.StatusUnauthorized)
        return
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
    if err != nil {
        http.Error(w, "Incorrect password or email password didnot matched", http.StatusUnauthorized)
        return 
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims {
        "sub": user.ID,
        "exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
    })


    tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

    if err != nil {
        http.Error(w, "failed to create the token", http.StatusInternalServerError)
        return
    }
    
    cookie := http.Cookie{
        Name : os.Getenv("NAME"),
        Value : tokenString,
        Expires: time.Now().Add(24 * time.Hour),
        HttpOnly: true,
        Secure: true,
        SameSite : http.SameSiteLaxMode,
    }


     http.SetCookie(w, &cookie)
     w.WriteHeader(http.StatusOK)
    w.Write([]byte("Cookies are set successfully"))

}





func Validate (w http.ResponseWriter, r *http.Request){
    message := map[string]string{"message":"I am logged in now"}
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    json.NewEncoder(w).Encode(message)
}




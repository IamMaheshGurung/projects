package controllers




import(
    "net/http"
    "github.com/IamMaheshGurung/ecommerce/models"
    "github.com/golang-jwt/jwt/v5"
    "github.com/IamMaheshGurung/ecommerce/initializers"
    "github.com/IamMaheshGurung/ecommerce/middleware"
    "golang.org/x/crypto/bcrypt"
    "encoding/json"
    "log"
    "time"
    "os"





)





func Roothandler(w http.ResponseWriter, r * http.Request){
    switch r.Method {
    case http.MethodPost:
        if r.URL.Path == "/login"{
            Login(w,r)
        } else if r.URL.Path == "/signup" {
            Signup(w,r)
        } else {
            http.Error(w, "Invalid url path", http.StatusBadRequest)
            return 
        }

    case http.MethodGet:
        if r.URL.Path=="/validate"{

            middleware.RequireAuth(http.HandlerFunc(Validate)).ServeHTTP(w,r)
        }
    }








}



//creating signup handler 
func Signup(w http.ResponseWriter, r *http.Request){
    defer r.Body.Close()


    //need to retrive data frm reuwst by user 



    var body struct {
        Email string
        Password string

    }


    err := json.NewDecoder(r.Body).Decode(&body)
    if err != nil {
        http.Error(w, "Unable to decode the body:%s", http.StatusInternalServerError)
        return 
    }


    hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 15)
    if err != nil {
        http.Error(w, "Unable to generate encrypted password", http.StatusInternalServerError)
        return 
    }

    user := models.Users{Email : body.Email, Password: string(hash),}
    result := initializers.DB.Create(&user)
    if result.Error!= nil {
        log.Printf("Failed to create the user here")
    }
    w.WriteHeader(http.StatusOK)

}




func Login(w http.ResponseWriter, r * http.Request){
    var body struct {
        Email string 
        Password string
    }

//here getting input from the url and decoding it
    err := json.NewDecoder(r.Body).Decode(&body)
    if err != nil {
        log.Printf("Failed to deccode the body")
    }


    //declaring the variable for storing the received and decoded data
    var user models.Users
    

    //searching in the database if the user exist

    initializers.DB.First(&user, "email=?", body.Email)
    if user.ID == 0 {
        http.Error(w, "Incorrect Email", http.StatusUnauthorized)
        return 
    }
    

    // if the user exist here then we compare password with the already hashed password

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
    if err != nil {
        http.Error(w, "Incorrect password or email didnot matched", http.StatusUnauthorized)
        return 
    }
    


    //if everything gets right we get token  with the below method

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims {
        "sub" : user.ID,
        "exp" : time.Now().Add(time.Hour * 24 * 30).Unix(),
    })
//stringify my token with the secret code that I have in my env file
    tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

    if err != nil {
        http.Error(w, "Failed to create the token", http.StatusInternalServerError)
        return
    }
//if everything goes right then lets form the cookie..
    cookie := http.Cookie {
        Name : os.Getenv("NAME"),
        Value : tokenString,
        Expires : time.Now().Add(24 * time.Hour),
        HttpOnly : true,
        Secure : true,
        SameSite : http.SameSiteLaxMode,
    }

    //with this lets set the cookie and with everything ok lets send the feedback ok

    http.SetCookie(w, &cookie)
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Cookies are set successfully"))
}


func Validate(w http.ResponseWriter, r*http.Request){
    message :=map[string]string{"message":"I am Logged in with the validation now"}
    w.Header().Set("Content-Type", "application/json")

    json.NewEncoder(w).Encode(message)
}

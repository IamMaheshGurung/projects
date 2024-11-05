package controllers




import(
    "net/http"
    "github.com/IamMaheshGurung/ecommerce/models"
    "github.com/IamMaheshGurung/ecommerce/initializers"
    "golang.org/x/crypto/bcrypt"
    "encoding/json"
    "log"




)





func Roothandler(w http.ResponseWriter, r * http.Request){








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


    err := json.NewDecoder(r.Body).Decode(&body)
    if err != nil {
        log.Printf("Failed to deccode the body")
    }

    var user models.User

    initializers.DB.First(&user, "email=?", body.Email)
    if user.ID == 0 {
        http.Error(w, "Incorrect Email", http.StatusUnauthorized)
        return 
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
    if err != nil {
        http.Error(w, "Incorrect password or email didnot matched", http.StatusUnauthorized)
        return 
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims {
        "sub" : user.ID,
        "exp" : time.Now().Add(time.Hour * 24 * 30).Unix(),
    })

    tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

    if err != nil {
        http.Error(w, "Failed to create the token", http.StatusInternalServerError)
        return
    }

    cookie := http.Cookie {
        Name : os.Getenv("NAME"),
        Value : tokenString,
        Expires : time.Now().Add(24 * time.Hour),
        HttpOnly : true,
        Secure : true,
        SameSite : http.SameSiteLaxMode,
    }

    http.SetCookie(w, &cookie)
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Cookies are set successfully"))
}


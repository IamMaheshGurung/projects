package controllers




import(
    "net/http"
    //"encoding/json"
    "os"
    "fmt"
    "time"
    "github.com/golang-jwt/jwt/v5"
    //"gorm.io/gorm"
    "html/template"
    "github.com/IamMaheshGurung/projects/hotelInventory/initializers"
    "golang.org/x/crypto/bcrypt"
    "github.com/IamMaheshGurung/projects/hotelInventory/models"
    //"strconv"
    //"github.com/gorilla/mux"

    "log"

)

func GetUser(r *http.Request) *models.User {
    user, ok := r.Context().Value("user").(*models.User)
    if !ok {
        fmt.Println("User in the context not found")
        return nil 
    }
    return user 
}

func Signup(w http.ResponseWriter, r *http.Request) {
    // Show the form to create a new user when the method is GET
    if r.Method == http.MethodGet {
        tmpl, err := template.ParseFiles("templates/user.html")
        if err != nil {
            log.Printf("Unable to parse template: %v", err)
            http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
            return
        }


        err = tmpl.Execute(w, nil) // Render the form with no data (empty form)
        if err != nil {
            log.Printf("Error executing template: %v", err)
            http.Error(w, "Execution error: "+err.Error(), http.StatusInternalServerError)
            return
        }
    } else if r.Method == http.MethodPost {
        // Handle form User not found or unauthorizedsubmission
        email := r.FormValue("email")
        password := r.FormValue("password")

        hash, err := bcrypt.GenerateFromPassword([]byte(password), 15)
        if err != nil {
            log.Printf("Unable to hash the password")
        }

        user:= models.User{Email:email, Password:string(hash)}

    
        result := initializers.DB.Create(&user)
        if result.Error != nil {
            log.Printf("Error creating item: %v", result.Error)
            http.Error(w, "User creation failed", http.StatusInternalServerError)
            return
        }
    }

        // Redirect to the inventory list after creating the item
        http.Redirect(w, r, "/", http.StatusFound) 
    }

func Login(w http.ResponseWriter, r *http.Request) {
    // Show the form to create a new user when the method is GET
    if r.Method == http.MethodGet {
        tmpl, err := template.ParseFiles("templates/login.html")
        if err != nil {
            log.Printf("Unable to parse template: %v", err)
            http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
            return
        }


        err = tmpl.Execute(w, nil) // Render the form with no data (empty form)
        if err != nil {
            log.Printf("Error executing template: %v", err)
            http.Error(w, "Execution error: "+err.Error(), http.StatusInternalServerError)
            return
        }
    } else if r.Method == http.MethodPost {
        // Handle form submission
        email := r.FormValue("email")
        password := r.FormValue("password")



        var user models.User

        initializers.DB.First(&user, "email = ? ", email)

        if user.ID == 0 {
            http.Error(w , "Invalid email or cannot be found", http.StatusNotFound)
            return 
        }



        

        err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

        if err != nil {

            RenderError(w, "Invalid email or password")
            http.Error(w, "Invalid Password or email id", http.StatusUnauthorized)
            log.Printf("Invaild password")
            return
        }

        token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
            "sub" : user.ID,
            "exp" : time.Now().Add(time.Hour * 24 * 30).Unix(),
        })

        tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

        if err != nil {
            http.Error(w, "Failed to create token", http.StatusInternalServerError)
            return 
        }

        
        cookie := http.Cookie{
            Name : "Authorization",
            Value: tokenString,
            Expires: time.Now().Add(24*time.Hour),
            HttpOnly: true,
            Secure: true,
            SameSite: http.SameSiteLaxMode,
        }

        http.SetCookie(w, &cookie)
        
    
        // Redirect to the inventory list after creating the item
        http.Redirect(w, r, "/inventory", http.StatusFound)
    }
    }






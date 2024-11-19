package controllers




import(
    "net/http"
    "encoding/json"
    //"gorm.io/gorm"
    "html/template"
    "github.com/IamMaheshGurung/projects/hotelInventory/initializers"
    "golang.org/x/crypto/bcrypt"
    "github.com/IamMaheshGurung/projects/hotelInventory/models"
    //"strconv"
    //"github.com/gorilla/mux"

    "log"

)



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
        // Handle form submission
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
        err = json.NewEncoder(w).Encode(result)
        if err!= nil {
            log.Printf("Unable to encode the json %s", err)
            return
        }
    }

        // Redirect to the inventory list after creating the item
        http.Redirect(w, r, "/", http.StatusFound) 
    }






package middleware


import(

    "net/http"
    "os"
    "time"
    "log"
    "fmt"
    "github.com/golang-jwt/jwt/v5"
    "context"
    "github.com/IamMaheshGurung/projects/hotelInventory/models"
    "github.com/IamMaheshGurung/projects/hotelInventory/initializers"
)



type contextKey string 


const userContextKey = contextKey("user")

func RequireAuth(next http.Handler) http.Handler{
    return http.HandlerFunc(func(w http.ResponseWriter, r * http.Request){
        tokenCookie, err := r.Cookie("Authorization")
        if err != nil {
            log.Printf("Unable to get the token string")
        }
        tokenString  := tokenCookie.Value

        token , err := jwt.Parse(tokenString, func(token *jwt.Token)(interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("Unexpected signing method : %v", token.Header["alg"])
            }

            return []byte(os.Getenv("SECRET")), nil 
        })
        if err != nil {
            log.Fatal(err)
        }

        if claims, ok := token.Claims.(jwt.MapClaims); ok {
            if float64(time.Now().Unix())> claims["exp"].(float64) {
                log.Printf("The token has been expired %s", err)
            }

            var user models.User

            initializers.DB.First(&user, claims["sub"])

            if user.ID == 0 {
                log.Printf("Failed to get the sub")
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return 
            }

            ctx:= context.WithValue(r.Context(), userContextKey, &user)
            next.ServeHTTP(w, r.WithContext(ctx))

        }

        http.Error(w , "Unauthorized", http.StatusUnauthorized)



    })
}

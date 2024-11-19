package middleware


import(

    "net/http"
    "os"
    "time"
    "log"
    "fmt"
    "github.com/golang-jwt/jwt/v5"
    "context"
)



type contextKey string 


const userContextKey = contextKey("user")

func RequireAuth(next http.Handler) http.Handler{
    return http.HandlerFunc(func(w http.ResponseWriter, r * http.Request){
        tokenCookie, err := r.Cookie("Authorization")
        if err != nil {
            log.Printf("Unable to get the token string")
        }
        tokenString  := token.Value

        token , err := jwt.Parse(tokenString, func(token *jwt.Token)(interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("Unexpected signing method : %v", token.Header["alg"])
            }

            return []byte(os.Getenv("SECRET")), nil 
        })
        if err != nil {
            log.Fatal(err)
        }

    }
}

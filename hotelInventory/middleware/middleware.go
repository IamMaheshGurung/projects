

package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/IamMaheshGurung/projects/hotelInventory/models"
	"github.com/IamMaheshGurung/projects/hotelInventory/initializers"
    "github.com/IamMaheshGurung/projects/hotelInventory/controllers"
)

type contextKey string

const userContextKey = contextKey("user")

// RequireAuth middleware to check for JWT token and validate the user

func RequireAuth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tokenCookie, err := r.Cookie("Authorization")
        if err != nil {
            log.Printf("Unable to get the token from cookie: %v", err)
            controllers.RenderError(w, "Unable to get the token from the coookie")
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        tokenString := tokenCookie.Value
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return []byte(os.Getenv("SECRET")), nil
        })

        if err != nil {
            log.Printf("Error parsing token: %v", err)
            controllers.RenderError(w, "Error in parsing the template")
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        if claims, ok := token.Claims.(jwt.MapClaims); ok {
            expirationTime := claims["exp"].(float64)
            if float64(time.Now().Unix()) > expirationTime {
                log.Printf("Token expired")
                controllers.RenderError(w, "token has been expired")
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }

            userID := int64(claims["sub"].(float64)) // Convert to int64
            var user models.User
            result := initializers.DB.First(&user, userID)

            if result.Error != nil || user.ID == 0 {
                log.Printf("User not found or unauthorized: %v", result.Error)
                controllers.RenderError(w, "User not found or password incorrect")
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }

            log.Printf("User found: %+v", user) // Log the user found

            // Store the user in the context
            ctx := context.WithValue(r.Context(), "user", &user)
            log.Printf("User stored in context: %+v", user) // Verify the user in context

            next.ServeHTTP(w, r.WithContext(ctx))
            return
        }

        http.Error(w, "Unauthorized", http.StatusUnauthorized)
    })
}


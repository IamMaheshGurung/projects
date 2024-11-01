package helper




import(
     "github.com/IamMaheshGurung/restaurant-management/database"
    /* "go.mongodb.org/mongo-driver/v2/mongo"
     "go.mongodb.org/mongo-driver/bson"
     "go.mongodb.org/mongo-driver/bson/primitive"
     "go.mongodb.org/mongo-driver/mongo/options"*/
     "os"
     "github.com/golang-jwt/jwt/v5"
     "net/http"
     "github.com/gin-gonic/gin"
     "time"
     



)





type MyCustomClaims struct {
    Email string         `json:"email"`
    First_name string   `json:"firstname"`
    Last_name string    `json:"lastname"`
    Uid string          `json:"uid"`

    jwt.RegisteredClaims
}


var userCollection  = database.OpenCollection(database.Client, "user")



var SECRET_KEY string = os.Getenv("SECRET_KEY")


func GenerateAllTokens(c *gin.Context, email string, firstname string, lastname string, uid string) (signedToken string, signedRefreshToken string, err error){
    claims := MyCustomClaims{
        Email : email,
        First_name: firstname,
        Last_name: lastname,
        Uid: uid,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }

   refreshedClaims := MyCustomClaims{
        Email:      email,
        First_name: firstname,
        Last_name:  lastname,
        Uid:        uid,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // Refresh token expiration
        },    }
    token :=jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    signedToken, err = token.SignedString(SECRET_KEY)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to signed the token"})
        return 
    }

    refreshedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshedClaims)
    refreshedSignedToken, err  := refreshedToken.SignedString(SECRET_KEY)
     if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to signed the token"})
        return 
    }
    return signedToken, refreshedSignedToken, nil 
}








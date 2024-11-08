package initializers

 import (
     "fmt"
     "github.com/joho/godotenv"


 )

 func LoadEnvVariables() {
     err := godotenv.Load()

     if err != nil {
         fmt.Println("Error while loading .env file")
     }

 }

 

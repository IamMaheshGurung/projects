package main 



import(
    "os"
    "github.com/gin-gonic/gin"
    "github.com/IamMaheshGurung/restaurant-management/database"
    "github.com/IamMaheshGurung/restaurant-managemant/routes"
    "github.com/IamMaheshGurung/restaurant-management/middleware"
    "go.mongodb.org/mongo-driver/mongo"
)


var  foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func main(){
    port := os.Genenv("PORT") 


    if port == "" {
        port = "8080"
    }


    router := gin.New()
    router.Use(gin.Logger())


    router.UserRoutes(router)

    router.Use(middleware.Authentication())

    routes.FoodRoutes(router)
    routes.MenuRoutes(router)
    routes.TableRoutes(router)
    routes.OrderRoutes(router)
    routes.OrderItemRoutes(router)
    routes.InvoiceRoutes(router)

    router.Run(":" + port)



}




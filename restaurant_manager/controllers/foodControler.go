package controller


import(



)



var foodCollection = database.OpenCollection(database.Client, "food")

var validate = validator.New()



func GetFoods()(w http.ResponseWriter, r * http.Request){
    ctx, cancel() := context.WithTimeout(context.Background(), 100 * time.Second)
    defer cancel()

    //Get Pagination parameters 

    recordPerPage, err := strconv.Atoi(r.URL.Query().Get("recordPerPage"))
    if err != nil || recordPerPage < 1 {
        recordPerPage  = 10 
    }


    page, err := strconv.Atoi(r.URL.Query().Get("Page"))
    if err != nil || page < 1 {
        page = 1 
    }

    startIndex := (page -1) * recordPerPage// Calculate the starting index.



//mongoDB pipeline stages

    matchStage := bson.D{{"$match", bson.D{{}}}}

    groupStage := bson.D{{"$group", bson.D{{"_id", "null"}}}, {"total_count", bson.D{{"$sum", 1}}}, {"data", bson.D{{"$push", "$$ROOT"}}}}}
    projectStage := bson.D{
        {
            "$project", bson.D{
                {"_id", 0},
                {"total_count", 1},
                {"food_items", bson.D{{"$slice", []interface{
                }}}

                result, err := foodCollection.Aggregate(ctx, mongo.Pipeline{
                    matchStage, groupStage, projectStage})
                    defer cancel()

                    if err != nil {
                        http.Error(w, "Error", http.StatusInternalServerError)
                        return 
                    }
                    var allFood []bson.M
                    if err = result.All(ctx, &allFoods); err != nil {
                        log.Fatal(err)
                    }
                    w.WriteHeader(http.StatusOK)
                }
            }
        }
    }


package controller


import(



)



var foodCollection = database.OpenCollection(database.Client, "food")

var validate = validator.New()



func GetFoods()(w http.ResponseWriter, r * http.Request){
    ctx, cancel() := context.WithTimeout(context.Background(), 100 * time.Second)
    defer cancel()

    //Get Pagination parameters 

    recordPerPage, err := strconv.Atoi(r.URL.Query().Get("recordPerPage")
    if err != nil || recordPerPage < 1 {
        recordPerPage  = 10 
    }


    page, err := strconv.Atoi(r.URL.Query().Get("Page"))
    if err != nil || page < 1 {
        page = 1 
    }

    startIndex := (page -1) * recordPerPage



    cursor, err := foodCollection.Find(ctx, bson.M{})
    if err != nil {
        http.Error(w, "Error occured while listing food items", http.StatusInternalServerError)
        return 
    }
    defer cursor.Close()

    var allFoods database.*:

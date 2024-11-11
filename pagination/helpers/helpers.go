package helpers



import 
(
    "github.com/IamMaheshGurung/pagination/initializers"
    "math"

)


type PaginationData struct {
    NextPage int
    PreviousPage int
    TwoPageUp int
    TwoPageDown int
    CurrentPage int
    TotalPages int
    Offset int
}



func GetPaginationData(page int, perPage int, model interface{}) PaginationData{
     //Calculating totalPages
    var totalRows int64
    initializers.DB.Model(model).Count(&totalRows)
    totalPages := math.Ceil(float64(totalRows / int64(perPage))) 
    //calculating offset
    offset := (page-1) * perPage

    
    return PaginationData {
            NextPage : page + 1,
            PreviousPage: page -1,
            TwoPageUp: page + 2,
            TwoPageDown: page - 2,
            CurrentPage: page,
            TotalPages : int(totalPages),
            Offset  : offset,
        }
        


}

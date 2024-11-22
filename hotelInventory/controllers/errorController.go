package controllers



import (
    "log"
    "net/http"
    "html/template"
    "os"


)


func RenderError(w http.ResponseWriter, errorMessage string) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("Unable to get current working directory:", err)
	}
	log.Println("Current directory:", dir)

	// Parse the error template from the templates directory
	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		// Log the error and return an internal server error response
		log.Printf("Unable to parse the error template: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	err = tmpl.Execute(w, errorMessage)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Execution error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

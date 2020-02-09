package main

import (
	"./models"
	"./routes"
	"./utils"
	"net/http"
)

func main() {

	models.Init()

	utils.LoadTemplates("templates/*.html")

	r := routes.NewRouter()
	http.Handle("/", r)
	http.ListenAndServe(":8000", nil)
}

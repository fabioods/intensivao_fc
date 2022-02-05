package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Product struct {
	UUID    string  `json:"uuid"`
	Product string  `json:"product"`
	Price   float64 `json:"price,string"`
}

type Products struct {
	Products []Product `json:"products"`
}

var productsURL string

func init() {
	productsURL = os.Getenv("PRODUCT_URL")
}

func loadProducts() []Product {
	response, err := http.Get(productsURL + "/products")
	if err != nil {
		fmt.Println("Error de http", err.Error())
	}
	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)
	var products Products
	json.Unmarshal(data, &products)
	return products.Products
}

func ListProducts(w http.ResponseWriter, r *http.Request) {
	products := loadProducts()
	t := template.Must(template.ParseFiles("templates/catalog.html"))
	t.Execute(w, products)
}

func ShowProducts(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	response, err := http.Get(productsURL + "/product/" + id)
	if err != nil {
		fmt.Println("Error de http", err.Error())
	}
	data, _ := ioutil.ReadAll(response.Body)
	var product Product
	json.Unmarshal(data, &product)
	t := template.Must(template.ParseFiles("templates/view.html"))
	t.Execute(w, product)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", ListProducts)
	r.HandleFunc("/product/{id}", ShowProducts)
	http.ListenAndServe(":8082", r)
}

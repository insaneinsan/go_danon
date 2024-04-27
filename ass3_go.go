package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

const (
	username = "root"
	password = "password"
	hostname = "127.0.0.1"
	port     = 5432
	database = "postgres"
)

var db *sql.DB
var rdClient *redis.Client

type Product struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Price    string `json:"price"`
	Category string `json:"category"`
}

func addProductHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	err1 := r.ParseForm()
	if err1 != nil {
		return
	}
	product := &Product{
		Name:     r.Form.Get("name"),
		Price:    r.Form.Get("price"),
		Category: r.Form.Get("category"),
	}
	err2 := addProductToDb(product)
	if err2 != nil {
		return
	}
}
func addProductToDb(p *Product) error {
	var _, err = db.Exec("INSERT INTO products(name, price, category) VALUES ($1, $2, $3)", p.Name, p.Price, p.Category)
	return err
}

func getProduct(r *http.Request) (*Product, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}
	var product Product
	productId := r.Form.Get("id")
	cacheKey := fmt.Sprintf("product%s", productId)
	cacheResult, err := rdClient.Get(context.Background(), cacheKey).Bytes()
	if err != nil {
		product, err := getProductById(productId)
		if err != nil {
			return nil, err
		}
		jsonProduct, err := json.Marshal(product)
		if err != nil {
			return nil, err
		}
		err = rdClient.Set(context.Background(), cacheKey, jsonProduct, 100*time.Second).Err() // Set expiry to 5 minutes
		if err != nil {
			return nil, err
		}
		return product, nil
	} else {
		err := json.Unmarshal(cacheResult, &product)
		if err != nil {
			return nil, err
		}
		return &product, nil
	}
}

func getProductHandle(w http.ResponseWriter, r *http.Request) {
	product, err := getProduct(r)
	if err != nil {
		return
	}
	if product == nil {
		return
	}
	jsonData, err := json.Marshal(product)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func getProductById(id string) (*Product, error) {
	row := db.QueryRow(fmt.Sprintf("SELECT * FROM products WHERE productId = %s", id))
	product := &Product{}
	err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Category)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func main() {
	DSN := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", username, password, hostname, port, database)
	var err error
	db, err = sql.Open("postgres", DSN)
	if err != nil {
		return
	}

	rdClient = redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	_, err2 := rdClient.Ping(context.Background()).Result()
	if err2 != nil {
		return
	}
	http.HandleFunc("/add-product", addProductHandle)
	http.HandleFunc("/get-product", getProductHandle)
	fmt.Println("Server is running on port 8088...")
	_ = http.ListenAndServe(":8088", nil)
}

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Product struct {
    ProductID      int    `json:"product_id"`
    ProductName    string `json:"product_name"`
    BrandName      string `json:"brand_name"`
    ImagePaths      []string `json:"image_path,omitempty"`
    ProductCreated string `json:"product_created_at"`
    ProductUpdated string `json:"product_updated_at"`
}

type RequestProduct struct {
    ProductID      int      `json:"product_id,omitempty"`
    ProductName    string   `json:"product_name"`
    BrandID        int      `json:"brand_id"`
    ImagePaths     []string `json:"image_paths,omitempty"`
    ProductCreated string   `json:"product_created_at,omitempty"`
    ProductUpdated string   `json:"product_updated_at,omitempty"`
}

func main() {
    err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbName := os.Getenv("DB_NAME")
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "GET":
        products, err := fetchProducts(db)
        if err != nil {
            http.Error(w, "Failed to fetch products.", http.StatusInternalServerError)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(products)
    case "POST":
        var product RequestProduct
        if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        if err := addProduct(db, product); err != nil {
            http.Error(w, "Failed to add product.", http.StatusInternalServerError)
            return
        }
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(product)
    default:
        http.Error(w, "Unsupported request method.", http.StatusMethodNotAllowed)
    }
})
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func fetchProducts(db *sql.DB) ([]Product, error) {
    query := `
        SELECT
            p.id AS product_id,
            p.name AS product_name,
            b.name AS brand_name,
            i.path AS image_path,
            p.created_at AS product_created_at,
            p.updated_at AS product_updated_at
        FROM
            products p
        JOIN
            brands b ON p.brand_id = b.id
        LEFT JOIN
            images i ON p.id = i.product_id
        ORDER BY
            p.id, i.path;
    `
    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var products []Product
    var currentProduct *Product

    for rows.Next() {
        var id int
        var imagePath sql.NullString
        var product Product

        err := rows.Scan(&id, &product.ProductName, &product.BrandName, &imagePath, &product.ProductCreated, &product.ProductUpdated)
        if err != nil {
            return nil, err
        }

        if currentProduct == nil || currentProduct.ProductID != id {
            if currentProduct != nil {
                products = append(products, *currentProduct)
            }
            currentProduct = &product
            currentProduct.ProductID = id
            currentProduct.ImagePaths = []string{}
        }

        if imagePath.Valid {
            currentProduct.ImagePaths = append(currentProduct.ImagePaths, imagePath.String)
        }
    }
    if currentProduct != nil {
        products = append(products, *currentProduct)
    }

    return products, nil
}

func addProduct(db *sql.DB, requestProduct RequestProduct) error {
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		defer tx.Rollback()

		productQuery := `
        INSERT INTO products (name, brand_id) VALUES (?, ?)
    `
		res, err := tx.Exec(productQuery, requestProduct.ProductName, requestProduct.BrandID)
		if err != nil {
			tx.Rollback()
			return err
		}

		productID, err := res.LastInsertId()
		if err != nil {
			tx.Rollback()
			return err
		}

		imageQuery := `
				INSERT INTO images (product_id, path) VALUES (?, ?)
		`
		for _, imagePath := range requestProduct.ImagePaths {
			_, err := tx.Exec(imageQuery, productID, imagePath)
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		return tx.Commit()
}
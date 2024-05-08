package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
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

func HandleProducts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			products, err := FetchProducts(db)
			if err != nil {
				ErrorResponse(w, "Failed to fetch products", http.StatusInternalServerError)
				return
			}
			SuccessResponse(w, products, http.StatusOK)
		case "POST":
			var reqProduct RequestProduct
			if err := json.NewDecoder(r.Body).Decode(&reqProduct); err != nil {
				ErrorResponse(w, "Failed to decode request body", http.StatusBadRequest)
				return
			}
			if err := AddProduct(db, reqProduct); err != nil {
				ErrorResponse(w, "Failed to add product", http.StatusInternalServerError)
				return
			}
			SuccessResponse(w, map[string]string{"message": "Product added"}, http.StatusCreated)
		default:
			ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func FetchProducts(db *sql.DB) ([]Product, error) {
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

func AddProduct(db *sql.DB, requestProduct RequestProduct) error {
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
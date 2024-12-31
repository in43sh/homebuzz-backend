package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/in43sh/homebuzz-backend/database"
)

type Product struct {
	Image        string  `json:"image" binding:"required"`
	ProductTitle string  `json:"product_title" binding:"required"`
	Price        float64 `json:"price" binding:"required"`
	Unit         string  `json:"unit" binding:"required"`
	Rating       int     `json:"rating" binding:"required,gte=1,lte=5"`
}

func AddProduct(ctx *gin.Context) {
	var product Product

	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `INSERT INTO products (image, product_title, price, unit, rating)
              VALUES ($1, $2, $3, $4, $5)`
	_, err := database.Db.Exec(query, product.Image, product.ProductTitle, product.Price, product.Unit, product.Rating)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not insert product into database"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product added successfully!"})
}

func GetProducts(ctx *gin.Context) {
	rows, err := database.Db.Query("SELECT image, product_title, price, unit, rating FROM products")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Couldn't fetch products"})
		return
	}
	defer rows.Close()

	products := []Product{}

	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.Image, &product.ProductTitle, &product.Price, &product.Unit, &product.Rating); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error reading user data"})
			return
		}
		products = append(products, product)
	}

	ctx.JSON(http.StatusOK, gin.H{"products": products})
}

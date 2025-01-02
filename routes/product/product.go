package routes

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/in43sh/homebuzz-backend/database"
)

type Product struct {
	ID           int64   `bun:",pk,autoincrement"`
	Image        string  `bun:"image,notnull" json:"image" binding:"required"`
	ProductTitle string  `bun:"product_title,notnull" json:"product_title" binding:"required"`
	Price        float64 `bun:"price,notnull" json:"price" binding:"required"`
	Unit         string  `bun:"unit,notnull" json:"unit" binding:"required"`
	Rating       int     `bun:"rating,notnull" json:"rating" binding:"required,gte=1,lte=5"`
}

// SuccessResponse for consistent success responses
type SuccessResponse struct {
	Message string `json:"message" example:"Product added successfully!"`
}

// ErrorResponse for consistent error responses
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid input"`
}

// @Summary Add a new product
// @Description Add a new product by providing image, title, price, unit, and rating
// @Tags Products
// @Accept  json
// @Produce  json
// @Param product body Product true "Product information"
// @Success 200 {object} SuccessResponse "Product added successfully!"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 500 {object} ErrorResponse "Could not insert product into database"
// @Router /products [post]
func AddProduct(ctx *gin.Context) {
	var product Product

	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := database.BunDB.NewInsert().Model(&product).Exec(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not insert product into database"})
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse{Message: "Product added successfully!"})
}

// @Summary Get all products
// @Description Retrieve a list of all products in the system
// @Tags Products
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{} "List of products"
// @Failure 500 {object} ErrorResponse "Couldn't fetch products"
// @Router /products [get]
func GetProducts(ctx *gin.Context) {
	var products []Product

	err := database.BunDB.NewSelect().
		Model(&products).
		Scan(context.Background())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: "Couldn't fetch products"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"products": products})
}

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

	ctx.JSON(http.StatusOK, gin.H{"message": "Product added successfully!"})
}

func GetProducts(ctx *gin.Context) {
	var products []Product

	err := database.BunDB.NewSelect().
		Model(&products).
		Scan(context.Background())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Couldn't fetch products"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"products": products})
}

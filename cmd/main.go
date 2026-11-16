package main

import (
	"context"
	"customer-api/internal/handler"
	"customer-api/internal/repository"
	"customer-api/internal/usecase"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := os.Getenv("MYSQL_DSN_CLI")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatal("Gagal koneksi database: ", err)
	}

	fmt.Println("Database connected!")

	repo := repository.NewCustomerAPIRepository(db)
	uc := usecase.NewCustomerAPIUsecase(repo)
	h := handler.NewCustomerAPIHandler(uc)

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "API is running"})
	})
	r.GET("/api/v1/customers", h.GetAllCustomers)
	r.GET("/api/v1/customers/:id", h.GetCustomerByID)
	r.POST("/api/v1/customers", h.AddNewCustomer)
	r.PUT("/api/v1/customers/:id", h.UpdateCustomerByID)
	r.DELETE("/api/v1/customers/:id", h.DeleteCustomerByID)
	r.Run(":8080")

	port := os.Getenv("PORT")
	log.Fatal(r.Run(":" + port))
}

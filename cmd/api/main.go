package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/snck/book-keeper-api/db"
	"github.com/snck/book-keeper-api/handler"
	"github.com/snck/book-keeper-api/repository"
	"github.com/snck/book-keeper-api/service"
)

func main() {

	godotenv.Load()

	err := db.Connect()
	if err != nil {
		fmt.Println("error connecting to DB", err)
		panic("error connecting to DB")
	}
	defer db.Close()

	expenseRepo := repository.NewExpenseRepository(db.DB)
	expenseService := service.NewExpenseService(expenseRepo)
	expenseHandler := handler.NewExpenseHandler(expenseService)

	categoryRepo := repository.NewCategoryRepository(db.DB)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	authRepo := repository.NewAuthenticationRepository(db.DB)
	authService := service.NewAuthenticationService(authRepo)
	authHandler := handler.NewAuthenticationHandler(authService)

	r := gin.Default()

	r.POST("/signup", authHandler.Signup)
	r.POST("/login", authHandler.Login)
	r.POST("/logout", authHandler.Logout)

	r.GET("/expenses", expenseHandler.GetExpenses)
	r.POST("/expenses", expenseHandler.CreateExpense)
	r.PUT("/expenses/:id", expenseHandler.UpdateExpense)
	r.DELETE("/expenses/:id", expenseHandler.DeleteExpense)
	r.GET("/categories", categoryHandler.GetCategories)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	err = r.Run(":8080")
	if err != nil {
		fmt.Println("error starting server: ", err)
	}
}

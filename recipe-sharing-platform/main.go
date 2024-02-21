package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Recipe struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Ingredients  string `json:"ingredients"`
	Instructions string `json:"instructions"`
	Author       string `json:"author"`
}

var db *sql.DB

func main() {

	var err error
	db, err = sql.Open("mysql", "joshua468:170821002@tcp(127.0.0.1:3306)/recipe_db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := gin.Default()

	r.GET("/recipes", getRecipes)
	r.POST("/recipes", createRecipe)
	r.GET("/recipes/:id", getRecipe)
	r.PUT("/recipes/:id", updateRecipe)
	r.DELETE("/recipes/:id", deleteRecipe)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func getRecipes(c *gin.Context) {
	rows, err := db.Query("SELECT * FROM recipes")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var recipes []Recipe
	for rows.Next() {
		var recipe Recipe
		if err := rows.Scan(&recipe.ID, &recipe.Name, &recipe.Description, &recipe.Ingredients, &recipe.Instructions, &recipe.Author); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		recipes = append(recipes, recipe)
	}

	c.JSON(http.StatusOK, recipes)
}

func createRecipe(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO recipes (name, description, ingredients, instructions, author) VALUES (?, ?, ?, ?, ?)",
		recipe.Name, recipe.Description, recipe.Ingredients, recipe.Instructions, recipe.Author)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	recipe.ID = int(id)
	c.JSON(http.StatusCreated, recipe)
}

func getRecipe(c *gin.Context) {
	id := c.Param("id")
	row := db.QueryRow("SELECT * FROM recipes WHERE id = ?", id)

	var recipe Recipe
	if err := row.Scan(&recipe.ID, &recipe.Name, &recipe.Description, &recipe.Ingredients, &recipe.Instructions, &recipe.Author); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	c.JSON(http.StatusOK, recipe)
}

func updateRecipe(c *gin.Context) {
	id := c.Param("id")

	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE recipes SET name=?, description=?, ingredients=?, instructions=?, author=? WHERE id=?",
		recipe.Name, recipe.Description, recipe.Ingredients, recipe.Instructions, recipe.Author, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recipe updated successfully"})
}

func deleteRecipe(c *gin.Context) {
	id := c.Param("id")

	_, err := db.Exec("DELETE FROM recipes WHERE id=?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recipe deleted successfully"})
}

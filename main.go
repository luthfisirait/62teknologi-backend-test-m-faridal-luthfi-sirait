package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/yelp")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	router := gin.Default()

	router.POST("/business", func(c *gin.Context) {
		// Baca parameter dari body request
		var req struct {
			Name        string   `json:"name" binding:"required"`
			Address     string   `json:"address" binding:"required"`
			City        string   `json:"city" binding:"required"`
			State       string   `json:"state" binding:"required"`
			ZipCode     string   `json:"zip_code" binding:"required"`
			Phone       string   `json:"phone" binding:"required"`
			Latitude    float64  `json:"latitude" binding:"required"`
			Longitude   float64  `json:"longitude" binding:"required"`
			Rating      float64  `json:"rating" binding:"required"`
			ReviewCount int      `json:"review_count" binding:"required"`
			Categories  []string `json:"categories" binding:"required"`
			URL         string   `json:"url" binding:"required"`
		}
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "invalid request body",
			})
			return
		}

		// Insert data bisnis baru ke database
		stmt, err := db.Prepare("INSERT INTO businesses (name, address, city, state, zip_code, phone, latitude, longitude, rating, review_count, categories, url) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "failed to prepare SQL statement",
			})
			return
		}
		defer stmt.Close()

		categoriesJSON, _ := json.Marshal(req.Categories)
		_, err = stmt.Exec(req.Name, req.Address, req.City, req.State, req.ZipCode, req.Phone, req.Latitude, req.Longitude, req.Rating, req.ReviewCount, categoriesJSON, req.URL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "failed to insert data to database",
			})
			return
		}

		// Response success
		c.JSON(http.StatusCreated, gin.H{
			"status":  "success",
			"message": "data has been inserted to database",
		})
	})

	router.PUT("/business/:id", func(c *gin.Context) {
		// Baca parameter id dari URL
		businessID := c.Param("id")

		// Baca parameter dari body request
		var req struct {
			Name        string   `json:"name" binding:"required"`
			Address     string   `json:"address" binding:"required"`
			City        string   `json:"city" binding:"required"`
			State       string   `json:"state" binding:"required"`
			ZipCode     string   `json:"zip_code" binding:"required"`
			Phone       string   `json:"phone" binding:"required"`
			Latitude    float64  `json:"latitude" binding:"required"`
			Longitude   float64  `json:"longitude" binding:"required"`
			Rating      float64  `json:"rating" binding:"required"`
			ReviewCount int      `json:"review_count" binding:"required"`
			Categories  []string `json:"categories" binding:"required"`
			URL         string   `json:"url" binding:"required"`
		}
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "invalid request body",
			})
			return
		}

		// Update data bisnis di database
		stmt, err := db.Prepare("UPDATE businesses SET name=?, address=?, city=?, state=?, zip_code=?, phone=?, latitude=?, longitude=?, rating=?, review_count=?, categories=?, url=? WHERE id=?")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "failed to prepare SQL statement",
			})
			return
		}
		defer stmt.Close()

		categoriesJSON, _ := json.Marshal(req.Categories)
		_, err = stmt.Exec(req.Name, req.Address, req.City, req.State, req.ZipCode, req.Phone, req.Latitude, req.Longitude, req.Rating, req.ReviewCount, categoriesJSON, req.URL, businessID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "failed to update data to database",
			})
			return
		}

		// Response success
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "data has been updated to database",
		})
	})
	router.DELETE("/business/:id", func(c *gin.Context) {
		// Ambil id dari path parameter
		id := c.Param("id")

		// Hapus data bisnis dari database
		stmt, err := db.Prepare("DELETE FROM businesses WHERE id = ?")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "failed to prepare SQL statement",
			})
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "failed to delete data from database",
			})
			return
		}

		// Response success
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "data has been deleted from database",
		})
	})
	router.GET("/businesses", func(c *gin.Context) {
		// Baca query parameter
		query := c.Query("pencarian")

		// Cari data bisnis berdasarkan categories, city, atau name yang mengandung query
		rows, err := db.Query("SELECT * FROM businesses WHERE categories LIKE ? OR city LIKE ? OR name LIKE ? ORDER BY review_count DESC", "%"+query+"%", "%"+query+"%", "%"+query+"%")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "failed to query data from database",
			})
			return
		}
		defer rows.Close()

		// Parse data dari hasil query ke slice business
		businesses := make([]map[string]interface{}, 0)
		for rows.Next() {
			var id int
			var name string
			var address string
			var city string
			var state string
			var zip_code string
			var phone string
			var latitude float64
			var longitude float64
			var rating float64
			var review_count int
			var categories string
			var url string
			err = rows.Scan(&id, &name, &address, &city, &state, &zip_code, &phone, &latitude, &longitude, &rating, &review_count, &categories, &url)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"message": "failed to parse data from database",
				})
				return
			}

			business := make(map[string]interface{})
			business["id"] = id
			business["name"] = name
			business["address"] = address
			business["city"] = city
			business["state"] = state
			business["zip_code"] = zip_code
			business["phone"] = phone
			business["latitude"] = latitude
			business["longitude"] = longitude
			business["rating"] = rating
			business["review_count"] = review_count
			business["categories"] = categories
			business["url"] = url
			businesses = append(businesses, business)
		}

		// Response data bisnis
		c.JSON(http.StatusOK, gin.H{
			"status":     "success",
			"businesses": businesses,
		})
	})
	err = router.Run(":8080")

	if err != nil {
		fmt.Println(err)
	}
}

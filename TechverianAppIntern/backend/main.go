package main

import (
    "log"
    "net/http"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
)

type Item struct {
    ID    uint    `json:"id" gorm:"primary_key"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}

type Response struct {
    TotalCost   float64 `json:"total_cost"`
    AverageCost float64 `json:"average_cost"`
}

var db *gorm.DB

func init() {
    var err error
    dsn := "host=localhost user=EnterYourUserName password=PleaseInsertYourDbPass dbname=YourDatabase`       port=5432 sslmode=disable"
    db, err = gorm.Open("postgres", dsn)
    if err != nil {
        log.Fatalf("failed to connect to database: %v", err)
    }
    db.AutoMigrate(&Item{})
}

func main() {
    router := gin.Default()
    router.Use(cors.Default())

    router.POST("/calculate", func(c *gin.Context) {
        var items []Item
        if err := c.BindJSON(&items); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        for _, item := range items {
            db.Create(&item)
        }

        var totalCost float64
        db.Model(&Item{}).Select("SUM(price)").Scan(&totalCost)

        var count int
        db.Model(&Item{}).Count(&count)

        averageCost := 0.0
        if count > 0 {
            averageCost = totalCost / float64(count)
        }

        response := Response{
            TotalCost:   totalCost,
            AverageCost: averageCost,
        }

        c.JSON(http.StatusOK, response)
    })

    router.Run(":8080")
}

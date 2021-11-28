package route_handlers

import (
	"fmt"
	"linkbook/cli/db"
	"linkbook/cli/db/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetAllCategory(c *gin.Context) {

	collection := db.DbClient.Database(db.Database).Collection("categorys")
	cur, currErr := collection.Find(db.DbCtx, bson.D{})
	if currErr != nil {
		panic(currErr)
	}
	var category []models.Category
	if err := cur.All(db.DbCtx, &category); err != nil {
		panic(err)
	}
	defer cur.Close(db.DbCtx)
	c.JSON(200, gin.H{
		"data": category,
	})
}

func AddCategory(c *gin.Context) {
	collection := db.DbClient.Database(db.Database).Collection("categorys")
	var category models.Category
	c.BindJSON(&category)
	fmt.Println(category)
	res, err := collection.InsertOne(db.DbCtx, bson.D{
		{Key: "name", Value: category.Name},
	})
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{"Data": res})

}

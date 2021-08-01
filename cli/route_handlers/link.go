package route_handlers

import (
	"linkbook/cli/db"
	"linkbook/cli/db/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetAllLink(c *gin.Context) {

	collection := db.DbClient.Database(db.Database).Collection("links")
	cur, currErr := collection.Find(db.DbCtx, bson.D{})
	if currErr != nil {
		panic(currErr)
	}

	var links []models.Link
	if err := cur.All(db.DbCtx, &links); err != nil {
		panic(err)
	}

	c.JSON(200, gin.H{
		"data": links,
	})
}

func PostSingleLink(c *gin.Context) {
	var linkData models.Link
	c.BindJSON(&linkData)
	collection := db.DbClient.Database(db.Database).Collection("links")

	res, err := collection.InsertOne(db.DbCtx, bson.D{
		{Key: "link", Value: linkData.Link},
	})
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{"Data": res})
}

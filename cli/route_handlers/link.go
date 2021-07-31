package route_handlers

import (
	"linkbook/cli/db"
	"linkbook/cli/db/models"

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
		c.JSON(200, gin.H{
			"error": err,
		})
	}

	c.JSON(200, gin.H{
		"data": links,
	})
}

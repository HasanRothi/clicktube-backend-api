package route_handlers

import (
	"fmt"
	"linkbook/cli/db"
	"linkbook/cli/db/models"
	"log"
	"net/http"
	"strconv"
	"time"

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
	// fmt.Println(reflect.TypeOf(links))
	defer cur.Close(db.DbCtx)
	c.JSON(200, gin.H{
		"data": links,
	})
}

func GetSingleLink(c *gin.Context) {

	// linkID, err := primitive.ObjectIDFromHex(c.Param("id"))
	// if err != nil {
	// 	panic(err)
	// }
	collection := db.DbClient.Database(db.Database).Collection("links")
	filterCursor, err := collection.Find(db.DbCtx, bson.M{"shortLink": c.Param("id")})
	if err != nil {
		log.Fatal(err)
	}
	var links []models.Link
	if err = filterCursor.All(db.DbCtx, &links); err != nil {
		log.Fatal(err)
	}
	result, err := collection.UpdateOne(
		db.DbCtx,
		bson.M{"_id": links[0].ID},
		bson.D{
			{"$set", bson.D{{"views", links[0].Views + 1}}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
	c.Redirect(302, links[0].Link)
}

func PostSingleLink(c *gin.Context) {
	var linkData models.Link
	c.BindJSON(&linkData)
	collection := db.DbClient.Database(db.Database).Collection("links")
	cur, currErr := collection.Find(db.DbCtx, bson.D{})
	if currErr != nil {
		panic(currErr)
	}

	var links []models.Link
	if err := cur.All(db.DbCtx, &links); err != nil {
		panic(err)
	}
	shortLink := genarateSortLink(len(links))
	res, err := collection.InsertOne(db.DbCtx, bson.D{
		{Key: "link", Value: linkData.Link},
		{Key: "shortLink", Value: shortLink},
		{Key: "date", Value: time.Now()},
	})
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	// defer collection.Close()
	c.JSON(http.StatusOK, gin.H{"Data": res})
}

func genarateSortLink(next int) string {
	currentDate := time.Now()
	dateString := currentDate.String()
	dateSlice := dateString[0:10]
	CAMPUS_CODE := "UITS"
	DEPT_CODE := "IT"
	nextLink := next + 1
	shortUrl := dateSlice + "-" + CAMPUS_CODE + "-" + DEPT_CODE + "-" + strconv.Itoa(nextLink)
	return shortUrl
}

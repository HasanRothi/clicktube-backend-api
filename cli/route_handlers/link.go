package route_handlers

import (
	"fmt"
	"linkbook/cli/db"
	"linkbook/cli/db/models"
	"linkbook/cli/services"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllLink(c *gin.Context) {

	collection := db.DbClient.Database(db.Database).Collection("links")
	//pagination 1
	// 	opts.SetSkip(0)
	// opts.SetLimit(2)
	//pagination 2
	// 	skip := int64(0)
	// limit := int64(10)
	// opts := options.FindOptions{
	//   Skip: skip,
	//   Limit: limit
	// }
	//v1
	// opts := options.Find()
	// opts.SetSort(bson.D{{"date", -1}})
	// cur, currErr := collection.Find(db.DbCtx, bson.D{}, opts)
	// if currErr != nil {
	// 	panic(currErr)
	// }

	// var links []models.Link
	// if err := cur.All(db.DbCtx, &links); err != nil {
	// 	panic(err)
	// }
	// fmt.Println(reflect.TypeOf(links))
	//v2
	lookupStage := bson.D{{"$lookup", bson.D{{"from", "users"}, {"localField", "author"}, {"foreignField", "_id"}, {"as", "author"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$author"}, {"preserveNullAndEmptyArrays", false}}}}

	linkWithAuthorCur, err := collection.Aggregate(db.DbCtx, mongo.Pipeline{lookupStage, unwindStage})
	if err != nil {
		panic(err)
	}
	var links []bson.M
	if err = linkWithAuthorCur.All(db.DbCtx, &links); err != nil {
		panic(err)
	}
	defer linkWithAuthorCur.Close(db.DbCtx)
	c.JSON(200, gin.H{
		"data": links,
	})
}

func GetSingleLink(c *gin.Context) {

	collection := db.DbClient.Database(db.Database).Collection("links")
	filterCursor, err := collection.Find(db.DbCtx, bson.M{"shortLink": c.Param("id"), "published": true})
	if err != nil {
		log.Fatal(err)
	}
	var links []models.Link
	if err = filterCursor.All(db.DbCtx, &links); err != nil {
		log.Fatal(err)
	}
	if len(links) == 0 {
		panic("Url Not Found")
	} else {
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

}

func GetPopularLinks(c *gin.Context) {

	collection := db.DbClient.Database(db.Database).Collection("links")
	opts := options.Find()
	opts.SetSort(bson.D{{"views", -1}})
	sortCursor, err := collection.Find(db.DbCtx, bson.D{}, opts)
	if err != nil {
		log.Fatal(err)
	}
	var links []bson.M
	if err = sortCursor.All(db.DbCtx, &links); err != nil {
		log.Fatal(err)
	}
	c.JSON(200, gin.H{
		"data": links,
	})
}
func GetPendingLinks(c *gin.Context) {

	collection := db.DbClient.Database(db.Database).Collection("links")
	sortCursor, err := collection.Find(db.DbCtx, bson.M{"published": false})
	if err != nil {
		log.Fatal(err)
	}
	var links []bson.M
	if err = sortCursor.All(db.DbCtx, &links); err != nil {
		log.Fatal(err)
	}
	c.JSON(200, gin.H{
		"data": links,
	})
}

func PostSingleLink(c *gin.Context) {
	var linkData models.Link
	c.BindJSON(&linkData)
	collection := db.DbClient.Database(db.Database).Collection("links")
	filterCursor, err := collection.Find(db.DbCtx, bson.M{"link": linkData.Link})
	if err != nil {
		log.Fatal(err)

	}
	var links []models.Link
	if err = filterCursor.All(db.DbCtx, &links); err != nil {
		log.Fatal(err)
	}
	if len(links) > 0 {
		c.JSON(400, gin.H{"Data": "Link Already Exist"})
	} else {
		author := LoadSingleUser(linkData.Author)
		urlKey := services.UrlKey() + "-" + author[0].University + "-" + author[0].Dept
		cur, currErr := collection.Find(db.DbCtx, bson.M{"urlKey": urlKey})
		if currErr != nil {
			panic(currErr)
		}

		if err := cur.All(db.DbCtx, &links); err != nil {
			panic(err)
		}
		shortLink := services.GenarateSortLink(len(links), author[0].University, author[0].Dept)
		res, err := collection.InsertOne(db.DbCtx, bson.D{
			{Key: "link", Value: linkData.Link},
			{Key: "shortLink", Value: shortLink},
			{Key: "date", Value: time.Now()},
			{Key: "published", Value: linkData.Published},
			{Key: "urlKey", Value: urlKey},
			{Key: "author", Value: author[0].ID},
		})
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
		c.JSON(http.StatusOK, gin.H{"Data": res})
	}
}

func PublishedSingleLink(c *gin.Context) {
	var linkData models.Link
	c.BindJSON(&linkData)
	collection := db.DbClient.Database(db.Database).Collection("links")
	result, err := collection.UpdateOne(
		db.DbCtx,
		bson.M{"_id": linkData.ID},
		bson.D{
			{"$set", bson.D{{"published", linkData.Published}}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("UpdateOne() result:", result)
	// fmt.Println("UpdateOne() result TYPE:", reflect.TypeOf(result))
	// fmt.Println("UpdateOne() result MatchedCount:", result.MatchedCount)
	// fmt.Println("UpdateOne() result ModifiedCount:", result.ModifiedCount)
	// fmt.Println("UpdateOne() result UpsertedCount:", result.UpsertedCount)
	// fmt.Println("UpdateOne() result UpsertedID:", result.UpsertedID)
	if linkData.Published == true && result.ModifiedCount == 1 {
		filterCursor, err := collection.Find(db.DbCtx, bson.M{"_id": linkData.ID})
		if err != nil {
			log.Fatal(err)
			panic("Link Not Found")
		}
		var links []models.Link
		if err = filterCursor.All(db.DbCtx, &links); err != nil {
			log.Fatal(err)
		}
		author := LoadSingleUser(linkData.Author)
		fmt.Println(author)
		services.SendMail(author[0].Name, author[0].Gmail, links[0].Link, links[0].ShortLink, time.Now())
	}
	c.JSON(200, gin.H{
		"data": "Link Published",
	})
}

package route_handlers

import (
	// "encoding/json"
	"fmt"
	"linkbook/cli/db"
	"linkbook/cli/db/models"
	"linkbook/cli/helpers"
	"linkbook/cli/services"
	"log"
	"net/http"

	// "net/url"
	// "os"
	// "strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
	// lookupStage := bson.D{{"$lookup", bson.D{{"from", "users"}, {"localField", "author"}, {"foreignField", "_id"}, {"as", "author"}}}}
	// unwindStage := bson.D{{"$unwind", bson.D{{"path", "$author"}, {"preserveNullAndEmptyArrays", false}}}}
	// pipeline := bson.D{{
	// 	"$project",
	// 	bson.D{
	// 		{"pubslihed", true},
	// 	},
	// }}

	// linkWithAuthorCur, err := collection.Aggregate(db.DbCtx, mongo.Pipeline{pipeline, lookupStage})
	// if err != nil {
	// 	panic(err)
	// }
	// var links []bson.M
	// if err = linkWithAuthorCur.All(db.DbCtx, &links); err != nil {
	// 	panic(err)
	// }
	// defer linkWithAuthorCur.Close(db.DbCtx)
	// c.JSON(200, gin.H{
	// 	"data": links,
	// })
	sortCursor, err := collection.Find(db.DbCtx, bson.M{"published": true})
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

func GetSingleLink(c *gin.Context) {
	val, err := db.DbRedisClient.Get(c.Param("id")).Result()
	if err != nil {
		fmt.Println(err)
	}
	if len(val) > 0 {
		c.Redirect(302, val)
	} else {
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
			err = db.DbRedisClient.Set(c.Param("id"), links[0].Link, 0).Err()
			// accountSid := os.Getenv("TWILO_SID")
			// authToken := os.Getenv("TWILO_TOKEN")
			// a := "https://api.twilio.com/2010-04-01/Accounts/"
			// d := "/Messages.json"
			// urlStr := a + accountSid + d
			// msgData := url.Values{}
			// msgData.Set("To", os.Getenv("TO_PHONE_NUMBER"))
			// msgData.Set("From", os.Getenv("TWILIO_PHONE_NUMBER"))
			// msgData.Set("Body", "Hi from Golang,Rothi")
			// msgDataReader := *strings.NewReader(msgData.Encode())
			// client := &http.Client{}
			// req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
			// req.SetBasicAuth(accountSid, authToken)
			// req.Header.Add("Accept", "application/json")
			// req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			// resp, _ := client.Do(req)
			// if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			// 	var data map[string]interface{}
			// 	decoder := json.NewDecoder(resp.Body)
			// 	err := decoder.Decode(&data)
			// 	if err == nil {
			// 		fmt.Println(data["sid"])
			// 	}
			// } else {
			// 	fmt.Println(resp.Status)
			// }
			c.Redirect(302, links[0].Link)
		}
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
	error := helpers.SchemaValidator(linkData)
	if len(error) > 0 {
		panic("Link Validation Falied")
	}
	urlValid := helpers.UrlValidator(linkData.Link)
	if urlValid == false {
		panic("Link is not valid")
	}
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
		// fmt.Println(author)
		urlKey := services.UrlKey() + "-" + author[0].University + "-" + author[0].Dept
		cur, currErr := collection.Find(db.DbCtx, bson.M{"urlKey": urlKey})
		if currErr != nil {
			panic(currErr)
		}

		if err := cur.All(db.DbCtx, &links); err != nil {
			panic(err)
		}
		shortLink := services.GenarateSortLink(len(links), author[0].University, author[0].Dept)
		// fmt.Println(linkData)
		res, err := collection.InsertOne(db.DbCtx, bson.D{
			{Key: "link", Value: linkData.Link},
			{Key: "title", Value: linkData.Title},
			{Key: "description", Value: linkData.Description},
			{Key: "shortLink", Value: shortLink},
			{Key: "date", Value: time.Now()},
			{Key: "published", Value: linkData.Published},
			{Key: "urlKey", Value: urlKey},
			{Key: "author", Value: author[0].ID},
			{Key: "category", Value: linkData.Category},
		})
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
		collection := db.DbClient.Database(db.Database).Collection("users")
		_, err = collection.UpdateOne(
			db.DbCtx,
			bson.M{"_id": linkData.Author},
			bson.M{"$push": bson.M{"links": res.InsertedID}},
		)
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(result)
		// authorData := LoadSingleUser(linkData.Author)
		collection = db.DbClient.Database(db.Database).Collection("links")
		filterLink, err := collection.Find(db.DbCtx, bson.M{"_id": res.InsertedID})
		if err != nil {
			log.Fatal(err)

		}
		var link []models.Link
		if err = filterLink.All(db.DbCtx, &link); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"Data": link})
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

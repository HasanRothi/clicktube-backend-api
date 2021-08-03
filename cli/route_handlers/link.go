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
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllLink(c *gin.Context) {

	collection := db.DbClient.Database(db.Database).Collection("links")
	opts := options.Find()
	opts.SetSort(bson.D{{"date", -1}})
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
	cur, currErr := collection.Find(db.DbCtx, bson.D{}, opts)
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
	cur, currErr := collection.Find(db.DbCtx, bson.M{"date": time.Now()})
	if currErr != nil {
		panic(currErr)
	}

	var links []models.Link
	if err := cur.All(db.DbCtx, &links); err != nil {
		panic(err)
	}
	fmt.Println(len(links))
	// shortLink := genarateSortLink(len(links))
	// res, err := collection.InsertOne(db.DbCtx, bson.D{
	// 	{Key: "link", Value: linkData.Link},
	// 	{Key: "shortLink", Value: shortLink},
	// 	{Key: "date", Value: time.Now()},
	// 	{Key: "published", Value: linkData.Published},
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// 	panic(err)
	// }
	// defer collection.Close()
	c.JSON(http.StatusOK, gin.H{"Data": "lol"})
}

func PublishedSingleLink(c *gin.Context) {
	var linkData models.Link
	c.BindJSON(&linkData)
	collection := db.DbClient.Database(db.Database).Collection("links")
	result, err := collection.UpdateOne(
		db.DbCtx,
		bson.M{"_id": linkData.ID},
		bson.D{
			{"$set", bson.D{{"published", true}}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
	c.JSON(200, gin.H{
		"data": "Link Published",
	})
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

// import gomail "gopkg.in/mail.v2"

// m := gomail.NewMessage()
// m.SetHeader("From", "hasanrothi99@gmail.com")
// m.SetHeader("To", "hasanrothi@gmail.com")
// m.SetAddressHeader("Cc", "dan@example.com", "Dan")
// m.SetHeader("Subject", "Hello!")
// m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")
// m.Attach("/home/Alex/lolcat.jpg")

// d := gomail.NewDialer("smtp.example.com", 587, "hasanrothi99@gmail.com", "kemonaso99")
// This is only needed when SSL/TLS certificate is not valid on server.
// In production this should be set to false.
//   d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

// // Send the email to Bob, Cora and Dan.
// if err := d.DialAndSend(m); err != nil {
// 	panic(err)
// }

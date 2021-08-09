package route_handlers

import (
	"linkbook/cli/db"
	"linkbook/cli/db/models"
	"linkbook/cli/helpers"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllUser(c *gin.Context) {

	collection := db.DbClient.Database(db.Database).Collection("users")
	opts := options.Find()
	opts.SetProjection(bson.M{"name": 1, "gmail": 1, "university": 1, "campusId": 1, "dept": 1, "links": 1, "role": 1})
	cur, currErr := collection.Find(db.DbCtx, bson.D{}, opts)
	if currErr != nil {
		panic(currErr)
	}
	// cur, currErr := collection.Find(db.DbCtx, bson.D{})
	// if currErr != nil {
	// 	panic(currErr)
	// }

	var users []models.User
	if err := cur.All(db.DbCtx, &users); err != nil {
		panic(err)
	}
	defer cur.Close(db.DbCtx)
	c.JSON(200, gin.H{
		"data": users,
	})
}

func GetSingleUser(c *gin.Context) {
	userID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		panic(err)
	}
	user := LoadSingleUser(userID)
	c.JSON(200, gin.H{
		"data": user,
	})

}
func LoadSingleUser(userID primitive.ObjectID) []models.User {
	collection := db.DbClient.Database(db.Database).Collection("users")
	cur, err := collection.Find(db.DbCtx, bson.M{"_id": userID})
	if err != nil {
		log.Fatal(err)
		panic("Author | User Not Found")
	}
	var user []models.User
	if err = cur.All(db.DbCtx, &user); err != nil {
		log.Fatal(err)
	}
	return user
}
func GetSingleUserLinks(c *gin.Context) {
	collection := db.DbClient.Database(db.Database).Collection("users")
	lookupStage := bson.D{{"$lookup", bson.D{{"from", "links"}, {"localField", "links"}, {"foreignField", "_id"}, {"as", "links"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$links"}, {"preserveNullAndEmptyArrays", false}}}}

	linkWithAuthorCur, err := collection.Aggregate(db.DbCtx, mongo.Pipeline{lookupStage, unwindStage})
	if err != nil {
		panic(err)
	}
	var user []bson.M
	if err = linkWithAuthorCur.All(db.DbCtx, &user); err != nil {
		panic(err)
	}
	defer linkWithAuthorCur.Close(db.DbCtx)
	c.JSON(200, gin.H{"Data": user})
}

func PostSingleUser(c *gin.Context) {
	var userData models.User
	c.BindJSON(&userData)
	error := helpers.SchemaValidator(userData)
	if len(error) > 0 {
		panic("User Validation Falied")
	}
	collection := db.DbClient.Database(db.Database).Collection("users")
	cur, err := collection.Find(db.DbCtx, bson.M{"gmail": userData.Gmail})
	if err != nil {
		log.Fatal(err)
		panic("Try Again")
	}
	var user []models.User
	if err = cur.All(db.DbCtx, &user); err != nil {
		log.Fatal(err)
	}
	if len(user) > 0 {
		c.JSON(400, gin.H{"Data": "User Already Exist"})
	} else {
		hashPass, _ := helpers.HashPassword(userData.Password)
		res, err := collection.InsertOne(db.DbCtx, bson.D{
			{Key: "gmail", Value: userData.Gmail},
			{Key: "name", Value: userData.Name},
			{Key: "password", Value: hashPass},
			{Key: "university", Value: userData.University},
			{Key: "campusId", Value: userData.CampusID},
			{Key: "dept", Value: userData.Dept},
		})
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
		c.JSON(http.StatusOK, gin.H{"Data": res})
	}

}

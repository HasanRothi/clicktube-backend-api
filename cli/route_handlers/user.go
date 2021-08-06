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
)

func GetAllUser(c *gin.Context) {

	collection := db.DbClient.Database(db.Database).Collection("users")
	cur, currErr := collection.Find(db.DbCtx, bson.D{})
	if currErr != nil {
		panic(currErr)
	}

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
		panic("Author Not Found")
	}
	var user []models.User
	if err = cur.All(db.DbCtx, &user); err != nil {
		log.Fatal(err)
	}
	return user
}

func PostSingleUser(c *gin.Context) {
	var userData models.User
	c.BindJSON(&userData)
	collection := db.DbClient.Database(db.Database).Collection("users")
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

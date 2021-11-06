package route_handlers

import (
	"linkbook/cli/db"
	"linkbook/cli/db/models"
	"linkbook/cli/helpers"
	"linkbook/cli/services"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func Login(c *gin.Context) {
	var userData models.User
	c.BindJSON(&userData)
	collection := db.DbClient.Database(db.Database).Collection("users")
	cur, err := collection.Find(db.DbCtx, bson.M{"gmail": userData.Gmail})
	if err != nil {
		log.Fatal(err)
	}
	var user []models.User
	if err = cur.All(db.DbCtx, &user); err != nil {
		log.Fatal(err)
	}
	if len(user) == 0 {
		c.JSON(404, gin.H{
			"message": "User Not Found",
		})
	} else {
		match := helpers.CheckPasswordHash(userData.Password, user[0].Password)
		if match == true {
			token := services.GenarateJwtToken(user[0].Gmail)
			// c.Set("isAuth", true)
			c.JSON(200, gin.H{
				"message":  "Logged in",
				"id":       user[0].ID,
				"gmail":    user[0].Gmail,
				"CampusId": user[0].CampusID,
				"token":    token,
			})
		} else {
			c.JSON(500, gin.H{
				"message": "try Again",
			})
		}

	}
}

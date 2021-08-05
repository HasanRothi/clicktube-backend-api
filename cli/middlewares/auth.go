package middlewares

import (
	"fmt"
	"linkbook/cli/db"
	"linkbook/cli/db/models"
	"linkbook/cli/services"
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func UserAuth(c *gin.Context) {
	const BEARER_SCHEMA = "Bearer "
	authHeader := c.GetHeader("Authorization")
	tokenString := authHeader[len(BEARER_SCHEMA):]
	token, err := services.VerifyJwtToken(tokenString)

	if token.Valid {
		// claims := token.Claims.(jwt.MapClaims)
		// fmt.Println("Claims[Gmail]: ", claims["gmail"])
		c.Next()
	} else {
		fmt.Println(err)
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
		c.Abort()
	}
}

func AdminAuth(c *gin.Context) {
	const BEARER_SCHEMA = "Bearer "
	authHeader := c.GetHeader("Authorization")
	tokenString := authHeader[len(BEARER_SCHEMA):]
	token, err := services.VerifyJwtToken(tokenString)

	if token.Valid {
		claims := token.Claims.(jwt.MapClaims)
		collection := db.DbClient.Database(db.Database).Collection("users")
		cur, err := collection.Find(db.DbCtx, bson.M{"gmail": claims["gmail"]})
		if err != nil {
			log.Fatal(err)
		}
		var user []models.User
		if err = cur.All(db.DbCtx, &user); err != nil {
			log.Fatal(err)
		}
		if len(user) == 0 {
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
		} else {
			if user[0].Role == "Admin" || user[0].Role == "SuperAdmin" {
				c.Next()
			} else {
				c.JSON(401, gin.H{
					"message": "Unauthenticated",
				})
				c.Abort()
			}
		}
		c.Next()
	} else {
		fmt.Println(err)
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
		c.Abort()
	}
}
func SuperAdminAuth(c *gin.Context) {
	const BEARER_SCHEMA = "Bearer "
	authHeader := c.GetHeader("Authorization")
	tokenString := authHeader[len(BEARER_SCHEMA):]
	token, err := services.VerifyJwtToken(tokenString)

	if token.Valid {
		claims := token.Claims.(jwt.MapClaims)
		collection := db.DbClient.Database(db.Database).Collection("users")
		cur, err := collection.Find(db.DbCtx, bson.M{"gmail": claims["gmail"]})
		if err != nil {
			log.Fatal(err)
		}
		var user []models.User
		if err = cur.All(db.DbCtx, &user); err != nil {
			log.Fatal(err)
		}
		if len(user) == 0 {
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
		} else {
			if user[0].Role == "SuperAdmin" {
				c.Next()
			} else {
				c.JSON(401, gin.H{
					"message": "Unauthenticated",
				})
				c.Abort()
			}
		}
		c.Next()
	} else {
		fmt.Println(err)
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
	}
}

// func SomeHandler() gin.HandlerFunc {
// 	fn := func(c *gin.Context) {
// 		var userData models.User
// 		c.BindJSON(&userData)
// 		fmt.Println(userData)
// 	}

// 	return gin.HandlerFunc(fn)
// }

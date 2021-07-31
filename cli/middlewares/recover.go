package middlewares

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func Recover(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			//Print error stack information
			log.Printf("panic: %v\n", r)
			debug.PrintStack()
			//Package general json return
			//c.JSON(http.StatusOK, Result.Fail(errorToString(r)))
			//Result.Fail is not the focus of this example, so use the following code instead
			c.JSON(http.StatusOK, gin.H{
				"code": "1",
				"msg":  errorToString(r),
				"data": nil,
			})
			// Terminate subsequent interface calls, if not added, after recovering to an exception, the subsequent code in the interface will continue to execute
			c.Abort()
		}
	}()
	//After loading defer recover, continue with subsequent interface calls
	c.Next()
}

// recover error, turn to string
func errorToString(r interface{}) string {
	switch v := r.(type) {
	case error:
		return v.Error()
	default:
		return r.(string)
	}
}

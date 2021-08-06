package middlewares

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
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
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": "500",
				"msg":  errorToString(r),
				"data": nil,
			})

			//senrty
			if hub := sentrygin.GetHubFromContext(c); hub != nil {
				hub.WithScope(func(scope *sentry.Scope) {
					scope.SetExtra("unwantedQuery1", "someQueryDataMaybe1")
					hub.CaptureMessage(errorToString(r))
				})
			}
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

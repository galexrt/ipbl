package routes

import (
	"net/http"
	"strconv"

	"github.com/galexrt/ipbl/pkg/db"
	"github.com/galexrt/ipbl/pkg/models"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int16
	Error   error
	Message string
	Result  interface{}
}

// Register registers routes
func Register(e *gin.Engine) {
	e.GET("/", func(c *gin.Context) {
		lists := []models.List{}
		if err := db.DBCon.Select(&lists, "SELECT ID, Name, Comment, Created, Updated FROM ipbl.List;"); err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Code:  http.StatusInternalServerError,
				Error: err,
			})
			c.Error(err)
			return
		}

		c.IndentedJSON(http.StatusOK, Response{
			Code:    http.StatusOK,
			Message: "Welcome to galexrt/ipbl server! The result key contains the available IP lists.",
			Result:  lists,
		})
	})
	e.GET("/ipbl/:ListID", func(c *gin.Context) {
		var err error
		var listID int
		paramListID := c.Param("ListID")
		if listID, err = strconv.Atoi(paramListID); err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Code:  http.StatusInternalServerError,
				Error: err,
			})
			c.Error(err)
			return
		}

		ips := []models.IP{}
		if err := db.DBCon.Select(&ips, "SELECT ID, ListID, INET6_NTOA(Address) AS Address, Network, Comment, Created, Updated FROM ipbl.IPAddress WHERE ListID = ?;", listID); err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Code:  http.StatusInternalServerError,
				Error: err,
			})
			c.Error(err)
			return
		}

		c.IndentedJSON(http.StatusOK, Response{
			Code:   http.StatusOK,
			Result: ips,
		})
	})
}

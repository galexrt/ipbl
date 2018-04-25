package routes

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/galexrt/ipbl/pkg/db"
	"github.com/galexrt/ipbl/pkg/models"
	"github.com/gin-gonic/gin"
)

func registerList(r *gin.Engine) {
	r.GET("/ipbl", ListLists)
	r.POST("/ipbl", CreateList)
	r.DELETE("/ipbl", DeleteList)
}

func ListLists(c *gin.Context) {
	outputRenderer := getOutputRenderer(c)

	lists := []models.List{}
	if err := db.DBCon.Select(&lists, "SELECT ID, Name, Comment, Created, Updated FROM ipbl.List;"); err != nil {
		outputRenderer(http.StatusInternalServerError, Response{
			Code:  http.StatusInternalServerError,
			Error: err,
		})
		c.Error(err)
		return
	}
	outputRenderer(http.StatusOK, Response{
		Code:   http.StatusOK,
		Result: lists,
	})
}

func CreateList(c *gin.Context) {
	outputRenderer := getOutputRenderer(c)

	list := models.List{}
	now := time.Now()
	list.Created = now
	list.Updated = now
	var err error
	var result sql.Result
	if result, err = db.DBCon.NamedExec("INSERT INTO ipbl.List (Name, Comment) VALUES (:Name, :Comment);", &list); err != nil {
		outputRenderer(http.StatusInternalServerError, Response{
			Code:   http.StatusInternalServerError,
			Error:  err,
			Result: list,
		})
		c.Error(err)
		return
	}
	list.ID, err = result.LastInsertId()
	if err != nil {
		outputRenderer(http.StatusInternalServerError, Response{
			Code:  http.StatusInternalServerError,
			Error: err,
		})
		c.Error(err)
		return
	}

	outputRenderer(http.StatusOK, Response{
		Code:   http.StatusOK,
		Result: list,
	})
}

func DeleteList(c *gin.Context) {
	outputRenderer := getOutputRenderer(c)

	outputRenderer(http.StatusNotImplemented, Response{
		Code:   http.StatusOK,
		Result: nil,
	})

	// TODO
}

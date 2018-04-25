package routes

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
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
	var err error
	var listID int
	paramListID := c.Param("ListID")
	if listID, err = strconv.Atoi(paramListID); err != nil {
		err = fmt.Errorf("no, empty or invalid ListID given")
		outputRenderer(http.StatusBadRequest, Response{
			Code:  http.StatusBadRequest,
			Error: err,
		})
		c.Error(err)
		return
	}

	var result sql.Result
	if result, err = db.DBCon.Exec("DELETE FROM ipbl.List WHERE ID = ?;", listID); err != nil {
		outputRenderer(http.StatusInternalServerError, Response{
			Code:  http.StatusInternalServerError,
			Error: err,
		})
		c.Error(err)
		return
	}

	affected, err := result.RowsAffected()

	outputRenderer(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: fmt.Sprintf("%d row(s) affected", affected),
		Result: models.List{
			ID: int64(listID),
		},
		Error: err,
	})
}

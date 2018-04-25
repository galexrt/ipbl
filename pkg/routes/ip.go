package routes

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/galexrt/ipbl/pkg/db"
	"github.com/galexrt/ipbl/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func registerIP(r *gin.Engine) {
	r.GET("/ipbl/:ListID", ListIPsFromList)
	r.POST("/ipbl/:ListID", InsertIPIntoList)
	r.DELETE("/ipbl/:ListID", DeleteIPFromList)
}

func ListIPsFromList(c *gin.Context) {
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

	list := []models.List{}
	if err := db.DBCon.Select(&list, "SELECT ID, Name, Comment, Created, Updated FROM ipbl.List WHERE ID = ? LIMIT 1;", listID); err != nil {
		outputRenderer(http.StatusInternalServerError, Response{
			Code:  http.StatusInternalServerError,
			Error: err,
		})
		c.Error(err)
		return
	}
	if len(list) == 0 || len(list) > 1 {
		outputRenderer(http.StatusInternalServerError, Response{
			Code:  http.StatusNotFound,
			Error: fmt.Errorf("no (unique) list found with given ID %d", listID),
		})
		return
	}

	ips := []models.IP{}
	if err := db.DBCon.Select(&ips, "SELECT ID, ListID, INET6_NTOA(Address) AS Address, Network, Comment, Created, Updated FROM ipbl.IPAddress WHERE ListID = ?;", listID); err != nil {
		outputRenderer(http.StatusInternalServerError, Response{
			Code:  http.StatusInternalServerError,
			Error: err,
		})
		c.Error(err)
		return
	}

	outputRenderer(http.StatusOK, Response{
		Code:   http.StatusOK,
		Result: ips,
	})
}

func InsertIPIntoList(c *gin.Context) {
	outputRenderer := getOutputRenderer(c)

	paramListID := c.Param("ListID")
	listID, err := strconv.ParseInt(paramListID, 10, 64)
	if err != nil || listID < 1 {
		err = fmt.Errorf("invalid ListID given")
		outputRenderer(http.StatusBadRequest, Response{
			Code:  http.StatusBadRequest,
			Error: err,
		})
		c.Error(err)
		return
	}

	ip := models.IP{
		Network: -1,
	}
	if err = c.ShouldBindJSON(&ip); err != nil {
		outputRenderer(http.StatusBadRequest, Response{
			Code:   http.StatusBadRequest,
			Error:  err,
			Result: ip,
		})
		c.Error(err)
		return
	}
	ip.ListID = listID

	parsedIP := net.ParseIP(ip.Address)
	if parsedIP == nil {
		outputRenderer(http.StatusBadRequest, Response{
			Code:   http.StatusBadRequest,
			Error:  fmt.Errorf("Address is not a valid IP version 4 nor 6"),
			Result: ip,
		})
		return
	}

	if ip.Network != -1 {
		if parsedIP.To4() != nil { // IPv4
			if ip.Network > 32 || ip.Network < 0 {
				outputRenderer(http.StatusBadRequest, Response{
					Code:   http.StatusBadRequest,
					Error:  fmt.Errorf("cidr notation is invalid for IPv4"),
					Result: ip,
				})
				return
			}
		} else if parsedIP.To16() != nil { // IPv6
			if ip.Network > 128 || ip.Network < 4 {
				outputRenderer(http.StatusBadRequest, Response{
					Code:   http.StatusBadRequest,
					Error:  fmt.Errorf("cidr notation is invalid for IPv6"),
					Result: ip,
				})
				return
			}
		} else {
			outputRenderer(http.StatusBadRequest, Response{
				Code:   http.StatusBadRequest,
				Error:  fmt.Errorf("cidr notation is invalid for IPV4 and IPv6"),
				Result: ip,
			})
			return
		}
	}

	now := time.Now()
	ip.Created = now
	ip.Updated = now
	var result sql.Result
	if result, err = db.DBCon.NamedExec("INSERT INTO ipbl.IPAddress (ListID, Address, Network, Comment) VALUES (:ListID, INET6_ATON(:Address), :Network, :Comment);", &ip); err != nil {
		var status int16
		if sqlErr, ok := err.(*mysql.MySQLError); !ok {
			status = http.StatusInternalServerError
			c.Error(err)
			err = fmt.Errorf("database error occured")
		} else {
			switch sqlErr.Number {
			case 1062:
				status = http.StatusAlreadyReported
				err = fmt.Errorf("address already in list")
			case 1452:
				status = http.StatusConflict
				err = fmt.Errorf("list with given ListID does not exist")
			default:
				c.Error(err)
				err = fmt.Errorf("database error occured")
			}
		}

		outputRenderer(int(status), Response{
			Code:   status,
			Error:  err,
			Result: ip,
		})
		return
	}
	ip.ID, err = result.LastInsertId()
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
		Result: ip,
	})
}

func DeleteIPFromList(c *gin.Context) {
	outputRenderer := getOutputRenderer(c)

	paramListID := c.Param("ListID")
	listID, err := strconv.ParseInt(paramListID, 10, 64)
	if err != nil || listID < 1 {
		err = fmt.Errorf("invalid ListID given")
		outputRenderer(http.StatusBadRequest, Response{
			Code:  http.StatusBadRequest,
			Error: err,
		})
		c.Error(err)
		return
	}

	ip := models.IP{
		Network: -1,
	}
	if err = c.ShouldBindJSON(&ip); err != nil {
		outputRenderer(http.StatusBadRequest, Response{
			Code:   http.StatusBadRequest,
			Error:  err,
			Result: ip,
		})
		c.Error(err)
		return
	}
	ip.ListID = listID
	args := map[string]interface{}{
		"ListID": listID,
	}

	if ip.ID != 0 {
		args["ID"] = ip.ID
	}

	if ip.Address != "" {
		parsedIP := net.ParseIP(ip.Address)
		if parsedIP == nil {
			outputRenderer(http.StatusBadRequest, Response{
				Code:   http.StatusBadRequest,
				Error:  fmt.Errorf("Address is not a valid IP version 4 nor 6"),
				Result: ip,
			})
			return
		}
		args["Address"] = ip.Address

		if ip.Network != -1 {
			if parsedIP.To4() != nil { // IPv4
				if ip.Network > 32 || ip.Network < 0 {
					outputRenderer(http.StatusBadRequest, Response{
						Code:   http.StatusBadRequest,
						Error:  fmt.Errorf("cidr notation is invalid for IPv4"),
						Result: ip,
					})
					return
				}
			} else if parsedIP.To16() != nil { // IPv6
				if ip.Network > 128 || ip.Network < 4 {
					outputRenderer(http.StatusBadRequest, Response{
						Code:   http.StatusBadRequest,
						Error:  fmt.Errorf("cidr notation is invalid for IPv6"),
						Result: ip,
					})
					return
				}
			} else {
				outputRenderer(http.StatusBadRequest, Response{
					Code:   http.StatusBadRequest,
					Error:  fmt.Errorf("cidr notation is invalid for IPV4 and IPv6"),
					Result: ip,
				})
				return
			}
		} else {
			ip.Network = 0
		}
		args["Network"] = ip.Network
	}

	if len(args) < 2 {
		outputRenderer(http.StatusBadRequest, Response{
			Code:   http.StatusBadRequest,
			Error:  fmt.Errorf("only ListID given for deletion, aborting"),
			Result: ip,
		})
		return
	}

	qParams := []interface{}{}
	query := "DELETE FROM ipbl.IPAddress WHERE "
	for key, param := range args {
		if key == "Address" {
			query += key + " = INET6_ATON(?) && "
		} else {
			query += key + " = ? && "
		}
		qParams = append(qParams, param)
	}

	var result sql.Result
	if result, err = db.DBCon.Exec(query[:len(query)-4]+";", qParams...); err != nil {
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
		Message: fmt.Sprintf("%d rows affected", affected),
		Result:  ip,
		Error:   err,
	})
}

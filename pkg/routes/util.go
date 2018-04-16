package routes

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/galexrt/ipbl/pkg/models"
	"github.com/galexrt/ipbl/pkg/render"
	"github.com/gin-gonic/gin"
)

func getOutputRenderer(c *gin.Context) func(code int, obj interface{}) {
	var outputRenderer func(code int, obj interface{})
	switch strings.ToLower(c.Query("renderer")) {
	case "json":
		outputRenderer = c.JSON
	case "identedjson":
		outputRenderer = c.IndentedJSON
	case "yaml":
		outputRenderer = c.YAML
	case "xml":
		outputRenderer = c.XML
	case "raw":
		outputRenderer = func(code int, obj interface{}) {
			if resp, ok := obj.(Response); ok {
				if ips, ok := resp.Result.([]models.IP); ok {
					var ipList bytes.Buffer
					for _, ip := range ips {
						ipList.WriteString(ip.Address)
						if ip.Network > 0 {
							ipList.WriteString("/" + strconv.Itoa(int(ip.Network)))
						}
						ipList.WriteString("\n")
					}
					c.String(code, "%s", ipList.String())
				}
			} else {
				c.String(code, "%+v", obj)
			}
		}
	case "ipset":
		outputRenderer = render.IPSet
	default:
		outputRenderer = c.JSON
	}
	return outputRenderer
}

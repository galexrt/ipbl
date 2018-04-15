package output

import (
	"io"

	"github.com/galexrt/ipbl/pkg/models"
)

type Output interface {
	Render(*io.Writer, models.List, []models.IP)
}

var Factories = make(map[string]func() (Output, error))

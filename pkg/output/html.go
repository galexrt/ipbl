package output

import (
	"io"

	"github.com/galexrt/ipbl/pkg/models"
)

type HTML struct {
	Output
}

func init() {
	Factories["html"] = NewHTML
}

func NewHTML() (Output, error) {
	return &HTML{}, nil
}

func (o *HTML) Render(*io.Writer, models.List, []models.IP) {

}

package output

import (
	"io"

	"github.com/galexrt/ipbl/pkg/models"
)

type JSON struct {
	Output
}

func init() {
	Factories["json"] = NewJSON
}

func NewJSON() (Output, error) {
	return &JSON{}, nil
}

func (o *JSON) Render(*io.Writer, models.List, []models.IP) {

}

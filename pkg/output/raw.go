package output

import (
	"io"

	"github.com/galexrt/ipbl/pkg/models"
)

type Raw struct {
	Output
}

func init() {
	Factories["raw"] = NewRaw
}

func NewRaw() (Output, error) {
	return &Raw{}, nil
}

func (o *Raw) Render(*io.Writer, models.List, []models.IP) {

}

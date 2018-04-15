package output

import (
	"io"

	"github.com/galexrt/ipbl/pkg/models"
)

type IPSet struct {
	Output
}

func init() {
	Factories["ipset"] = NewIPSet
}

func NewIPSet() (Output, error) {
	return &IPSet{}, nil
}

func (o *IPSet) Render(*io.Writer, models.List, []models.IP) {

}

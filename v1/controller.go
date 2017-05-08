package v1

import (
	"github.com/gin-gonic/gin"
)

type baseController struct {
	park string
	lang string
	id   string
}

func (control *baseController) initialization(c *gin.Context) {
	control.park = "1"
	if "land" != c.Param("park") {
		control.park = "2"
	}

	control.lang = c.Param("lang")
	if len(control.lang) < 1 {
		control.lang = "ja"
	}

	control.id = c.Param("id")
}

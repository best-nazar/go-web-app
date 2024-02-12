package model

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Id struct {
	ID int `json:"id" form:"id" gorm:"-" binding:"required,numeric"`
}

func (id Id) String() string {
	v := fmt.Sprintf("%v", id.ID)
	_, err := strconv.Atoi(v)

	if err != nil {
		return "0"
	}

	return v
}

func (obj *Id) ShouldBindId(c *gin.Context) error {
	idStr, exists := c.GetPostForm("ID")
    if !exists {
        // ID not present, check for JSON body
		if err := c.ShouldBind(obj); err != nil {
			return err
		}
    }

    // Convert the ID string to the desired data type (numeric in this case)
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        return err
    }

    // Set the converted ID to the struct
    obj.ID = int(id)

    return nil
}
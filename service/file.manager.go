// SERVICE
// File uploading shared service
package service

import (
	"errors"
	"net/http"
	"path/filepath"
	"slices"
	"strings"

	"github.com/best-nazar/web-app/helpers"
	"github.com/best-nazar/web-app/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FileService interface {
	GetFile(c *gin.Context)
	SaveFile(c *gin.Context) (string, error)
}

type ImageLoader struct {
	sizeLimit 	int64
	path 		string
	extentions 	[]string
}

func (config *ImageLoader) SaveFile(c *gin.Context) (string, error) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, config.sizeLimit)
	file, err := c.FormFile("file")
	destination := ""

    // The file cannot be received.
    if err != nil {
        return "", err
    }

	extension := filepath.Ext(file.Filename)

	if slices.Contains(config.extentions, extension) {
		// Generate random file name for the new upload
		newFileName := uuid.New().String() + extension
		destination = config.path + newFileName
		gin.Default().MaxMultipartMemory = config.sizeLimit

		if err := c.SaveUploadedFile(file, destination); err != nil {
			return "", err
		}
	} else {
		return "", errors.New("Extention is not allowed. Allowed extentioon" + strings.Join(config.extentions, ","))
	}

	return destination, nil
}

func (config *ImageLoader) LoadDefaults(c *gin.Context) {
	appConfig := c.MustGet("config").(model.Config)

	if config.extentions == nil {
		config.extentions = helpers.InterfaceArray(appConfig.ImageConfig["extentions"])
	}

	if config.path == "" {
		config.path = appConfig.ImageConfig["path"].(string)
	}

	if config.sizeLimit == 0 {
		s := appConfig.ImageConfig["size_limit"].(int) // comes from config.yaml as interface of int
		n := int64(s)

		config.sizeLimit = n
	}
}

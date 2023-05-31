package actions

import (
	"io"
	"mime"
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo"
)

// FileServerMiddleware serves uploaded files
func FileServerMiddleware(c buffalo.Context) error {
	filePath := c.Param("path")

	file, err := os.Open(filepath.Join("./uploads", filePath))
	if err != nil {
		return c.Error(404, err)
	}
	defer file.Close()

	// Set appropriate Content-Type header based on the file extension
	contentType := mime.TypeByExtension(filepath.Ext(filePath))
	c.Response().Header().Set("Content-Type", contentType)

	// Serve the file as the response
	_, err = io.Copy(c.Response(), file)
	if err != nil {
		return c.Error(500, err)
	}

	return nil
}

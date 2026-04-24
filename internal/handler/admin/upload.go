package admin

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// saveUploadedImage reads a multipart file from fileField. If a file is present
// it saves it to web/static/uploads/ and returns the public URL.
// If no file is uploaded, it falls back to the urlField form value, then oldURL.
func saveUploadedImage(c *gin.Context, fileField, urlField, oldURL string) (string, error) {
	file, header, err := c.Request.FormFile(fileField)
	if err == nil {
		defer file.Close()

		ext := strings.ToLower(filepath.Ext(header.Filename))
		allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true}
		if !allowed[ext] {
			return oldURL, fmt.Errorf("format gambar tidak didukung (jpg, png, webp, gif)")
		}
		if header.Size > 5<<20 {
			return oldURL, fmt.Errorf("ukuran gambar maksimal 5MB")
		}

		filename := fmt.Sprintf("%d-%s%s", time.Now().UnixMilli(), uuid.New().String()[:8], ext)
		savePath := filepath.Join("web", "static", "uploads", filename)

		if err := os.MkdirAll(filepath.Dir(savePath), 0755); err != nil {
			return oldURL, fmt.Errorf("gagal menyimpan gambar")
		}
		dst, err := os.Create(savePath)
		if err != nil {
			return oldURL, fmt.Errorf("gagal menyimpan gambar")
		}
		defer dst.Close()
		if _, err := io.Copy(dst, file); err != nil {
			return oldURL, fmt.Errorf("gagal menyimpan gambar")
		}

		return "/static/uploads/" + filename, nil
	}

	if v := strings.TrimSpace(c.PostForm(urlField)); v != "" {
		return v, nil
	}
	return oldURL, nil
}

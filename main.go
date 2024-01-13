package main

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type hasName interface {
	Name() string
}

func main() {
	e := echo.New()
	e.Use(middleware.CORS())

	e.GET("/hello", handleHello)
	e.POST("/upload_txt", handleUploadTxt)

	slog.Info("Server is running on http://localhost:8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatal("error starting the server:", err)
	}
}

func handleHello(c echo.Context) error {
	return c.JSON(http.StatusOK, "hello")
}

func handleUploadTxt(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": fmt.Sprintf("error: %v", err),
		})
	}
	urls, err := getURLsFromTxt(file)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": fmt.Sprintf("error: %v", err),
		})
	}

	// 	nowStr := strings.ReplaceAll(time.Now().Format(time.RFC3339), ":", "-")
	// 	folderPath := filepath.Join("files", nowStr)

	// 	if err := os.Mkdir(folderPath, 0777); err != nil {
	// 		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
	// 			"message": fmt.Sprintf("error: %v", err),
	// 		})
	// 	}

	// 	var errURLs []string

	// UrlLoop:
	// 	for i, url := range urls {
	// 		resp, err := http.Get(url)
	// 		if err != nil {
	// 			errURLs = handleUrlError(err, url, errURLs)
	// 			continue
	// 		}
	// 		defer resp.Body.Close()

	// 		// fileType, err := checkFileType(resp.Body)

	// 		// típus ellenőrzése (.fit, .tcx, .gpx)

	// 		f, ok := resp.Body.(hasName)
	// 		if !ok {
	// 			log.Fatal(errors.New("type assertion failed"))
	// 		}
	// 		fmt.Println("Name:", f.Name())

	// 		if resp.StatusCode != http.StatusOK {
	// 			errURLs = handleUrlError(err, url, errURLs)
	// 			continue UrlLoop
	// 		}

	// 		destination := filepath.Join(folderPath, fmt.Sprintf("%d.fit", i))
	// 		out, err := os.Create(destination)
	// 		if err != nil {
	// 			errURLs = handleUrlError(err, url, errURLs)
	// 			continue
	// 		}
	// 		defer out.Close()

	// 		if _, err := io.Copy(out, resp.Body); err != nil {
	// 			errURLs = handleUrlError(err, url, errURLs)
	// 			continue
	// 		}
	// 	}

	// 	if !(len(urls) > len(errURLs)) {
	// 		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
	// 			"status":  "error",
	// 			"message": "coudn't save any files",
	// 		})
	// 	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
		"urls":   urls,
	})
}

// func handleUrlError(err error, url string, errURLs []string) []string {
// 	slog.Error(fmt.Sprintf("URL: %s\nERROR: %v", url, err))
// 	errURLs = append(errURLs, url)
// 	return errURLs
// }

func getURLsFromTxt(file *multipart.FileHeader) ([]string, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	urls := strings.Split(string(data), "\r\n")
	return urls, nil
}

// func checkFileType(reader io.ReadCloser) (string, error) {
// 	buffer := make([]byte, 512)

// 	if _, err := reader.Read(buffer); err != nil && err != io.EOF {
// 		return "", err
// 	}

// 	if _, err := reader.Seek(0, io.SeekStart); err != nil {
// 		return "", err
// 	}

// 	fileType := http.DetectContentType(buffer)
// 	return fileType, nil
// }

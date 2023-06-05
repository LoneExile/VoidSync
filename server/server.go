package server

import (
	"net/http"
	"os"
	"path/filepath"
	"voidsync/api"
	"voidsync/storage"
	"voidsync/sync"
	"voidsync/utils"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

func StartServer(client storage.Storage, syncer sync.Syncer) {
	router := gin.Default()

	apiInstance := api.NewAPI(client, syncer)
	router.POST("/remote-files", func(c *gin.Context) {
		var requestBody struct {
			RemotePath string `json:"remotePath"`
		}
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		remoteFiles, err := apiInstance.GetRemoteFileList(requestBody.RemotePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, remoteFiles)
	})

	router.POST("/sync", func(c *gin.Context) {
		var requestBody struct {
			LocalPath  string `json:"localPath"`
			RemotePath string `json:"remotePath"`
		}
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		err := apiInstance.Sync(requestBody.LocalPath, requestBody.RemotePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Sync successful"})
	})

	router.POST("/download-in-server", func(c *gin.Context) {
		var requestBody struct {
			LocalPath  string `json:"localPath"`
			RemotePath string `json:"remotePath"`
		}
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		err := apiInstance.DownloadObjectsInServer(requestBody.RemotePath, requestBody.LocalPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Download successful"})
	})

	router.POST("/download-all", func(c *gin.Context) {
		var requestBody struct {
			RemotePath string `json:"remotePath"`
		}
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		tmpDir, err := apiInstance.DownloadAllObjects(requestBody.RemotePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer os.RemoveAll(tmpDir)

		c.Header("Content-Type", "application/zip")
		c.Header("Content-Disposition", "attachment; filename=objects.zip")

		err = utils.CreateZipArchive(c.Writer, tmpDir)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	})

	router.POST("/upload-all", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		remotePath := form.Value["remotePath"][0]
		contentType := form.Value["contentType"][0]
		// log.Println(remotePath, contentType)

		tmpDir := utils.MkTmpDir()
		defer os.RemoveAll(tmpDir)

		// get the files
		files := form.File["files"]
		for i, file := range files {
			if file.Filename == ".DS_Store" {
				continue
			}
			relativePath := form.Value["relativePaths"][i]

			// log.Println(filepath.Join(tmpDir, filepath.Dir(relativePath)))
			path := filepath.Join(tmpDir, filepath.Dir(relativePath))
			pathFile := filepath.Join(tmpDir, relativePath)

			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			err := c.SaveUploadedFile(file, pathFile)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		err = apiInstance.UploadDirClient(tmpDir, remotePath, contentType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Upload successful"})
	})

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Replace with the domain of your Next.js app
		AllowCredentials: true,
	})

	http.ListenAndServe(":8080", corsMiddleware.Handler(router))

	// router.Run(":8080")
}

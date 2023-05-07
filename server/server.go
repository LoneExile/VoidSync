package server

import (
	"net/http"
	"voidsync/api"
	"voidsync/storage"
	"voidsync/sync"

	"github.com/gin-gonic/gin"
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

	router.POST("/download-all", func(c *gin.Context) {
		var requestBody struct {
			LocalPath  string `json:"localPath"`
			RemotePath string `json:"remotePath"`
		}
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		err := apiInstance.DownloadAllObjects(requestBody.RemotePath, requestBody.LocalPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Download successful"})
	})

	router.Run(":8080")
}

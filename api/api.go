package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	. "github.com/sys520084/namenode/internal"
)

var namenodeata *NameNodeData = new(NameNodeData)

type UploadNodeForm struct {
	Name  string `form:"name" binding:"required"`
	IsDir bool   `form:"isdir" binding:"required"`
	Size  int    `form:"size" binding:"required"`
}

// Setup Router
func SetupRouter() *gin.Engine {

	r := gin.Default()

	// Ping test
	r.GET("/ping/", func(c *gin.Context) {
		c.String(http.StatusOK, "Pong")

	})

	// dump all info from Namenode
	r.GET("/info/", func(c *gin.Context) {
		c.String(http.StatusOK, "get all nfo done")
	})

	// Head a dataset

	// Add a node to NameNode data
	r.POST("/namenode/:dataset/", func(c *gin.Context) {
		var uploadNodeForm UploadNodeForm
		dataset := c.Param("dataset")
		// get node info from query
		if err := c.ShouldBind(&uploadNodeForm); err != nil {
			fmt.Println(err)
			c.String(http.StatusBadRequest, "%v", err)

		} else {
			size := uploadNodeForm.Size
			isdir := uploadNodeForm.IsDir
			name := uploadNodeForm.Name

			// add node
			node := NewNode(name, size, isdir)
			namenodeata.AddDataSetData(dataset, node)
		}
	})

	// Get dirs info from dataset at Namenode data
	r.GET("/namenode/:dataset/dirinfo/", func(c *gin.Context) {
		c.String(http.StatusOK, "get dir info done")
	})

	// Get files info from dataset at Namenode data
	r.GET("/namenode/:dataset/fileinfo/", func(c *gin.Context) {
		c.String(http.StatusOK, "get dir info done")
	})

	// return response
	return r
}

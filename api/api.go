package api

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	. "github.com/sys520084/namenode/internal"
)

type NameNodeData struct {
	datasets map[string]*NodeTree
	lock     sync.RWMutex
}

func (d *NameNodeData) AddDataSetData(dataset string, name string, size int) {
	d.lock.Lock()
	defer d.lock.Unlock()

	// check dataset
	_, ok := d.datasets[dataset]
	if ok {
		d.datasets[dataset].Nodes = AddToTree(d.datasets[dataset].Nodes, strings.Split(name, "/"), size)
	} else {
		newdata := make(map[string]*NodeTree)
		nameNodeTree := &NodeTree{}
		nameNodeTree.Nodes = AddToTree(nameNodeTree.Nodes, strings.Split(name, "/"), size)
		newdata[dataset] = nameNodeTree
		d.datasets = newdata
	}
}

func (d *NameNodeData) PrintData(dataset string) {
	d.lock.Lock()
	defer d.lock.Unlock()

	fmt.Println(d.datasets[dataset])
}

//var nameNodeTree = &NodeTree{}
var nameNodeData = &NameNodeData{}

type UploadNodeForm struct {
	Name string `form:"name" binding:"required"`
	Size int    `form:"size" binding:"required"`
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
			name := uploadNodeForm.Name
			// add node
			nameNodeData.AddDataSetData(dataset, name, size)
		}
	})

	// Get files info from dataset floder at Namenode data
	r.POST("/namenode/:dataset/fileinfo/", func(c *gin.Context) {
		c.String(http.StatusOK, "get dir info done")
	})

	// return response
	return r
}

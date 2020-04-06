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

func (d *NameNodeData) GetPrefixChildren(dataset string, names []string) []Node {
	_, ok := d.datasets[dataset]
	if ok {
		//d.PrintData(dataset)
		datatree := d.datasets[dataset]
		for _, v := range names {
			for _, p := range datatree.Nodes {
				if p.Name == v {
					return p.Children
				}
			}
		}

	}
	return nil
}

//var nameNodeTree = &NodeTree{}
var nameNodeData = &NameNodeData{}

type UploadNodeForm struct {
	Name string `form:"name" binding:"required"`
	Size int    `form:"size" binding:"required"`
}

type GetNodeChildrenForm struct {
	Prefix string `form:"prefix" binding:"required"`
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
	r.POST("/namenode/:dataset/getprefixnode/", func(c *gin.Context) {
		var getNodeChildrenForm GetNodeChildrenForm
		dataset := c.Param("dataset")

		if err := c.ShouldBind(&getNodeChildrenForm); err != nil {
			fmt.Println(err)
			c.String(http.StatusBadRequest, "%v", err)
		} else {
			prefix := getNodeChildrenForm.Prefix

			nodes := nameNodeData.GetPrefixChildren(dataset, strings.Split(prefix, "/"))
			fmt.Println(nodes)
			c.String(http.StatusOK, "ok")
		}

	})

	// return response
	return r
}

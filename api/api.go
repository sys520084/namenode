package api

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	. "github.com/sys520084/namenode/internal"
	"github.com/sys520084/namenode/internal/middleware"
)

type NameNodeData struct {
	datasets map[string]*NodeTree
	lock     sync.RWMutex
}
type PrefixNode struct {
	Prefix string `json:"prefix"`
}
type NodeOwner struct {
	DisPlayName string `json:"displayname"`
	ID          string `json:"id"`
}
type ContentsNode struct {
	ETag         string `json:"etag"`
	Key          string `json:"key"`
	Marker       string `json:"marker"`
	Owner        NodeOwner
	LastModified time.Time `json:"lastModified"`
	Size         int       `json:"size"`
	StorageClass string    `json:"storageclass"`
}

//var nameNodeTree = &NodeTree{}
var nameNodeData = &NameNodeData{}

func (d *NameNodeData) GetDatasetList() []string {
	datasetList := []string{}
	for k, _ := range d.datasets {
		datasetList = append(datasetList, k)
	}
	return datasetList
}

func (d *NameNodeData) AddDataSetData(dataset string, name string, size int) {
	d.lock.Lock()
	defer d.lock.Unlock()

	// check dataset
	_, ok := d.datasets[dataset]
	if ok {
		d.datasets[dataset].Nodes = AddToTree(d.datasets[dataset].Nodes, strings.Split(name, "/"), size)
	} else {
		if d.datasets == nil {
			d.datasets = make(map[string]*NodeTree)
		}

		nameNodeTree := &NodeTree{}
		nameNodeTree.Nodes = AddToTree(nameNodeTree.Nodes, strings.Split(name, "/"), size)
		d.datasets[dataset] = nameNodeTree
	}

}

func (d *NameNodeData) GetPrefixChildrenNodes(dataset string, names string) []Node {
	_, ok := d.datasets[dataset]
	fmt.Println("starting get prefixchild nodes...")
	if ok {
		if names == "/" {
			return d.datasets[dataset].Nodes
		} else {
			names = strings.TrimRight(names, "/")
			fmt.Println("geting %v children:", names)
			nodes := GetNodeChildren(d.datasets[dataset].Nodes, strings.Split(names, "/"))
			//			fmt.Println("get children nodes is", nodes)
			return nodes
		}
	}
	return nil
}

func (d *NameNodeData) NodeToOutput(nodes []Node, prefix string) ([]PrefixNode, []ContentsNode) {
	var dirs []PrefixNode
	var files []ContentsNode
	var i int

	for i = 0; i < len(nodes); i++ {
		if nodes[i].Size == 0 {
			if prefix == "/" {
				dirs = append(dirs, PrefixNode{Prefix: nodes[i].Name + "/"})
			} else {
				dirs = append(dirs, PrefixNode{Prefix: prefix + nodes[i].Name + "/"})
			}
		}

		if nodes[i].Size > 0 {
			owner := NodeOwner{DisPlayName: "admin",
				ID: "admin",
			}
			if prefix == "/" {
				files = append(files, ContentsNode{ETag: "\"4ba83b7c4afb769bd584709defc09f68\"",
					Key: nodes[i].Name,
					//LastModified: "2020-04-12 16:51:26.722 +0000 UTC,",
					Owner:        owner,
					Size:         nodes[i].Size,
					StorageClass: "STANDARD",
				})
			} else {
				files = append(files, ContentsNode{ETag: "\"4ba83b7c4afb769bd584709defc09f68\"",
					Key:          prefix + nodes[i].Name,
					Size:         nodes[i].Size,
					StorageClass: "STANDARD",
				})
			}
		}

	}
	return dirs, files
}

type UploadNodeForm struct {
	Name string `form:"name" binding:"required"`
	Size int    `form:"size" binding:"required"`
}

type GetNodeChildrenForm struct {
	Prefix string `form:"prefix" binding:"required"`
	Marker string `form:"marker" binding:"required"`
}

// Setup Router
func SetupRouter() *gin.Engine {

	r := gin.Default()
	r.Use(middleware.Logger())

	// Ping test
	r.GET("/ping/", func(c *gin.Context) {
		c.String(http.StatusOK, "Pong")

	})

	// dump all info from Namenode
	r.GET("/info/", func(c *gin.Context) {
		c.String(http.StatusOK, "get all nfo done")
	})

	// get datasets list from namenode
	r.GET("/namenode/", func(c *gin.Context) {
		datasets := nameNodeData.GetDatasetList()
		c.JSON(http.StatusOK, gin.H{"Datasets": datasets})
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
			fmt.Println("upload file name is:", name)
			nameNodeData.AddDataSetData(dataset, name, size)
		}
	})

	// Get files info from dataset floder at Namenode data
	r.POST("/namenode/:dataset/getprefixnode/", func(c *gin.Context) {
		var getNodeChildrenForm GetNodeChildrenForm
		nodes := []Node{}
		var nodeprefix string
		var prefix string
		var marker string
		dataset := c.Param("dataset")

		if err := c.ShouldBind(&getNodeChildrenForm); err != nil {
			fmt.Println(err)
			c.String(http.StatusBadRequest, "%v", err)
		} else {
			prefix = getNodeChildrenForm.Prefix
			marker = getNodeChildrenForm.Marker
			fmt.Println("prefix request is: ", prefix)
			fmt.Println("mark request is: ", marker)
			if prefix == " " && marker == " " {
				nodes = nameNodeData.GetPrefixChildrenNodes(dataset, "/")
				nodeprefix = "/"
			} else {
				if marker != " " {
					nodes = nameNodeData.GetPrefixChildrenNodes(dataset, marker)
					nodeprefix = marker
				} else {
					nodes = nameNodeData.GetPrefixChildrenNodes(dataset, prefix)
					nodeprefix = prefix
				}
			}
		}
		dirPrefixes, fileContents := nameNodeData.NodeToOutput(nodes, nodeprefix)

		c.JSON(http.StatusOK, gin.H{"CommonPrefixes": dirPrefixes,
			"Contents":          fileContents,
			"ContinuationToken": "",
			"Prefix":            prefix,
			"Marker":            marker,
			"Delimiter":         "/",
			"IsTruncated":       false,
			//	"KeyCount":          0,
			"MaxKeys": 1000,
			"Name":    dataset,
			//	"StartAfter":        "",
		})
	})

	// return response
	return r
}

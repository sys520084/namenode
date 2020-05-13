package api

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	. "github.com/sys520084/namenode/internal"
	"github.com/sys520084/namenode/internal/log"
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

func (d *NameNodeData) GetPrefixChildrenNodes(dataset string, names string) ([]Node, bool) {
	set, ok := d.datasets[dataset]
	if !ok {
		return nil, false
	}
	names = strings.TrimRight(names, "/")
	nodes := GetNodeChildren(set.Nodes, strings.Split(names, "/"))
	return nodes, true
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
				dirs = append(dirs, PrefixNode{Prefix: prefix[1:] + nodes[i].Name + "/"})
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
					Key:          prefix[1:] + nodes[i].Name,
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
	Size int    `form:"size"`
}

type GetNodeChildrenForm struct {
	Prefix *string `form:"prefix" binding:"required"`
	Marker *string `form:"marker"`
}

type Prefixnode struct {
	CommonPrefixes    []PrefixNode
	Contents          []ContentsNode
	ContinuationToken string
	Prefix            *string
	Marker            *string
	Delimiter         string
	KeyCount          int
	IsTruncated       bool
	MaxKeys           int
	Name              string
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
		entry := log.GetApiLogEntry(c)
		dataset := c.Param("dataset")

		if err := c.ShouldBind(&getNodeChildrenForm); err != nil {
			entry.Errorf("c.ShouldBind failed:%s", err)
			c.String(http.StatusBadRequest, "%v", err)
			return
		}

		prefix := String(getNodeChildrenForm.Prefix)
		marker := String(getNodeChildrenForm.Marker)
		entry.Infof("request dataset:%s prefix:%s mark:%s", dataset, prefix, marker)

		if len(prefix) > 0 {
			if prefix[0] != '/' {
				prefix = fmt.Sprintf("/%s", prefix)
			}
		}
		if len(marker) > 0 {
			if marker[0] != '/' {
				marker = fmt.Sprintf("/%s", marker)
			}
		}

		if getNodeChildrenForm.Marker != nil {
			entry.Infof("request marker:%s", marker)
			nodes, _ = nameNodeData.GetPrefixChildrenNodes(dataset, marker)
			nodeprefix = marker
		} else if getNodeChildrenForm.Prefix != nil {
			entry.Infof("request prefix:%s", prefix)
			nodes, _ = nameNodeData.GetPrefixChildrenNodes(dataset, prefix)
			nodeprefix = prefix
		} else {
			nodes, _ = nameNodeData.GetPrefixChildrenNodes(dataset, "/")
			nodeprefix = "/"
		}

		dirPrefixes, fileContents := nameNodeData.NodeToOutput(nodes, nodeprefix)

		prefixnode := Prefixnode{
			CommonPrefixes:    dirPrefixes,
			Contents:          fileContents,
			ContinuationToken: "",
			Prefix:            getNodeChildrenForm.Prefix,
			Marker:            getNodeChildrenForm.Marker,
			Delimiter:         "/",
			IsTruncated:       false,
			MaxKeys:           1000,
			KeyCount:          len(fileContents),
			Name:              dataset,
		}

		c.JSON(http.StatusOK, prefixnode)
	})

	// return response
	return r
}

func String(p *string) string {
	if p != nil {
		return *p
	}
	return ""
}

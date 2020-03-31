package internal

import (
	"sync"
)

type NameNodeDataNode struct {
	key   string
	size  int
	IsDir bool
}

func NewNode(name string, size int, isdir bool) *NameNodeDataNode {
	node := &NameNodeDataNode{
		key:   name,
		size:  size,
		IsDir: isdir,
	}
	return node
}

type DataSet struct {
	name  string
	files []*NameNodeDataNode
	dirs  []*NameNodeDataNode
}

type NameNodeData struct {
	datasets map[string]*DataSet
	lock     sync.RWMutex
}

func (d *NameNodeData) AddDataSetData(dataset string, node *NameNodeDataNode) {
	d.lock.Lock()
	defer d.lock.Unlock()

	// check dataset
	_, ok := d.datasets[dataset]

	if ok {
		d.datasets[dataset].name = dataset
		if node.IsDir {
			d.datasets[dataset].dirs = append(d.datasets[dataset].dirs, node)
		} else {
			d.datasets[dataset].files = append(d.datasets[dataset].files, node)
		}
	} else {
		nodelist := []*NameNodeDataNode{}
		newdataset := &DataSet{}
		newdata := make(map[string]*DataSet)
		nodelist = append(nodelist, node)

		if node.IsDir {
			newdataset = &DataSet{
				name: dataset,
				dirs: nodelist,
			}
		} else {
			newdataset = &DataSet{
				name:  dataset,
				files: nodelist,
			}
		}
		newdata[dataset] = newdataset
		d.datasets = newdata
	}
}

func (d *NameNodeData) DataSetDirNum(dataset string) (num int) {
	d.datasets[dataset].name = dataset
	_, ok := d.datasets[dataset]
	num = 0
	if ok {
		num = len(d.datasets[dataset].dirs)
	}
	return num
}

func (d *NameNodeData) DataSetFileNum(dataset string) (num int) {
	d.datasets[dataset].name = dataset
	_, ok := d.datasets[dataset]
	num = 0
	if ok {
		num = len(d.datasets[dataset].files)
	}
	return num
}

func (d *NameNodeData) DataSetNum(dataset string) (num int) {
	d.datasets[dataset].name = dataset
	_, ok := d.datasets[dataset]
	num = 0
	if ok {
		num = len(d.datasets[dataset].files) + len(d.datasets[dataset].dirs)
	}
	return num
}

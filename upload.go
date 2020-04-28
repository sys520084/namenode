package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func myUsage() {
	fmt.Printf("Usage: %s [OPTIONS] argument ...\n", os.Args[0])
	flag.PrintDefaults()
}

// Upload a dataset directory to a new bucket
// Usage:
// upload <parm>
//        -dataset <dataset name> //required
//        -path <path> //required
type Work struct {
	wg   sync.WaitGroup
	work chan func() //任务队列
	f	chan string
}
var WorkPoolCount = 20
var SimpleWork = &Work{}
var DatasetName string
var Prefixpath string
var s3Session *session.Session
// 添加任务
func (p *Work) Add(fn func()) {
	p.work <- fn
}

// 执行
func (p *Work) Run() {
	close(p.work)
	p.wg.Wait()
}

func NewWorkPoll(workers int, fileList []string) *Work{
	SimpleWork.wg = sync.WaitGroup{}
	SimpleWork.work = make(chan func())
	SimpleWork.f = make(chan string, len(fileList))
	SimpleWork.wg.Add(workers)
	// 将文件名放入chan中
	for _, f := range fileList{
		SimpleWork.f <- f
	}
	//根据指定的并发量去读取管道并执行
	for i := 0; i < workers; i++ {
		go func() {
			defer func() {
				// 捕获异常 防止waitGroup阻塞
				if err := recover(); err != nil {
					fmt.Println(err)
					SimpleWork.wg.Done()
				}
			}()
			// 从workChannel中取出任务执行
			for fn := range SimpleWork.work {
				fn()
			}
			SimpleWork.wg.Done()
		}()
	}
	return SimpleWork
}

func RunWorkPool(fileList []string) {
	p := NewWorkPoll(WorkPoolCount, fileList)
	for i := 0; i < len(fileList); i++ {
		p.Add(AddfileToS3)
	}
	p.Run()
}

func main() {
	flag.Usage = myUsage

	datasetPtr := flag.String("dataset", "", "dataset to upload to")
	pathPtr := flag.String("path", "", "path of directory to be synced")
	flag.Parse()
	//	bucket := os.Args[1]
	runtime.GOMAXPROCS(runtime.NumCPU())

	bucket := *datasetPtr
	dataset := *datasetPtr
	floderPath := *pathPtr

	if dataset == "" {
		exitErrorf("Dataset name missing!\nUsage: %s --help", os.Args[0])
	}

	if floderPath == "" {
		exitErrorf("local data path missing!\nUsage: %s --help", os.Args[0])
	}
	// Initialize a session in us-west-1 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.

	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("us-east-1"),
		Endpoint: aws.String("http://abc-storage.ainirobot.net:8080"),
		//Endpoint:         aws.String("http://10.60.78.109:8080"),
		S3ForcePathStyle: aws.Bool(true),
		//LogLevel:         aws.LogLevel(aws.LogDebug | aws.LogDebugWithRequestErrors),
	})
	if err != nil {
		exitErrorf("session.NewSession error: %v", err)
	}

	s3Session = sess
	// Create S3 service client
	svc := s3.New(sess)

	// Create the S3 Bucket
	_, err = svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		exitErrorf("Unable to create dataset %q, %v", bucket, err)
	}

	// Wait until bucket is created before finishing
	fmt.Printf("Waiting for dataset %q to be created...\n", bucket)

	err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		exitErrorf("Error occurred while waiting for bucket to be created, %v", bucket)
	}
	fmt.Printf("Bucket %q successfully created\n", bucket)

	// scan floder

	fmt.Printf("waiting for  %q scaning\n", floderPath)
	fileList, err := ScanDir(floderPath)
	if err != nil {
		exitErrorf("Error occurred while waiting for floder to be scaned, %v", floderPath)
	}
	if len(fileList) <1{
		exitErrorf("no found file in this path. %v", floderPath)
	}
	fmt.Printf("%q successfully scaned\n", floderPath)

	// paddle upload
	DatasetName = dataset
	Prefixpath = floderPath
	RunWorkPool(fileList)

}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
func AddfileToS3() {
	s := s3Session
	fileDir := <- SimpleWork.f
	// open the file for use
	file, err := os.Open(fileDir)
	if err != nil {
		fmt.Println("Open file err: ", err.Error())
		return
	}
	// Get file size and read file content info a buffer
	fileInfo, err := file.Stat()
	if err != nil{
		return
	}
	fmt.Println(fileInfo.Name())
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	_, err = file.Read(buffer)
	if err != nil{
		return
	}

	//define key in bucket
	key := ""
	if fileInfo.IsDir() {
		key = strings.TrimPrefix(fileDir, Prefixpath)
		key = key + "/"
	} else {
		key = strings.TrimPrefix(fileDir, Prefixpath)
	}

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:             aws.String(DatasetName),
		Key:                aws.String(key),
		ACL:                aws.String("private"),
		Body:               bytes.NewReader(buffer),
		ContentLength:      aws.Int64(size),
		ContentType:        aws.String(http.DetectContentType(buffer)),
		ContentDisposition: aws.String("attachment"),
		//ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil{
		return
	}
	// upload info to namenode
	namenodeURL := "http://10.60.78.116:8080"
	UploadURL := namenodeURL + "/" + "namenode" + "/" + DatasetName + "/"
	v := url.Values{}
	v.Set("name", key)
	v.Add("size", strconv.FormatInt(size, 10))
	resp, err:= http.Post(UploadURL, "application/x-www-form-urlencoded", strings.NewReader(v.Encode()))
	if err != nil{
		return
	}
	fmt.Println("upload %v is done", UploadURL)
	file.Close()
	resp.Body.Close()
	time.Sleep(2*time.Second)
}

func ScanDir(floder string) ([]string, error) {
	files := []string{}

	e := filepath.Walk(floder, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			files = append(files, path)
		}
		return err
	})

	if e != nil {
		panic(e)
	}

	for _, file := range files {
		fmt.Println(file)
	}

	return files, nil
}

package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Upload a dataset directory to a new bucket
// Usage:
// upload <parm>
//        -dataset <dataset name> //required
//        -path <path> //required

func main() {
	datasetPtr := flag.String("dataset", "", "dataset to upload to")
	pathPtr := flag.String("path", "", "path of directory to be synced")
	flag.Parse()
	//	if len(os.Args) != 2 {
	//		exitErrorf("Bucket name missing!\nUsage: %s bucket_name", os.Args[0])
	//	}
	//
	//	bucket := os.Args[1]
	bucket := *datasetPtr
	// Initialize a session in us-west-1 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.

	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("us-east-1"),
		Endpoint: aws.String("http://abc-storage.ainirobot.net:8080"),
		//Endpoint:         aws.String("http://10.60.78.109:8080"),
		S3ForcePathStyle: aws.Bool(true),
		//LogLevel:         aws.LogLevel(aws.LogDebug | aws.LogDebugWithRequestErrors),
	})

	// Create S3 service client
	svc := s3.New(sess)

	// Create the S3 Bucket
	_, err = svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		exitErrorf("Unable to create bucket %q, %v", bucket, err)
	}

	// Wait until bucket is created before finishing
	fmt.Printf("Waiting for bucket %q to be created...\n", bucket)

	err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		exitErrorf("Error occurred while waiting for bucket to be created, %v", bucket)
	}
	fmt.Printf("Bucket %q successfully created\n", bucket)

	// scan floder
	dataset := *datasetPtr
	floderPath := *pathPtr
	fmt.Printf("waiting for  %q scaning\n", floderPath)
	fileList, err := ScanDir(floderPath)
	if err != nil {
		exitErrorf("Error occurred while waiting for floder to be scaned, %v", floderPath)
	}
	fmt.Printf("%q successfully scaned\n", floderPath)

	// upload
	for _, file := range fileList {
		fmt.Printf("Waiting for %s successfully upload\n", file)
		err = go AddfileToS3(sess, file, dataset, floderPath)
		if err != nil {
			fmt.Printf("uploading a file err: %s\n", err)
		} else {
			fmt.Printf("upload %s successfully uploaded\n", file)
		}
	}

}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func AddfileToS3(s *session.Session, fileDir string, datasetName string, prefixpath string) error {
	// open the file for use
	file, err := os.Open(fileDir)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get file size and read file content info a buffer
	fileInfo, _ := file.Stat()
	fmt.Println(fileInfo.Name())
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	//define key in bucket
	key := ""
	if fileInfo.IsDir() {
		key = strings.TrimPrefix(fileDir, prefixpath)
		key = key + "/"
	} else {
		key = strings.TrimPrefix(fileDir, prefixpath)
	}

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:             aws.String(datasetName),
		Key:                aws.String(key),
		ACL:                aws.String("private"),
		Body:               bytes.NewReader(buffer),
		ContentLength:      aws.Int64(size),
		ContentType:        aws.String(http.DetectContentType(buffer)),
		ContentDisposition: aws.String("attachment"),
		//ServerSideEncryption: aws.String("AES256"),
	})

	// upload info to namenode
	namenodeURL := "http://10.60.78.116:8080"
	UploadURL := namenodeURL + "/" + "namenode" + "/" + datasetName + "/"
	v := url.Values{}
	v.Set("name", key)
	v.Add("size", strconv.FormatInt(size, 10))
	resp, _ := http.Post(UploadURL, "application/x-www-form-urlencoded", strings.NewReader(v.Encode()))
	defer resp.Body.Close()
	return err
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

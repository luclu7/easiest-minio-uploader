package main

import (
	"github.com/minio/minio-go/v6"
	"log"
	"os"
)

func main() {
	endpoint := os.Getenv("S3_HOST")
	accessKeyID := os.Getenv("S3_KEY")
	secretAccessKey := os.Getenv("S3_SECRET")
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln(err)
	}

	// Make a new bucket called mymusic.
	bucketName := os.Args[1]

	err = minioClient.MakeBucket(bucketName, os.Getenv("S3_REGION"))
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	minioClient.SetBucketPolicy(bucketName, "download")

	// Upload the zip file
	objectName := os.Args[3]
	filePath := os.Args[2]
	contentType := "application/octet-stream"

	// Upload the zip file with FPutObject
	n, err := minioClient.FPutObject(bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, n)
}

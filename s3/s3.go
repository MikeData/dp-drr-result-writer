package s3

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"log"
)

func UploadSource(file io.Reader, filename, aws_region string) error {

	fmt.Println("here")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(aws_region)},
	)
	if err != nil {
		log.Fatal("Unable to create session for s3 upload.")
	}

	uploader := s3manager.NewUploader(sess)
	bucket := "drr-test" // TODO was import-bucket". Should be from env

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		// Print the error and exit.
		log.Fatal("Unable to upload %q to %q, %v", filename, bucket, err)
	}

	return nil

}

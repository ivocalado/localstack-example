package main

import (
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Define the /v1/status endpoint
	router.GET("/v1/summarized_report/:id", func(c *gin.Context) {

		id := c.Param("id")
		presignedUrl, err := generatePresignedURL(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.String(http.StatusOK, presignedUrl)
	})

	// Start the server
	router.Run(":9090")
}

func retrieveSession() (*session.Session, error) {

	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")
	endpoint := os.Getenv("AWS_ENDPOINT_URL")
	s3ForcePathStyle := os.Getenv("S3_FORCE_PATH_STYLE") == "true" // Expecting "true" or "false"

	// Initialize AWS session with credentials and region
	return session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: aws.Bool(s3ForcePathStyle),
	})
}

// Define a function to generate a presigned URL given an object key
func generatePresignedURL(key string) (string, error) {
	bucket := os.Getenv("AWS_BUCKET")
	session, err := retrieveSession()
	if err != nil {
		return "", err
	}

	svc := s3.New(session)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	return req.Presign(15 * time.Minute)
}

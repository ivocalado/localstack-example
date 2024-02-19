package main

import (
	"bytes"
	"net/http"
	"os"

	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Define the /v1/status endpoint
	router.GET("/v1/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "Service is up and running",
		})
	})

	// Define the /v1/submit endpoint
	router.POST("/v1/submit", func(c *gin.Context) {
		var request struct {
			ID      string `json:"id"`
			Payload string `json:"payload"`
		}

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Print the ID
		fmt.Println("Received ID:", request.ID)
		fmt.Printf("Payload: %s\n", request.Payload)

		session, err := retrieveSession()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = uploadToS3(session, request.ID, []byte(request.Payload))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = sendToQueue(session, request.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	// Start the server
	router.Run(":8080")
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

func uploadToS3(session *session.Session, key string, data []byte) error {
	bucket := os.Getenv("AWS_BUCKET")

	svc := s3.New(session)
	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(key),
		Body:          bytes.NewReader(data),
		ContentLength: aws.Int64(int64(len(data))),
	})

	return err
}

func sendToQueue(session *session.Session, payload string) error {
	sqsClient := sqs.New(session)
	queueURL := os.Getenv("AWS_SQS_URL")
	_, err := sqsClient.SendMessage(&sqs.SendMessageInput{
		MessageBody:  aws.String(payload),
		QueueUrl:     &queueURL,
		DelaySeconds: aws.Int64(0),
	})
	return err
}

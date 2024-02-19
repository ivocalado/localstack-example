package main

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
)

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

func pollSQS(session *session.Session) {
	queueURL := os.Getenv("AWS_SOURCE_SQS_URL")

	// Create an SQS service client
	sqsClient := sqs.New(session)

	for {
		// Receive messages from SQS with long polling
		result, err := sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            &queueURL,
			MaxNumberOfMessages: aws.Int64(1),
			WaitTimeSeconds:     aws.Int64(20), // Long polling timeout (adjust as needed)
		})

		if err != nil {
			fmt.Println("Error receiving message:", err)
			time.Sleep(1 * time.Second)
			continue
		}

		// Process received messages
		for _, message := range result.Messages {
			//Read the CSV idata from the message
			bucketObjectId := *message.Body
			csvData, err := readFromS3(session, bucketObjectId)
			if err != nil {
				fmt.Println("Error reading from S3:", err)
				continue
			}

			summarizedCsvData, err := processCsvData(string(csvData))
			if err != nil {
				fmt.Println("Error processing message:", err)
				continue
			}

			err = uploadToS3(session, "summarized_"+bucketObjectId, []byte(summarizedCsvData))
			if err != nil {
				fmt.Println("Error uploading to S3:", err)
				continue
			}

			err = sendToQueue(session, "summarized_"+bucketObjectId)
			if err != nil {
				fmt.Println("Error sending to SQS:", err)
				continue
			}

			// Delete the received message from the queue
			_, err = sqsClient.DeleteMessage(&sqs.DeleteMessageInput{
				QueueUrl:      &queueURL,
				ReceiptHandle: message.ReceiptHandle,
			})
			if err != nil {
				fmt.Println("Error deleting message:", err)
				continue
			}
			fmt.Println("Deleted message from the queue.")
		}

		// Sleep before the next poll
		time.Sleep(1 * time.Second)
	}
}

func processCsvData(payload string) (string, error) {
	// Process the CSV data here
	reader := csv.NewReader(strings.NewReader(payload))
	records, err := reader.ReadAll()
	if err != nil {
		return "", err
	}

	//Iterate over the CSV records,
	//parse the salary (third column) and calculate the total salary per user
	//Define a map to store the total salary per user
	//Iterate over the CSV records, parse the salary (third column) and calculate the total salary per user
	//Define a map to store the total salary per user
	totalSalaryPerUser := make(map[string]float64)
	for _, record := range records {
		if len(record) != 3 {
			return "", errors.New("each record must have 3 columns")
		}
		name := record[1]
		salary, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return "", errors.New("invalid salary format")
		}
		totalSalaryPerUser[name] += salary
	}

	//Generate a csv string with the total salary per user
	var csvData strings.Builder
	csvData.WriteString("Name,Total Salary\n")
	for name, totalSalary := range totalSalaryPerUser {
		csvData.WriteString(fmt.Sprintf("%s,%.2f\n", name, totalSalary))
	}

	return csvData.String(), nil
}

func uploadToS3(session *session.Session, key string, data []byte) error {
	bucket := os.Getenv("AWS_DESTINATION_BUCKET")

	svc := s3.New(session)
	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(key),
		Body:          bytes.NewReader(data),
		ContentLength: aws.Int64(int64(len(data))),
	})

	return err
}

// Read object from S3
func readFromS3(session *session.Session, key string) ([]byte, error) {
	bucket := os.Getenv("AWS_SOURCE_BUCKET")

	svc := s3.New(session)
	result, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	defer result.Body.Close()
	return io.ReadAll(result.Body)
}

func sendToQueue(session *session.Session, payload string) error {
	sqsClient := sqs.New(session)
	queueURL := os.Getenv("AWS_DESTINATION_SQS_URL")
	_, err := sqsClient.SendMessage(&sqs.SendMessageInput{
		MessageBody:  aws.String(payload),
		QueueUrl:     &queueURL,
		DelaySeconds: aws.Int64(0),
	})
	return err
}

func main() {
	session, err := retrieveSession()
	if err != nil {
		fmt.Println("Error creating session:", err)
		return
	}

	go pollSQS(session)

	// Keep the main goroutine running
	select {}
}

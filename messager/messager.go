package messager

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func SendMsg(sourceJson []byte, qURL string) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := sqs.New(sess)

	result, err := svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageBody:  aws.String(string(sourceJson)),
		QueueUrl:     &qURL,
	})

	if err != nil {
		fmt.Println("Error sending message to: "+qURL, err)
	}

	fmt.Println("Success", *result.MessageId)

}

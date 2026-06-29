package main

import (
	// "fmt"
	"fmt"
	"log"
	"net/smtp"
	"sync"
	"time"
)

// id just to debug , we should know alag workers char rahe hai as alag alag consumers run ho rahe hai
func emailWorker(id int , ch chan Recipient, wg *sync.WaitGroup ){
	defer wg.Done();

	for recipient:=range ch{
		// fmt.Println(id,recipient)
		// call smtp service
		// local smtp use karo and it is mailpit
		// acts as real smtp 

		smtpHost:="localhost" // localhost hai jo chal rahi docker ke andhar
		smtpPort:="1025"

		// formattedMessage:=fmt.Sprintf("To: %s\r\nSubject: Test Email\r\n\r\n%s\r\n",recipient.Email,"Just Testing our Email Campaign")
		// msg:=[]byte(formattedMessage)
		msg,err:=executeTemplate(recipient)

		if err!=nil{
			fmt.Printf("Worker: %d Error Parsing template for %s",id,recipient.Email)
			// ADD to DLQ
			continue
		}

		fmt.Printf("Worker %d: Sending Email to %s \n",id,recipient.Email)
		// we dont have auth , means api key and all , so nil
		err=smtp.SendMail(smtpHost+":"+smtpPort,nil,"suyash@gmail.com",[]string{recipient.Email},[]byte(msg))

		if err!=nil{
			log.Fatal(err)
		}

		// dont overwhelm , to be safe for rate limit , 
		time.Sleep(50*time.Millisecond)

		fmt.Printf("Worker %d: Sent Email to %s \n",id,recipient.Email)


		
	}

// 	From: sender@example.com
// To: recipient@example.com
// Subject: Email Subject

// This is the body of the email.
// It can contain multiple lines of text.

}

// we are using the producer consumer architecture where we do this thing , like producer reads the data and consumer processes it
package main

import (
	"bytes"
	"sync"
	"html/template"
)





type Recipient struct{
	Name string
	Email string
}


func main(){

	var wg sync.WaitGroup

	recipientChannel:=make(chan Recipient);
	// we dont have to add the wg here in producer becuase wo already blocked hain  jabtak koi read nahi karta , 
	go func(){
		loadRecipient("./emails.csv",recipientChannel);

	}()
	workerCount:=5

	for i:=1;i<=workerCount;i++{
		wg.Add(1)
		go emailWorker(i,recipientChannel,&wg)
	}

	wg.Wait()
	

	
}

func executeTemplate(r Recipient) (string,error){
	t,err:=template.ParseFiles("email.tmpl")
	if err!=nil{
		return "",err
	}
	var tpl bytes.Buffer
	err= t.Execute(&tpl,r)
	if err!=nil{
		return "",err
	}

	return tpl.String(),nil
}

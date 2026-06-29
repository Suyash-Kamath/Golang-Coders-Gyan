package main

import (
	"encoding/csv"
	// "fmt"
	"os"
)

// capital nahi kiya because same package ke andhar hai so
// Go ki ek speciality hai ki hum error throw nahi , return kartein hai
func loadRecipient(filePath string,ch chan Recipient) error{
	
	defer close(ch)
	
	file,err:=os.Open(filePath)
	if err!=nil{
		return err
	}

	defer file.Close();

	r:=csv.NewReader(file)

	records,err:=r.ReadAll()
	if err!=nil{
		return err
	}
	// returns index, record , and 1: is becuase we dont want column header
	for _,record := range records[1:]{
		// fmt.Println(record)
		// this is slice
		// get and send -> consumer 
		// one more concept , we want queue system , so we use channels  
		ch <- Recipient{
			Name: record[0],
			Email: record[1],
		}
	}

	return nil


}

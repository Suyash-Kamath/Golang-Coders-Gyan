package main

import "fmt"
import "time"
import "sync"


type Result struct{
    Value string
    Err error
    
}

func worker(jobsChan chan string, wg *sync.WaitGroup, resultChan chan Result)  {
    
    defer wg.Done()
    
    
    for job:=range jobsChan{
       time.Sleep(time.Millisecond*50)
        fmt.Printf("Image processed: %s\n",job)
        resultChan <-Result{
        Value:job,
        Err:nil,
    }
    }
    
    fmt.Println("Worker shutting down")
    
    
    
}

// Channel is memory ke andhar shared space where go routines has access of it and does communication through memory 

func main(){
    // fmt.Println("Welcome to Go Concurrency")
    
    jobs:=[]string{
        "image_1.png",
        "image_2.png",
        "image_3.png",
        "image_4.png",
        "image_5.png",
        "image_6.png",
        "image_7.png",
        "image_8.png",
        "image_9.png",
        "image_10.png",
         "image_11.png",
        "image_12.png",
        "image_13.png",
        "image_14.png",
        "image_15.png",
        "image_16.png",
        "image_17.png",
        "image_18.png",
        "image_19.png",
        "image_20.png",
    }
    
    
    var wg sync.WaitGroup
    
    totalWorkers:= 5
    
    resultChan:=make(chan Result,50)
    jobsChan:= make(chan string,len(jobs))
    startTime:=time.Now()
    
   
   
   // problematic hai , what if we have 1000 jobs , 1000 goroutines suru ho jayenge 
   // best practice is to limit worker so as to avoid system performance degradation
   
   
   
//   for _,job:= range jobs{
//       wg.Add(1)
//       go worker(job,&wg,resultChan)
//   }

for i:=1;i<=totalWorkers;i++{
          wg.Add(1)
      go worker(jobsChan,&wg,resultChan)
}
// workers ko run karenge , usko warn karke rakhenge and peeche se job send karna start karenge  , we use channels 

    // new goroutine ke andhar jaake bolenge ki waha pe jaake wait karle kyuki hame jobs start karne hai and ye line ki vajase , wo neeche nahi aanewala hai 
    
    go func(){
        
    wg.Wait()
    close(resultChan)
    }()
    
    // send the jobs 
    
    for i:=0;i<len(jobs);i++{
        jobsChan <-jobs[i]
    }
    
    close(jobsChan)
    
    
    // Channel ke upar read karte jaa rahe ho , close karo warna deadlock aa jayega
    for result:=range resultChan{
        fmt.Printf(" Jobs Completed : %v\n",result)
       
    }
    
    fmt.Printf("It took %s ms.\n",time.Since(startTime))
    
}

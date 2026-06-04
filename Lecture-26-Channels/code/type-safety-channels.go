package main


import "fmt"
import "time"

func emailSender(emailChannel <-chan string,done chan<- bool){
    defer func(){done <-true}()
    /*
    
    
range on a channel does not run forever because the channel is buffered.

range on a channel runs until the channel is closed.


Process all available values.
When the channel becomes empty, it will block and wait.
If new values arrive later, it continues.
If no new values ever arrive, it waits forever.

So the most accurate statement is:

Without closing the channel, range has no way to know it should stop. It can wait forever.

Yes. Exactly the same for unbuffered channels as well.

    */
    for email:=range emailChannel{
        fmt.Println("Sending Email to ", email)
        time.Sleep(time.Second)
    }
    
}

func main(){
    
    emailChannel:=make(chan string,5)
    done:=make(chan bool)
    
       go emailSender(emailChannel,done)
       
    // emailChannel<-"suyash@gmail.com"
    // emailChannel<-"krishna@gmail.com"
    
    
    // fmt.Println(<-emailChannel)
    // fmt.Println(<-emailChannel)
    
    
    for i:=0;i<5;i++{
        emailChannel<- fmt.Sprintf("%d@gmail.com", i)
    }
        fmt.Println("Done sending.")
 
    
    close(emailChannel)
    <-done
    
    
    
    
}

package main

import "fmt"


func main(){


messageChannel:=make(chan string)

messageChannel <- "ping"

message:=<-messageChannel

fmt.Println(message)
}

/*

ERROR!
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan send]:
main.main()
	/tmp/Gf5xN1V2n9/main.go:11 +0x36
exit status 2

=== Code Exited With Errors ===

*/

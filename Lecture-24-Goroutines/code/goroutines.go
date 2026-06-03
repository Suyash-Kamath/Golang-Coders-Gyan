// Online Go compiler to run Golang program online
// Print "Start small. Ship something." message
// https://chatgpt.com/share/6a1fcae2-53c8-83e8-8b6d-f2becfd642e7

package main
import "fmt"
import "time"


// func task(id int){
//     fmt.Println("Doing task ",id)
    
// }
func main() {

for i:=0 ; i<=10;i++{
    // go task(i)
    
   go  func (i int){
       fmt.Println(i) // i is closure here , but good practice is receive karo
    }(i)
}

time.Sleep(time.Second*2)
}

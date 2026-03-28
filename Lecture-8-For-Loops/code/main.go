package main


import "fmt"

func main(){

	fmt.Println("====== While Loop Like Behaviour ========")
	i:=1
	for i<=3 {
		fmt.Println(i)
		// i+=1
		i=i+1
		// i++ , you cannot use this in Println() but can use here
	}
	fmt.Println("====== Infinite Loop ========")
	for{
		fmt.Println("Infinite")
		break
	}
	fmt.Println("====== Classic For Loop ========")

	for i:=0;i<=3;i++{
		fmt.Println("Hi Suyash")
	}

	fmt.Println("====== With Continue ========")

	for i:=0;i<=3;i++{
		if i==2{
			continue
		}
		fmt.Println(i)
	}

	fmt.Println("====== With Range ========")

	for i:=range 3{
		fmt.Println("Suyash",i)
	}

	fmt.Println("===== blank identifier =====")

	for _=range 3{
		fmt.Println("Hi Suyash")
	}

}
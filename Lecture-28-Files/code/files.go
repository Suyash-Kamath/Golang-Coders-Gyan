package main

import (
	// "bufio"
	"fmt"
	"os"
)

func main(){
	/*
	f,err:=os.Open("example.txt")
	if err !=nil{
		//log the error
		panic(err)
	}

	fileInfo,err:=f.Stat()

	if err!=nil{
		panic(err)
	}

	
	fmt.Println("File name ",fileInfo.Name())
	fmt.Println("File or Folder ",fileInfo.IsDir())
	fmt.Println("File Size ",fileInfo.Size())
	fmt.Println("File Mode ",fileInfo.Mode())
	fmt.Println("File Modified at ",fileInfo.ModTime())

	*/

	/*
	// Read file

	f,err:= os.Open("example.txt")
	if err!=nil{
		panic(err)
	}

	defer f.Close()

	buf:= make([]byte,12)
	d, err:=f.Read(buf)

	if err!=nil{
		panic(err)
	}
	for i:=0;i<len(buf);i++{
		fmt.Println("data ",d,string(buf[i]))
	}

*/

/*
	// this loads the entire file content in memory in one go , so don't use it 
	data,err:=os.ReadFile("example.txt")

	if err!=nil{
		panic(err)
	}

	fmt.Println(string(data))

	*/

	/*
	// read folders

	dir,err:=os.Open("../")

	if err!=nil{
		panic(err)
	}

	defer dir.Close()

	fileInfo,err:=dir.ReadDir(-1)

	for _,fileIn:=range fileInfo{
		fmt.Println(fileIn.Name(),fileIn.IsDir())
	}
*/


/*
	// create a file

	f,err:= os.Create("example2.txt")
	if err!=nil{
		panic(err)
	}

	defer f.Close()

	f.WriteString("Hi Go ")
	f.WriteString("Nice Language")

	fmt.Print()

	// file is nothing but binary data 
	// data ko alag tarike se bhi dala sakte hai 

	// bytes ka array bana sakte hai string se

	bytes:=[]byte("Hello Golang") //byte ki slice

	f.Write(bytes)

	*/


	/*

	// Transferring the data from one file to another file via streaming

	// read and write to another file (streaming fashion)

	sourceFile,err:=os.Open("example.txt")

	if err!=nil{
		panic(err)
	}

	defer sourceFile.Close()

	destFile,err:=os.Create("example2.txt")

	if err!=nil{
		panic(err)
	}

	defer destFile.Close()
	// buffer default size is 4096
	reader:=bufio.NewReader(sourceFile)
	writer:=bufio.NewWriter(destFile)

	for{
		b,err:=reader.ReadByte()

		if err!=nil{
			if err.Error() != "EOF"{

				panic(err)
			}
			break
		}


		e:= writer.WriteByte(b)

		if e!=nil{
			panic(e)
		}

	}

	// last me agar kuch bhi data bacha hai to flush kar dete hai

	writer.Flush()

	fmt.Println("Written to a new file successfully")

*/


// Delete a file



	err:= os.Remove("example2.txt")

	if err!=nil{
		panic(err)
	}

	fmt.Println("File Deleted Successfully")



}

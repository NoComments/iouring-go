package main

import (
	"fmt"
	"os"

	"github.com/iceber/iouring-go"
)

const entries uint = 64

var (
	str1 = "str1 str1 str1 str1\n"
	str2 = "str2 str2 str2 str2 str2\n"
)

func main() {
	iour, err := iouring.New(entries)
	if err != nil {
		panic(fmt.Sprintf("new IOURing error: %v", err))
	}

	file, err := os.Create("./tmp")
	if err != nil {
		panic(err)
	}

	writeRequest1 := iouring.RequestWithInfo(iouring.Write(int(file.Fd()), []byte(str1)), "write str1")
	writeRequest2 := iouring.RequestWithInfo(iouring.Pwrite(int(file.Fd()), []byte(str2), uint64(len(str1))), "write str2")

	buffer := make([]byte, len(str1)+len(str2))
	readRequest1 := iouring.RequestWithInfo(iouring.Read(int(file.Fd()), buffer), "read fd to buffer")
	readRequest2 := iouring.RequestWithInfo(iouring.Write(int(os.Stdout.Fd()), buffer), "read buffer to stdout")

	ch := make(chan *iouring.Result, 4)
	err = iour.SubmitLinkRequests(
		[]iouring.Request{
			writeRequest1,
			writeRequest2,
			readRequest1,
			readRequest2,
		},
		ch,
	)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 4; i++ {
		result := <-ch
		info := result.GetRequestInfo().(string)
		fmt.Println(info)
		if err := result.Err(); err != nil {
			fmt.Printf("error: %v\n", err)
		}
	}
}
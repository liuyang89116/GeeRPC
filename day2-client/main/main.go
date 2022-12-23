package main

import (
	"fmt"
	"geerpc"
	"log"
	"net"
	"sync"
	"time"
)

func startServer(addr chan string) {
	// pick a free port
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	geerpc.Accept(l)
}

/*
output:

start rpc server on [::]:33149
&{Foo.Sum 5 } geerpc req 2
&{Foo.Sum 2 } geerpc req 1
&{Foo.Sum 3 } geerpc req 0
&{Foo.Sum 4 } geerpc req 3
reply: geerpc response 4
reply: geerpc response 5
reply: geerpc response 2
reply: geerpc response 3
&{Foo.Sum 1 } geerpc req 4
reply: geerpc response 1
*/
func main() {
	log.SetFlags(0)
	addr := make(chan string)
	go startServer(addr)

	client, _ := geerpc.Dial("tcp", <-addr)
	defer func() {
		_ = client.Close()
	}()

	time.Sleep(time.Second)
	// send requests & receive response
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := fmt.Sprintf("geerpc req %d", i)
			var reply string
			if err := client.Call("Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}

			log.Println("reply:", reply)
		}(i)
	}
	wg.Wait()
}

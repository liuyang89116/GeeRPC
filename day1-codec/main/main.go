package main

import (
	"encoding/json"
	"fmt"
	"geerpc"
	"geerpc/codec"
	"log"
	"net"
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

start rpc server on [::]:35865
&{Foo.Sum 0 } geerpc req 0
reply: geerpc response 0
&{Foo.Sum 1 } geerpc req 1
reply: geerpc response 1
&{Foo.Sum 2 } geerpc req 2
reply: geerpc response 2
&{Foo.Sum 3 } geerpc req 3
reply: geerpc response 3
&{Foo.Sum 4 } geerpc req 4
reply: geerpc response 4
*/
func main() {
	log.SetFlags(0)
	addr := make(chan string)
	go startServer(addr)

	// we create a simple geerpc client
	conn, _ := net.Dial("tcp", <-addr)
	defer func() {
		_ = conn.Close()
	}()

	time.Sleep(time.Second)
	// send options
	_ = json.NewEncoder(conn).Encode(geerpc.DefaultOption)
	cc := codec.NewGobCodec(conn)
	// send request and receive response
	for i := 0; i < 5; i++ {
		h := &codec.Header{
			ServiceMethod: "Foo.Sum",
			Seq:           uint64(i),
		}
		_ = cc.Write(h, fmt.Sprintf("geerpc req %d", h.Seq))
		_ = cc.ReadHeader(h)
		var reply string
		_ = cc.ReadBody(&reply)
		log.Println("reply:", reply)
	}
}

package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
)

func main(){
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("proxy main init...")
	l,err :=net.Listen("tcp",":1088")
	if err !=nil{
		log.Panic(err)
	}
	for{
		client,err := l.Accept()
		if err != nil{
			log.Panic(err)
		}
		go handleClientRequest(client)
	}
}

func handleClientRequest(c net.Conn){
	if c==nil{
		return
	}
	defer c.Close()
	var b [1024]byte
	n,err:=c.Read(b[:])
	if err != nil{
		log.Println(err)
		return
	}
	//fmt.Println(b)
	if b[0]==0x05{
		c.Write([]byte{0x05,0x00})
		n,err=c.Read(b[:])
		//VER CMD RSV   ATYP  DST.ADDR DST.PORT
		//1   1   X'00' 1     Variable 2
		//fmt.Println(b[4:20])
		//fmt.Println("n:",n)
		var host,port string
		switch b[3] {
		case 0x01:
			//fmt.Println("0x01")
			host=net.IPv4(b[4],b[5],b[6],b[7]).String()
		case 0x03:
			//fmt.Println("0x03")
			host=string(b[5:n-2])
		case 0x04:
			//fmt.Println("0x04")
			host = net.IP{b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15], b[16], b[17], b[18], b[19]}.String()
		}
		port=strconv.Itoa(int(b[n-2])<<8|int(b[n-1]))
		log.Printf("net dail address %s:%s",host,port)
		server,err:=net.Dial("tcp",net.JoinHostPort(host,port))
		if err!=nil{
			log.Println(err)
			return
		}
		defer server.Close()
		//VER REP RSV   ATYP BND.ADDR BND.PORT
		//1   1   X'00' 1    Variable 2
		c.Write([]byte{0x05,0x00,0x00,0x01,0x00,0x00,0x00,0x00,0x00,0x00})
		go io.Copy(server,c)
		io.Copy(c,server)
	}
}

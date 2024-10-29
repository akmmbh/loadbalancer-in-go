package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)
type Server interface{
	Address() string
	IsALive() bool
	Server(rw http.ResponseWriter, r *http.Request)
}
type simpleServer struct {
	addr string 
	proxy *httputil.ReverseProxy
}
func newSimpleServer(addr string)*simpleServer{
	serverUrl, err := url.Parse(addr)
	handleErr(err)
	return  &simpleServer{
		addr:addr,
		proxy: httputil.NewSingleHostReverseProxy(serverUrl),
	}
}
type LoadBalancer struct{
	port string
	roundRobinCount int
	servers []Server
}
func NewLoadBalancer(port string, servers []Server)*LoadBalancer{
	return &LoadBalancer{
		port:port ,
		roundRobinCount:0,
		servers: servers,
	}
}

func handleErr(err error){
	if err!=nil{
		fmt.Printf("error:%v\n",err)
		os.Exit(1)
	
	}
}
func (s *simpleServer)Address()string{ return s.addr}
func (s *simpleServer)IsALive()bool{return true}
func( s *simpleServer)Serve(rw http.ResponseWriter, r*Request){
	s.proxy.ServerHTTP(rw,r)
}
func (lb *LoadBalancer)getNexAvailableServer()Server{}{
server:= lb.servers[lb.roundRobinCount%len(lb.servers)]
for !server.Alive(){
	lb.roundRobinCount++
	server=lb.server[lb.roundRobinCount%len(lb.server)]
}
lb.roundRobinCount++
return server
}
func(lb *LoadBalancer)serveProxy(rw http.ResponseWriter, r *http.Request){
targetServer:= lb.getNexAvailableServer()
fmt.Printf("forwarding request to address %q\n",targetServer.Address())
targetServer.Serve(rw,r)

}
func main(){
	servers:=[]Server{
		newSimpleServer("http://www.facebook.com"),
		newSimpleServer("http://www.bing.com"),
		newSimpleServer("http://www.duckduckgo.com"),
	}
	lb:=NewLoadBalancer("8000",servers)
	handleRedirect:=func(rw *http.ResponseWriter, r *http.Request){
		lb.serverProxy(rw,r)
	}
	http.HandleFunc("/",handleRedirect)
	fmt.Printf("serving request at localhots:%s\n",lb.port)
	http.ListenAndServer(":"+lb.port,nil)
}
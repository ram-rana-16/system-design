package main

import(
	"sync"
	"net/http/httputil"
	"sync/atomic"
	"net/http"
	"fmt"
	"flag"
)

// global serverpool
var serverPool ServerPool

// Backaned defines a backed properties
type Backaned struct{
	Name string
	URL string
	Alive bool
	mux sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

// ServerPool is a manager to manage active servers
type ServerPool struct {
	Backaneds []Backaned
	current uint32
}

func (b *Backaned)isAlive() bool {
	var alive bool
	b.mux.RLock()
	alive = b.Alive
	b.mux.RUnlock()
	return alive
}

func (b *Backaned)setAlive(alive bool) {
	b.mux.Lock()
	b.Alive = alive
	b.mux.Unlock()
}

func (s *ServerPool)nextIndex() int {
	return int(atomic.AddUint32(&s.current, uint32(1)) % uint32(len(s.Backaneds)))
}

func (s *ServerPool)getNextPeer() *Backaned {
	next := s.nextIndex()
	l := next + len(s.Backaneds) // iterate all next perr server
	for i := next; i < l; i++ {
		index := i % len(s.Backaneds)
		if s.Backaneds[index].Alive {
			atomic.StoreUint32(&s.current, uint32(index))
			return &s.Backaneds[index]
		}
	}
	return nil
}

func lb(w http.ResponseWriter, req *http.Request) {

}
 

func main() {
	port := flag.Int("port", 8080, "ort to run the server")
	
	server := http.Server{
		Addr: fmt.Sprintf(":%d", *port),
		Handler: http.HandlerFunc(lb),
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
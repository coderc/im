package discovery

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestServiceRegiste(t *testing.T) {
	ctx := context.Background()
	ser, err := NewServiceRegister(&ctx, "/web/node1", &EndPointInfo{
		IP:   "127.0.0.1",
		Port: "0101",
	}, 5)
	if err != nil {
		log.Fatalln(err)
	}
	go ser.ListenLeaseRespChan()
	<-time.After(20 * time.Second)
	ser.Close()
}

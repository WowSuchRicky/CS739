package nfs_client

import (
	"fmt"
	pb "github.com/Ricky54326/CS739/hw2/protos"
)

var CLIENT_WRITE_QUEUE_CAPACITY = int64(10)
var MAX_N_DELAYED_WRITES = int64(100)

type ClientWriteQueue struct {
	Queue []*pb.WriteArgs
	Size  int64
}

func InitializeClientWriteQueue() *ClientWriteQueue {
	return &ClientWriteQueue{
		Size: 0}
}

func (wq *ClientWriteQueue) Reinitialize() {
	wq.Queue = wq.Queue[:0] // empty it
	wq.Size = 0
	fmt.Printf("Inside Reinitialize, just reinitialized queue; size of queue is now: %v\n", wq.Size)
}

func (wq *ClientWriteQueue) InsertWrite(in *pb.WriteArgs) {
	wq.Queue = append(wq.Queue, in)
	wq.Size += in.Count
	fmt.Printf("Inside InsertWrite, just inserted; size of queue is now: %v\n", wq.Size)
}

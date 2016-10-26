package nfs

import (
	"container/list"
	pb "github.com/Ricky54326/CS739/hw2/protos"
)

var CLIENT_WRITE_QUEUE_CAPACITY = int64(100000)

/////////////////
// CLIENT WRITE QUEUE

type ClientWriteQueue struct {
	queue *list.List
	size  int64
}

func InitializeClientWriteQueue() *ClientWriteQueue {
	return &ClientWriteQueue{
		queue: list.New(),
		size:  0}
}

func (wq *ClientWriteQueue) Reinitialize() {
	wq.queue = list.New()
	wq.size = 0
}

func (wq *ClientWriteQueue) PersistAll() {
	// TODO: call NFS commit
	wq.Reinitialize()
}

func (wq *ClientWriteQueue) InsertWrite(in *pb.WriteArgs) {

	// if size of write is greater than queue capacity, issue a stable write and return (don't touch the queue)
	// this works --correctly-- only if  WriteArgs.Data size is not greater than WriteArgs.Count
	if in.Count > CLIENT_WRITE_QUEUE_CAPACITY {
		// TODO: Call NFS STABLE WRITE
	}

	if wq.size+in.Count > CLIENT_WRITE_QUEUE_CAPACITY {
		wq.PersistAll() // this sends NFS COMMIT to server and empties the queue
	}
	wq.queue.PushBack(in)
	wq.size = in.Count

	// TODO: call NFS WRITE w/ UNSTABLE; once receive acknowledgment, return
}

/////////////////
// SERVER WRITE QUEUE

type ServerWriteQueue struct {
	queue *list.List
}

type ServerWriteQueueEntry struct {
	args     *pb.WriteArgs
	filepath string
}

func InitializeServerWriteQueue() *ServerWriteQueue {
	return &ServerWriteQueue{
		queue: list.New()}
}

func (wq *ServerWriteQueue) Reinitialize() {
	wq.queue = list.New()
}

// this should be called for an NFS COMMIT
func (wq *ServerWriteQueue) ExecuteAllWrites() {
	for e := wq.queue.Front(); e != nil; e = e.Next() {
	}
	wq.Reinitialize()
}

// this should be called for an NFS UNSTABLE WRITE
func (wq *ServerWriteQueue) InsertWrite(in *pb.WriteArgs, filepath string) {
	entry := &ServerWriteQueueEntry{
		args:     in,
		filepath: filepath}
	wq.queue.PushBack(entry)
}

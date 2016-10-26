package nfs

import (
	pb "github.com/Ricky54326/CS739/hw2/protos"
)

type ServerWriteQueue struct {
	queue    []*ServerWriteQueueEntry
	size     int
	capacity int
}

type ServerWriteQueueEntry struct {
	args     *pb.WriteArgs
	filepath string
}

func InitializeServerWriteQueue() *ServerWriteQueue {
	return &ServerWriteQueue{
		size: 0,
	}
}

func (wq *ServerWriteQueue) Reinitialize() {
	wq.queue = wq.queue[:0]
	wq.size = 0

}

// this should be called when server receives COMMIT request
func (wq *ServerWriteQueue) ExecuteAllWrites() {
	for i := 0; i < wq.size; i++ {
		// here it is safe to ignore errors (i.e. file doens't exist anymore),
		// as it is no different than normal UNIX write semantics (no guarantees on when
		// writes are persisted)
		StableWrite(wq.queue[i].filepath, wq.queue[i].args)
	}
	wq.Reinitialize()
}

// this should be called when server receives an UNSTABLE WRITE request
func (wq *ServerWriteQueue) InsertWrite(in *pb.WriteArgs, filepath string) {
	entry := &ServerWriteQueueEntry{
		args:     in,
		filepath: filepath}
	wq.size++
	wq.queue = append(wq.queue, entry)
}

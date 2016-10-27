package nfs

import (
	pb "github.com/Ricky54326/CS739/hw2/protos"
)

// n_commit: we will return this to the client after each write or commit; if it's different
// than the last one they received, they know they should clear their buffer

type ServerWriteQueue struct {
	queue      []*ServerWriteQueueEntry
	size       int
	n_commit   int32
	writeverf3 int32
}

type ServerWriteQueueEntry struct {
	args     *pb.WriteArgs
	filepath string
}

func InitializeServerWriteQueue(writeverf3 int32) *ServerWriteQueue {
	return &ServerWriteQueue{
		size:       0,
		n_commit:   0,
		writeverf3: writeverf3}
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
		StableWrite(wq.queue[i].filepath, wq.queue[i].args, wq)
	}
	wq.Reinitialize()
	wq.n_commit += 1
}

// this should be called when server receives an UNSTABLE WRITE request
func (wq *ServerWriteQueue) InsertWrite(in *pb.WriteArgs, filepath string) {
	entry := &ServerWriteQueueEntry{
		args:     in,
		filepath: filepath}
	wq.size++
	wq.queue = append(wq.queue, entry)
}

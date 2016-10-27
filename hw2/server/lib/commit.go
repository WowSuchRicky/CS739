package nfs

import (
	"fmt"
	pb "github.com/Ricky54326/CS739/hw2/protos"
	"syscall"
)

func CommitNFS(in *pb.CommitArgs, wq *ServerWriteQueue) (*pb.CommitReturn, error) {

	wq.ExecuteAllWrites()

	// ensure all changes reach disk
	syscall.Sync()
	if EN_OUTPUT {
		fmt.Printf("Commit: successfully synced changes to disk\n")
	}

	return &pb.CommitReturn{
			Writeverf3: wq.writeverf3,
			NCommit:    wq.n_commit},
		nil
}

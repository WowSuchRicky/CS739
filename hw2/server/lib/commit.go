package nfs

import (
	pb "github.com/Ricky54326/CS739/hw2/protos"
)

func CommitNFS(in *pb.CommitArgs, wq *ServerWriteQueue) (*pb.CommitReturn, error) {
	wq.ExecuteAllWrites()
	return &pb.CommitReturn{
			Writeverf3: wq.writeverf3,
			NCommit:    wq.n_commit},
		nil
}

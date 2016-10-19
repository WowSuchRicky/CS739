package nfs

import (
	pb "github.com/Ricky54326/CS739/hw2/protos"
	"syscall"
)

func RemoveNFS(in *pb.RemoveArgs) (*pb.RemoveReturn, error) {

	// get full path of directory + file
	dir_path, err := InumToPath(int(in.Dirfh.Inode))
	full_path := dir_path + "/" + in.Name

	// TODO: figure out how we should use genum
	// dir_genum := in.Dirfh.Genum

	// delete that file
	err = syscall.Unlink(full_path)

	// TODO: in this case, we're returning err and status, but I think these
	// are essentially the same thing - figure out what we should return
	return &pb.RemoveReturn{Status: 52}, err
}

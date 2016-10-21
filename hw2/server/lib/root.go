package nfs

import (
	"errors"
	"fmt"
	pb "github.com/Ricky54326/CS739/hw2/protos"
	"os"
	"syscall"
)

// return the file handle of the root of server (what we want to make accessible to client)
// for now, we can say that is just the directory the server runs in
func GetRootNFS(in *pb.RootArgs) (*pb.RootReturn, error) {
	root_path := in.Path

	// get inode
	var f_info syscall.Stat_t
	err := syscall.Stat(root_path, &f_info)
	if err != nil {
		return &pb.RootReturn{}, errors.New("file does not exist")
	}
	root_inode := f_info.Ino

	// get genum
	root_genum, err := PathToGen(root_path)
	if err != nil {
		fmt.Println("Failed to get genum of root, fatal error")
		os.Exit(-1)
	}

	return &pb.RootReturn{Fh: &pb.FileHandle{Inode: root_inode, Genum: root_genum}, Attr: StatToAttr(&f_info)}, nil
}

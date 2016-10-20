package nfs

import (
	"errors"
	"fmt"
	pb "github.com/Ricky54326/CS739/hw2/protos"
	"os"
	"syscall"
)

func RemoveNFS(in *pb.RemoveArgs) (*pb.RemoveReturn, error) {
	dir_path, err := InumToPath(int(in.Dirfh.Inode))
	file_path := dir_path + "/" + in.Name

	// if InumToPath finds nothing, then that inode is unused, error
	// TODO: this relies on the fact that InumToPath returns something other than nil if there's no inode
	if err != nil {
		return &pb.RemoveReturn{Status: int32(-1)}, errors.New("inode not found")
	}

	// make sure genum exists and if genum for the found inode doesnt match fh, error
	fh_genum := in.Dirfh.Genum
	fs_genum, err := PathToGen(dir_path)
	if err != nil {
		fmt.Println("Genum not found, error.")
		os.Exit(-1)
	}
	if fh_genum != fs_genum {
		fmt.Printf("fh_genum: %v, fs_genum: %v", fh_genum, fs_genum)
		return &pb.RemoveReturn{Status: int32(-1)}, errors.New("genum mismatch")
	}

	err = syscall.Unlink(file_path)

	status := 0
	if err != nil {
		status = -1
	}

	return &pb.RemoveReturn{Status: int32(status)}, err
}

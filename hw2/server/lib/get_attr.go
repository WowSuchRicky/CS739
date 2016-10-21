package nfs

import (
	"errors"
	"fmt"
	pb "github.com/Ricky54326/CS739/hw2/protos"
	"os"
	"syscall"
)

func GetAttrNFS(in *pb.GetAttrArgs) (*pb.GetAttrReturn, error) {

	// get path for file
	file_path, err := InumToPath(int(in.Fh.Inode))
	if err != nil {
		return &pb.GetAttrReturn{}, errors.New("inode not found")
	}

	// confirm genums are same
	fh_genum := in.Fh.Genum
	fs_genum, err := PathToGen(file_path)
	if err != nil {
		fmt.Println("Genum not found, fatal err.")
		os.Exit(-1)
	}

	if fh_genum != fs_genum {
		return &pb.GetAttrReturn{}, errors.New("genum mismatch")
	}

	// stat the file
	var f_info syscall.Stat_t
	err = syscall.Stat(file_path, &f_info)
	if err != nil {
		return &pb.GetAttrReturn{}, errors.New("stat failed")
	}

	return &pb.GetAttrReturn{Attr: StatToAttr(&f_info)}, err
}

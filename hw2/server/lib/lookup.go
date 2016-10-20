package nfs

import (
	"errors"
	"fmt"
	pb "github.com/Ricky54326/CS739/hw2/protos"
	"os"
	"syscall"
)

func LookupNFS(in *pb.LookupArgs) (*pb.LookupReturn, error) {
	dir_path, err := InumToPath(int(in.Dirfh.Inode))
	full_path_to_file := dir_path + "/" + in.Name

	// error if inode wasn't found
	if err != nil {
		return &pb.LookupReturn{}, errors.New("inode not found")
	}

	// error if genum not found (fatal error, should never happen), or genum mismatch
	fh_genum := in.Dirfh.Genum
	fs_genum, err := PathToGen(dir_path)
	if err != nil {
		fmt.Println("Genum not found, error.")
		os.Exit(-1)
	}
	if fh_genum != fs_genum {
		fmt.Printf("fh_genum: %v, fs_genum: %v", fh_genum, fs_genum)
		return &pb.LookupReturn{}, errors.New("genum mismatch")
	}

	// get inode & genum of the actual file
	var f_info syscall.Stat_t
	err = syscall.Stat(full_path_to_file, &f_info)
	if err != nil {
		return &pb.LookupReturn{}, errors.New("file does not exist")
	}

	ret_inode := f_info.Ino
	ret_genum, err := PathToGen(full_path_to_file)
	if err != nil {
		fmt.Println("File exists but failed to retrieve genum, fatal error")
		os.Exit(-1)
	}

	// TODO: get attributes

	return &pb.LookupReturn{Fh: &pb.FileHandle{Inode: ret_inode, Genum: ret_genum}, Attr: &pb.Attribute{}}, err
}

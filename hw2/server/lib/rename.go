package nfs

import (
	"errors"
	"fmt"
	pb "github.com/Ricky54326/CS739/hw2/protos"
	"os"
	"syscall"
)

func RenameNFS(in *pb.RenameArgs) (*pb.RenameReturn, error) {
	orig_dir_path, err := InumToPath(int(in.Dirfh.Inode))
	dest_dir_path, err2 := InumToPath(int(in.Tofh.Inode))

	orig_full_path := orig_dir_path + "/" + in.Name
	dest_full_path := dest_dir_path + "/" + in.Toname

	// error if either inode wasn't found
	if err != nil || err2 != nil {
		return &pb.RenameReturn{}, errors.New("inode not found")
	}

	// error if either genum not found (fatal error, should never happen), or genum mismatch
	orig_fh_genum := in.Dirfh.Genum
	dest_fh_genum := in.Tofh.Genum

	orig_fs_genum, err := PathToGen(orig_dir_path)
	dest_fs_genum, err2 := PathToGen(dest_dir_path)
	if err != nil || err2 != nil {
		fmt.Println("Genum not found, fatal error.")
		os.Exit(-1)
	}
	if orig_fh_genum != orig_fs_genum || dest_fh_genum != dest_fs_genum {
		return &pb.RenameReturn{}, errors.New("genum mismatch")
	}

	err = syscall.Rename(orig_full_path, dest_full_path)

	if err != nil {
		return &pb.RenameReturn{}, err
	}

	return &pb.RenameReturn{Status: 0}, nil

}

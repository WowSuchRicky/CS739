package nfs

import (
	"errors"
	"fmt"
	pb "github.com/Ricky54326/CS739/hw2/protos"
	"io/ioutil"
	"os"
	"syscall"
)

// we're not using the syscalls directly here, so we're not dealing
// directly with file descriptors
// the reason is because we need to do fseek
// if we wanted to cache open file descriptors, we would just have to do
// that separately

func ReaddirNFS(in *pb.ReaddirArgs) (*pb.ReaddirReturn, error) {

	// readdirargs contains Dirfh and Count (uint64)
	// readdirreturn contains Entries, an array of Dirents, where each Dirent contains Inode (uint64) and Name (string)

	// get path for inode
	dir_path, err := InumToPath(int(in.Dirfh.Inode))
	if err != nil {
		return &pb.ReaddirReturn{}, errors.New("inode not found")
	}

	// ensure genums match
	fh_genum := in.Dirfh.Genum
	fs_genum, err := PathToGen(dir_path)
	if err != nil {
		fmt.Println("Genum of file not found, fatal error.")
		os.Exit(-1)
	}

	if fh_genum != fs_genum {
		return &pb.ReaddirReturn{}, errors.New("genum mismatch")
	}

	// use ioutil.ReadDir with that path to get the entries
	entries, err := ioutil.ReadDir(dir_path)
	if err != nil {
		// TODO: ERROR
	}

	// translate those entires into what we want to return
	var nfs_entries []*pb.Dirent
	nfs_entries = make([]*pb.Dirent, len(entries))

	for i := 0; i < len(entries); i++ {

		// call stat on each entry because we must get inode and know whether the file is a directory or not
		file_path := dir_path + "/" + entries[i].Name()
		var f_info syscall.Stat_t
		err = syscall.Stat(file_path, &f_info)
		if err != nil {
			fmt.Println("Stat failed on directory entry; fatal error")
			os.Exit(-1)
		}

		nfs_entries[i] = &pb.Dirent{Name: entries[i].Name(), Inode: f_info.Ino, Mode: f_info.Mode}
	}

	return &pb.ReaddirReturn{Entries: nfs_entries}, nil
}

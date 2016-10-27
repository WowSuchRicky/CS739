package nfs

import (
	"errors"
	"fmt"
	pb "github.com/Ricky54326/CS739/hw2/protos"
	"os"
	"syscall"
)

func WriteNFS(in *pb.WriteArgs, wq *ServerWriteQueue) (*pb.WriteReturn, error) {
	// we have same input args as Read but also have Data which is an array of bytes
	// see read.go for info on some decision made here

	// get path for the file
	filepath, err := InumToPath(int(in.Fh.Inode))
	if err != nil {
		return &pb.WriteReturn{}, errors.New("inode not found")
	}

	// ensure genums match
	fh_genum := in.Fh.Genum
	fs_genum, err := PathToGen(filepath)
	if err != nil {
		fmt.Println("Genum of file not found, error.")
		os.Exit(-1)
	}

	if fh_genum != fs_genum {
		return &pb.WriteReturn{}, errors.New("genum mismatch")
	}

	// NOTE: above this is the same for both stable and unstable,
	//       because we must notify caller immediately if the filehandle
	//       no longer refers to what they think it does

	if in.Stable {
		fmt.Printf("Stable write\n")
		return StableWrite(filepath, in, wq)
	} else {
		fmt.Printf("Unstable write\n")
		return UnstableWrite(filepath, in, wq)
	}
}

func StableWrite(filepath string, in *pb.WriteArgs, wq *ServerWriteQueue) (*pb.WriteReturn, error) {

	// get file object for that file (not an fd)
	f, err := os.OpenFile(filepath, os.O_WRONLY, 0)
	if err != nil {
		return &pb.WriteReturn{}, err
	}

	// write the data into it starting at in.Offset
	data_to_write := in.Data[0:in.Count]
	n_bytes_written, err := f.WriteAt(data_to_write, in.Offset)
	n_bytes_written = n_bytes_written // to supress compiler warning

	err = f.Sync()
	if err != nil {
		fmt.Printf("Fsync'd in stable write\n")
		return &pb.WriteReturn{}, err
	}

	// get attributes after writing it
	var f_info syscall.Stat_t
	err = syscall.Stat(filepath, &f_info)
	if err != nil {
		fmt.Println("Stat failed, FATAL error")
		os.Exit(-1)
	}

	return &pb.WriteReturn{
			Attr:       StatToAttr(&f_info),
			Writeverf3: wq.writeverf3,
			NCommit:    wq.n_commit},
		err

}

func UnstableWrite(filepath string, in *pb.WriteArgs, wq *ServerWriteQueue) (*pb.WriteReturn, error) {

	wq.InsertWrite(in, filepath)
	for i := 0; i < len(wq.queue); i++ {
		fmt.Printf("Entry %d: %v\n", i, *wq.queue[i])
	}
	return &pb.WriteReturn{
			Writeverf3: wq.writeverf3,
			NCommit:    wq.n_commit},
		nil
}

func ApplyWriteFromBuffer(filepath string, in *pb.WriteArgs) error {

	fmt.Printf("Applying write from buffer\n")

	// get file object for that file (not an fd)
	f, err := os.OpenFile(filepath, os.O_WRONLY, 0)
	if err != nil {
		return err
	}

	// write the data into it starting at in.Offset
	data_to_write := in.Data[0:in.Count]
	_, err = f.WriteAt(data_to_write, in.Offset)
	if err != nil {
		return err
	}

	return nil
}

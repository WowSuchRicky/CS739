package nfs

import (
	"errors"
	"fmt"
	pb "github.com/Ricky54326/CS739/hw2/protos"
	"os"
	"syscall"
)

// we're not using the syscalls directly here, so we're not dealing
// directly with file descriptors
// the reason is because we need to do fseek
// if we wanted to cache open file descriptors, we would just have to do
// that separately

func ReadNFS(in *pb.ReadArgs, wq *ServerWriteQueue) (*pb.ReadReturn, error) {

	// persist all writes before we read; this greatly simplifies things but
	// we could do better
	wq.ExecuteAllWrites()

	filepath, err := InumToPath(int(in.Fh.Inode))
	if err != nil {
		return &pb.ReadReturn{}, errors.New("inode not found")
	}

	fh_genum := in.Fh.Genum
	fs_genum, err := PathToGen(filepath)
	if err != nil {
		fmt.Println("Genum of file not found, fatal error.")
		os.Exit(-1)
	}

	if fh_genum != fs_genum {
		return &pb.ReadReturn{}, errors.New("genum mismatch")
	}

	// open that file to get file object for it (this isn't an fd)
	f, err := os.Open(filepath)
	if err != nil { // propagate the file open error back to client
		fmt.Println("os.Open returned err")
		return &pb.ReadReturn{}, err
	}

	// use fseek to go to offset position in the file
	new_offset, err := f.Seek(in.Offset, 0)
	if err != nil { // propagate fseek error back to client
		fmt.Println("Fseek returned error")
		return &pb.ReadReturn{}, err
	}
	new_offset = new_offset // to supress compiler warning

	// use read to read count bytes
	f_count := in.Count
	data := make([]byte, f_count, f_count)
	n_bytes_read, err := f.Read(data)
	if err != nil {
		fmt.Println("Read returned error")
		return &pb.ReadReturn{}, err
	}

	/*
		// TODO: handle the case here when we're trying to read the entire file...

		// TODO: go through queue, for each write-to-be-committed:
		//  1) check if inode is same as this file; if not, go to next item in queue
		//  2) check if offset + count have any overlap with the boundaries to be read; if not, go to next item in queue
		//  3) calculate the positions in data[] that need to be updated
		//  4) update that in-memory structure
		low_boundary_read := in.Offset             // int64
		high_boundary_read := in.Offset + in.Count // int64

		// is this a full file read?
		if low_boundary_read == 0 && high_boundary_read

		for i := 0; i < wq.size; i++ {
			write_inode := wq.queue[i].args.Fh.Inode
			if write_inode != in.Fh.Inode {
				continue
			}
			fmt.Printf("Not-yet-committed update found for the file being read\n")

			low_boundary_write := wq.queue[i].args.Offset
			high_boundary_write := low_boundary_write + wq.queue[i].args.Count
			write_data := wq.queue[i].args.Data

			if low_boundary_write > high_boundary_read || high_boundary_write < low_boundary_read {
				continue
			}

			// entirely contained within, make the changes
			if low_boundary_write >= low_boundary_read && high_boundary_write <= high_boundary_read {
				fmt.Printf("Write contained within read, updating in-memory\n")
				upper_data := data[high_boundary_write:high_boundary_read]
				data = append(data[low_boundary_read:low_boundary_write],
					write_data[low_boundary_write:high_boundary_write]...)
				data = append(data, upper_data...)
			}

			// write is entirely after the read (i.e. we're appending to the file)
			if low_boundary_write == high_boundary_read {
				data = append(data
			}

		}
	*/

	n_bytes_read = n_bytes_read // to supress compiler warning

	// TODO: if f_count is greater than the number of bytes to read in data, do
	// we need to do anything special?

	// stat file to fill attributes
	var f_info syscall.Stat_t
	err = syscall.Stat(filepath, &f_info)
	if err != nil {
		fmt.Println("Stat failed, FATAL error")
		os.Exit(-1)
	}

	return &pb.ReadReturn{Attr: StatToAttr(&f_info), Data: data}, nil
}

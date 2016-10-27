package main

import (
	"flag"
	"fmt"
	nfsc "github.com/Ricky54326/CS739/hw2/client/lib"
	pb "github.com/Ricky54326/CS739/hw2/protos"
	"google.golang.org/grpc"
	"log"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	_ "bazil.org/fuse/fs/fstestutil"
	"golang.org/x/net/context"
)

// @TODO: Remember that this is the VM IP address
const (
	address                 = "104.197.218.40:50051"
	err_grpc                = "rpc error: code = 14 desc = grpc: the connection is unavailable"
	ENABLE_WRITE_BUFFER_OPT = false
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s MOUNTPOINT\n", os.Args[0])
	flag.PrintDefaults()
}

// Global vars for NFS
var conn_pb pb.NFSClient
var wq *nfsc.ClientWriteQueue

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		usage()
		os.Exit(2)
	}
	mountpoint := flag.Arg(0)

	// GRPC Connection
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	conn_pb = pb.NewNFSClient(conn)

	// Local FS mount at Arg(0)
	c, err := fuse.Mount(
		mountpoint,
		fuse.FSName("helloworld"),
		//fuse.Subtype("hellofs"),
		fuse.LocalVolume(),
		fuse.VolumeName("Hello world!"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// TODO: this is the path RELATIVE to server which we want to make available remotely
	root_path_on_server := "test/"
	root_ret, err := conn_pb.Root(context.Background(), &pb.RootArgs{Path: root_path_on_server})
	nfs_fs := &FS{Fh: root_ret.Fh}

	writeverf3 := int32(-1)
	wq = nfsc.InitializeClientWriteQueue(writeverf3)

	err = fs.Serve(c, nfs_fs)
	if err != nil {
		log.Fatal(err)
	}

	// check if the mount process has an error to report
	<-c.Ready
	if err := c.MountError; err != nil {
		log.Fatal(err)
	}
}

type FS struct {
	Fh *pb.FileHandle
}

var _ fs.FS = (*FS)(nil)

func (f *FS) Root() (fs.Node, error) {

	if nfsc.EN_OUTPUT {
		fmt.Println("Root path: %v", f.Fh)
	}

	return &Dir{Fh: f.Fh}, nil
}

// Dir implements both Node and Handle for the root directory.
type Dir struct {
	Fh *pb.FileHandle
}

var _ fs.Node = (*Dir)(nil)

func (d *Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	r, err := conn_pb.Getattr(context.Background(),
		&pb.GetAttrArgs{
			Fh: d.Fh})

	// GRPC error, retry til things work
	for err != nil && err.Error() == err_grpc {
		r, err = conn_pb.Getattr(context.Background(),
			&pb.GetAttrArgs{Fh: d.Fh})
	}

	// non-grpc error, actual problem here, abort for now..
	if err != nil {
		fmt.Println("Error on file Attr()")
		os.Exit(1)
	}

	// we aren't sending inode in NFS get_attr since it's redundant but if the
	// attr call succeeds, the inode should be the same one that we passed in
	a.Inode = d.Fh.Inode
	a.Mode = os.ModeDir | os.FileMode(r.Attr.Mode)
	a.Size = r.Attr.Size
	a.Uid = r.Attr.Uid
	a.Gid = r.Attr.Gid

	if nfsc.EN_OUTPUT {
		fmt.Printf("Calling attr on dir\n")
	}
	// TODO: ignoring some other stuff in fuse Attr structure, that maybe we want...
	return nil
}

func (d *Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {

	r, err := conn_pb.Lookup(context.Background(),
		&pb.LookupArgs{
			Dirfh: d.Fh,
			Name:  name})

	for err != nil && err.Error() == err_grpc {
		r, err = conn_pb.Lookup(context.Background(),
			&pb.LookupArgs{Dirfh: d.Fh, Name: name})
		//fmt.Printf("Lookup, retrying...err: %v\n", err)
		// TODO: error
	}

	// non-grpc error, return nil
	if err != nil {
		return nil, fuse.ENOENT
	}

	if ModeToBoolIfDir(r.Attr.Mode) {
		return &Dir{Fh: r.Fh}, nil
	} else {
		return &File{Fh: r.Fh}, nil
	}

	return nil, fuse.ENOENT // TODO: is this the correct error to return?

}

func (d *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {

	// TODO: count is meaningless right now; server will always
	// read every entry
	count := 0

	r, err := conn_pb.Readdir(context.Background(),
		&pb.ReaddirArgs{Dirfh: d.Fh, Count: uint64(count)})

	for err != nil && err.Error() == err_grpc {
		r, err = conn_pb.Readdir(context.Background(),
			&pb.ReaddirArgs{Dirfh: d.Fh, Count: uint64(count)})

		//fmt.Printf("retrying ReadDirAll... err: %v\n", err)

	}

	if err != nil {
		return nil, fuse.ENOENT
	}

	// transfer our rpc Dirent into a fuse Dirent
	var dirDirs []fuse.Dirent
	dirDirs = make([]fuse.Dirent, len(r.Entries))
	for i := 0; i < len(r.Entries); i++ {

		fileType := fuse.DT_File
		if ModeToBoolIfDir(r.Entries[i].Mode) {
			fileType = fuse.DT_Dir
			// fmt.Println("File is recognized as being a directory: %v", r.Entries[i].Name)
		}

		dirDirs[i] = fuse.Dirent{Inode: r.Entries[i].Inode, Name: r.Entries[i].Name, Type: fileType}
	}

	return dirDirs, nil
}

func (d *Dir) Create(ctx context.Context, req *fuse.CreateRequest, resp *fuse.CreateResponse) (fs.Node, fs.Handle, error) {

	if nfsc.EN_OUTPUT {
		fmt.Println("Create called")
	}

	attr := &pb.Attribute{}
	attr.Mode = uint32(req.Mode)

	r, err := conn_pb.Create(context.Background(),
		&pb.CreateArgs{
			Dirfh: d.Fh,
			Name:  req.Name,
			Attr:  attr})

	// GRPC error, retry
	for err != nil && err.Error() == err_grpc {
		r, err = conn_pb.Create(context.Background(),
			&pb.CreateArgs{
				Dirfh: d.Fh,
				Name:  req.Name,
				Attr:  attr})
	}

	// non-GRPC error, something actually wrong..
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return &File{}, &File{}, err
	}

	created_file := &File{Fh: r.Newfh, Offset: 0}

	return created_file, created_file, nil
}

func (d *Dir) Mkdir(ctx context.Context, req *fuse.MkdirRequest) (fs.Node, error) {

	if nfsc.EN_OUTPUT {
		fmt.Printf("Mkdir called\n")
	}

	attr := &pb.Attribute{}
	attr.Mode = uint32(req.Mode)

	r, err := conn_pb.Mkdir(context.Background(),
		&pb.MkdirArgs{
			Dirfh: d.Fh,
			Name:  req.Name,
			Attr:  attr})

	// GRPC error, retry
	for err != nil && err.Error() == err_grpc {
		r, err = conn_pb.Mkdir(context.Background(),
			&pb.MkdirArgs{
				Dirfh: d.Fh,
				Name:  req.Name,
				Attr:  attr})
	}

	// non GRPC error, something actually wrong...
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return &Dir{}, err
	}

	return &Dir{Fh: r.Fh}, nil
}

func (d *Dir) Remove(ctx context.Context, req *fuse.RemoveRequest) error {

	if nfsc.EN_OUTPUT {
		fmt.Printf("Remove called\n")
	}

	// request contains Dir boolean, whcih is true if we're removing a dir; need to handle that
	_, err := conn_pb.Remove(context.Background(),
		&pb.RemoveArgs{
			Dirfh: d.Fh,
			Name:  req.Name,
			IsDir: req.Dir})

	// GRPC error, retry
	for err != nil && err.Error() == err_grpc {
		_, err = conn_pb.Remove(context.Background(),
			&pb.RemoveArgs{
				Dirfh: d.Fh,
				Name:  req.Name,
				IsDir: req.Dir})
	}

	// TODO: RemoveReturn in our nfs-like protocol actually
	// return status; we might not need to use it? because we have
	// err - think about it more
	// non-GRPC error, actual issue.
	if err != nil {
		fmt.Printf("Error on remove: %v\n", err)
		return err
	}

	return nil

}

func (d *Dir) Rename(ctx context.Context, req *fuse.RenameRequest, newDir fs.Node) error {

	if nfsc.EN_OUTPUT {
		fmt.Printf("Rename called\n")
	}

	var x *pb.FileHandle
	switch y := newDir.(type) {
	case *File:
		x = y.Fh
	case *Dir:
		x = y.Fh
	}

	_, err := conn_pb.Rename(context.Background(),
		&pb.RenameArgs{
			Dirfh:  d.Fh,
			Name:   req.OldName,
			Tofh:   x,
			Toname: req.NewName})

	// GRPC error, retry..
	for err != nil && err.Error() == err_grpc {
		_, err = conn_pb.Rename(context.Background(),
			&pb.RenameArgs{
				Dirfh:  d.Fh,
				Name:   req.OldName,
				Tofh:   x,
				Toname: req.NewName})
	}

	// non-GRPC Error, actual issue.
	if err != nil {
		fmt.Printf("Error on rename: %v\n", err)
		return err
	}

	return nil
}

// File implements both Node and Handle for the hello file.
type File struct {
	Fh     *pb.FileHandle
	Offset int64
}

var _ fs.Node = (*File)(nil)

func (f *File) Attr(ctx context.Context, a *fuse.Attr) error {
	r, err := conn_pb.Getattr(context.Background(),
		&pb.GetAttrArgs{
			Fh: f.Fh})

	// GRPC error, retry..
	for err != nil && err.Error() == err_grpc {
		r, err = conn_pb.Getattr(context.Background(),
			&pb.GetAttrArgs{Fh: f.Fh})
	}

	// non-GRPC error, something actually wrong..
	if err != nil {
		fmt.Println("Error on file Attr()")
		os.Exit(1)
	}

	// we aren't sending inode in NFS get_attr since it's redundant but if the
	// attr call succeeds, the inode should be the same one that we passed in
	a.Inode = f.Fh.Inode
	a.Mode = os.FileMode(r.Attr.Mode)
	a.Size = r.Attr.Size
	a.Uid = r.Attr.Uid
	a.Gid = r.Attr.Gid

	// it's possibly the file length is longer than getAttr returns, if we wrote
	// but that write hasn't persisted on server yet
	if ENABLE_WRITE_BUFFER_OPT {
		max_size := a.Size
		for i := 0; i < len(wq.Queue); i++ {
			write_inode := wq.Queue[i].Fh.Inode
			if write_inode != f.Fh.Inode {
				continue
			}
			new_size := uint64(wq.Queue[i].Offset + wq.Queue[i].Count)
			if new_size > max_size {
				max_size = new_size
			}
		}
		a.Size = max_size

	}

	// TODO: ignoring some other stuff in fuse Attr structure, that maybe we want...
	return nil
}

func (f *File) ReadAll(ctx context.Context) ([]byte, error) {

	if nfsc.EN_OUTPUT {
		fmt.Printf("Calling readall\n")
	}

	// force all changes on server to disk, resend changes if notice difference
	if ENABLE_WRITE_BUFFER_OPT {
		performCommit()
		wq.Reinitialize()
	}

	// get file size
	r, err := conn_pb.Getattr(context.Background(),
		&pb.GetAttrArgs{
			Fh: f.Fh})

	// GRPC error, retry.
	for err != nil && err.Error() == err_grpc {
		r, err = conn_pb.Getattr(context.Background(),
			&pb.GetAttrArgs{Fh: f.Fh})
	}

	// non-GRPC error, something actually wrong...
	if err != nil {
		fmt.Println("Error on NFS protocol getAttr()")
		os.Exit(1)
	}

	file_size := r.Attr.Size

	r2, err := conn_pb.Read(context.Background(),
		&pb.ReadArgs{
			Fh:     f.Fh,
			Offset: int64(0),
			Count:  int64(file_size)})

	// GRPC error, retry...
	for err != nil && err.Error() == err_grpc {
		r2, err = conn_pb.Read(context.Background(),
			&pb.ReadArgs{
				Fh:     f.Fh,
				Offset: int64(0),
				Count:  int64(file_size)})
	}

	// non-GRPC fatal error
	if err != nil {
		// TODO: what do we return for err
		return nil, fuse.EIO
	}

	if nfsc.EN_OUTPUT {
		log.Printf("read response: %v\n", r2)
		log.Printf("Errors: %v\n", err)
	}

	return r2.Data, nil
}

func (f *File) Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {

	if nfsc.EN_OUTPUT {
		fmt.Println("Write called")
	}

	// STABLE WRITE if can't fit write into queue no matter what
	// also, ALWAYS stable write if we disabled optimization
	if (int64(len(req.Data)) > nfsc.CLIENT_WRITE_QUEUE_CAPACITY) || !ENABLE_WRITE_BUFFER_OPT {
		return performWrite(ctx, req, resp, f, true)
	}

	// COMMIT and empty queue if queue is already too full
	if wq.Size+int64(len(req.Data)) > nfsc.CLIENT_WRITE_QUEUE_CAPACITY {
		performCommit()
		wq.Reinitialize()
	}

	// UNSTABLE WRITE (it will insert into queue)
	return performWrite(ctx, req, resp, f, false)
}

func (f *File) Fsync(ctx context.Context, req *fuse.FsyncRequest) error {

	if nfsc.EN_OUTPUT {
		fmt.Printf("Fsync called, persisting\n")
	}

	performCommit()
	wq.Reinitialize()

	return nil
}

// returns true if it's a directory
func ModeToBoolIfDir(mode uint32) bool {
	return (mode & 0040000) > 0
}

func performWrite(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse, f *File, stable bool) error {

	if nfsc.EN_OUTPUT {
		if stable {
			fmt.Printf("Stable write client request\n")
		} else {
			fmt.Printf("Unstable write client request\n")
		}
	}

	size_data := len(req.Data)
	nfs_req := &pb.WriteArgs{
		Fh:     f.Fh,
		Offset: req.Offset,
		Count:  int64(size_data),
		Data:   req.Data,
		Stable: stable}

	// retry until it sends
	r, err := conn_pb.Write(context.Background(), nfs_req)
	for err != nil && err.Error() == err_grpc {
		r, err = conn_pb.Write(context.Background(), nfs_req)
	}

	// handle retries if server crashed (thus lost its queue) =>
	if wq.Writeverf3 == -1 {
		wq.Writeverf3 = r.Writeverf3
	} else {
		if wq.Writeverf3 != r.Writeverf3 {
			sendAllWritesInQueue()
			wq.Writeverf3 = r.Writeverf3
		}
	}

	if !stable { // add to queue if unstable write
		wq.InsertWrite(nfs_req)
	}

	// non-GRPC error...
	if err != nil {
		fmt.Println("Write error!")
		return err
	}

	// TODO: might need to return size of data written; try Attr??
	resp.Size = size_data
	return nil
}

func performCommit() (*pb.CommitReturn, error) {

	if nfsc.EN_OUTPUT {
		fmt.Printf("Performing commit\n")
	}

	r, err := conn_pb.Commit(context.Background(),
		&pb.CommitArgs{})

	for err != nil && err.Error() == err_grpc {
		r, err = conn_pb.Commit(context.Background(),
			&pb.CommitArgs{})
	}

	if err != nil {
		fmt.Println("Commit error!")
		return nil, err
	}

	// check writeverf3; if receive different one than expected,
	// resend queue and then recommit
	if wq.Writeverf3 == -1 {
		wq.Writeverf3 = r.Writeverf3
	} else {
		if wq.Writeverf3 != r.Writeverf3 {
			sendAllWritesInQueue()
			wq.Writeverf3 = r.Writeverf3
			return performCommit() // try to commit again
		}
	}

	return nil, nil
}

func sendAllWritesInQueue() {

	if nfsc.EN_OUTPUT {
		fmt.Printf("Sending all writes in queue again\n")
	}

	// send each one
	for i := 0; i < len(wq.Queue); i++ {
		nfs_req := wq.Queue[i]

		if nfsc.EN_OUTPUT {
			fmt.Printf("Sending a write in queue again\n")
		}

		_, err := conn_pb.Write(context.Background(), nfs_req)
		for err != nil && err.Error() == err_grpc {
			_, err = conn_pb.Write(context.Background(), nfs_req)
		}
	}
}

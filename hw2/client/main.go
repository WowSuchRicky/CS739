// Hellofs implements a simple "hello world" file system.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	pb "github.com/Ricky54326/CS739/hw2/protos"
	"google.golang.org/grpc"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	_ "bazil.org/fuse/fs/fstestutil"
	"golang.org/x/net/context"
)

// @TODO: Change
const (
	address = "localhost:50051"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s MOUNTPOINT\n", os.Args[0])
	//flag.PrintDefaults()
}

// Global vars for NFS
var conn_pb pb.NFSClient

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
		fuse.Subtype("hellofs"),
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

// var _ fs.FS = (*FS)(nil)

func (f *FS) Root() (fs.Node, error) {
	fmt.Println("Root path: %v", f.Fh)
	return &Dir{Fh: f.Fh}, nil
}

// Dir implements both Node and Handle for the root directory.
type Dir struct {
	Fh *pb.FileHandle
}

// var _ fs.Node = (*Dir)(nil)

// TODO: attribute for directory
func (d *Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 4291
	a.Mode = os.ModeDir | 0555
	return nil
}

// TODO: lookup for directory
func (d *Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {

	r, err := conn_pb.Lookup(context.Background(),
		&pb.LookupArgs{
			Dirfh: d.Fh,
			Name:  name})

	if err != nil {
		// TODO: error
	}

	if ModeToBoolIfDir(r.Attr.Mode) {
		return &Dir{Fh: r.Fh}, nil
	} else {
		return &File{Fh: r.Fh}, nil
	}

	return nil, fuse.ENOENT // TODO: is this the correct error to return?

}

func (d *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {

	// TODO: count is meaningless right now, might need to change that on server side?
	count := 10
	r, err := conn_pb.Readdir(context.Background(),
		&pb.ReaddirArgs{Dirfh: d.Fh, Count: uint64(count)})

	if err != nil {
		// TODO: handle errors
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

// File implements both Node and Handle for the hello file.
type File struct {
	Fh *pb.FileHandle
}

func (f *File) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 2
	a.Mode = 0444
	a.Size = uint64(1)

	r, err := conn_pb.Getattr(context.Background(),
		&pb.GetAttrArgs{
			Fh: f.Fh})

	if err != nil {
		fmt.Println("Error on file Attr()")
		os.Exit(-1)
	}

	// we aren't sending inode in NFS get_attr since it's redundant but if the
	// attr call succeeds, the inode should be the same one that we passed in
	a.Inode = f.Fh.Inode
	a.Mode = os.FileMode(r.Attr.Mode)
	a.Size = r.Attr.Size
	a.Uid = r.Attr.Uid
	a.Gid = r.Attr.Gid

	// TODO: ignoring some other stuff in fuse Attr structure, that maybe we want...
	return nil
}

func (f *File) ReadAll(ctx context.Context) ([]byte, error) {

	// get file size
	r, err := conn_pb.Getattr(context.Background(),
		&pb.GetAttrArgs{
			Fh: f.Fh})

	if err != nil {
		fmt.Println("Error on NFS protocol getAttr()")
		os.Exit(-1)
	}

	file_size := r.Attr.Size

	r2, err := conn_pb.Read(context.Background(),
		&pb.ReadArgs{
			Fh:     f.Fh,
			Offset: int64(0),
			Count:  int64(file_size)})

	if err != nil {
		// TODO: what do we return for err
		return nil, fuse.EIO
	}

	log.Printf("read response: %v\n", r2)
	log.Printf("Errors: %v\n", err)

	return r2.Data, nil
}

// returns true if it's a directory
func ModeToBoolIfDir(mode uint32) bool {
	return (mode & 0040000) > 0
}

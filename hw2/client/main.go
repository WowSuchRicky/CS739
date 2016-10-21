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

	// call our NFS mount function to retrieve the filehandle of server's root
	// we associate that file handle with the FS structure

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

// FS implements the hello world file system.
type FS struct {
	Fh *pb.FileHandle
}

var _ fs.FS = (*FS)(nil)

func (f *FS) Root() (fs.Node, error) {
	fmt.Println("Root path: %v", f.Fh)
	return &Dir{Fh: f.Fh}, nil
}

// Dir implements both Node and Handle for the root directory.
type Dir struct {
	Fh *pb.FileHandle
}

var _ fs.Node = (*Dir)(nil)

// TODO: attribute for directory
func (d *Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 4291
	a.Mode = os.ModeDir | 0555
	return nil
}

// TODO: lookup for directory
func (d *Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	/*	if name == "hello" {
			return File{}, nil
		}
		return nil, fuse.ENOENT
	*/

	r, err := conn_pb.Lookup(context.Background(),
		&pb.LookupArgs{
			Dirfh: d.Fh,
			Name:  name})

	// TODO: we're assuming the result is a file, but it could be a dir, in which case
	// we need to return a &Dir
	if err == nil {
		return &File{Fh: r.Fh}, nil
	}

	return nil, fuse.ENOENT // TODO: is this the correct error to return?

}

// var dirDirs = []fuse.Dirent{
//	{Inode: 2, Name: "hello", Type: fuse.DT_File},
// }

// TODO: readdir (need to implement on server side too)
func (d *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {

	// TODO: count is meaningless right now, might need to change that on server side
	count := 10
	count = count

	r, err := conn_pb.Readdir(context.Background(),
		&pb.ReaddirArgs{Dirfh: d.Fh, Count: 10})

	if err != nil {
		// TODO: handle errors
	}

	var dirDirs []fuse.Dirent
	dirDirs = make([]fuse.Dirent, len(r.Entries))

	for i := 0; i < len(r.Entries); i++ {
		dirDirs[i] = fuse.Dirent{Inode: r.Entries[i].Inode, Name: r.Entries[i].Name, Type: fuse.DT_File}
	}

	return dirDirs, nil
}

// File implements both Node and Handle for the hello file.
type File struct {
	Fh *pb.FileHandle
}

const greeting = "hello, world\n"

func (f *File) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 2
	a.Mode = 0444
	a.Size = uint64(len(greeting))
	return nil
}

func (f *File) ReadAll(ctx context.Context) ([]byte, error) {

	// TODO: need to know file size, which we put in count

	r, err := conn_pb.Read(context.Background(),
		&pb.ReadArgs{
			Fh:     f.Fh,
			Offset: int64(0),
			Count:  int64(5)})

	if err != nil {
		// TODO: what do we return for err
		return nil, fuse.EIO
	}

	log.Printf("read response: %v\n", r)
	log.Printf("Errors: %v\n", err)

	return r.Data, nil
}

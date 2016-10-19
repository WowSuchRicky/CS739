/*
 *
 * Copyright 2015, Google Inc.
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are
 * met:
 *
 *     * Redistributions of source code must retain the above copyright
 * notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above
 * copyright notice, this list of conditions and the following disclaimer
 * in the documentation and/or other materials provided with the
 * distribution.
 *     * Neither the name of Google Inc. nor the names of its
 * contributors may be used to endorse or promote products derived from
 * this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package main

import (
	"log"
	"os"
	"strconv"

	pb "github.com/Ricky54326/CS739/hw2/protos"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server. (taken from helloworld)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewNFSClient(conn)

	if len(os.Args) < 2 {
		log.Printf("lookup/create/remove/read/write [arg0] [arg1] [...]\n")
		os.Exit(1)
	}

	call := os.Args[1]
	if call == "lookup" {

		// this is some temporary crap for testing purpose
		var inode int64
		var genum int64
		var name string

		inode = 1052236
		genum = 0
		name = "test.txt"

		if len(os.Args) < 5 {
			//log.Printf("lookup inode genum filename\n")
			//os.Exit(1)
		} else {
			inode, _ = strconv.ParseInt(os.Args[2], 0, 32)
			genum, _ = strconv.ParseInt(os.Args[3], 0, 32)
			name = os.Args[4]
		}

		r, _ := c.Lookup(context.Background(),
			&pb.LookupArgs{
				Dirfh: &pb.FileHandle{Inode: uint64(inode), Genum: uint64(genum)},
				Name:  name})

		log.Printf("lookup response: %v\n", r)

	} else if call == "create" {
		var inode int64
		var genum int64
		var name string

		inode = 1052236
		genum = 0
		name = "test_new.txt"

		if len(os.Args) < 5 {
			//log.Printf("create dir_inode dir_genum filename [attribute, add later]\n")
			//os.Exit(1)
		} else {
			inode, _ = strconv.ParseInt(os.Args[2], 0, 32)
			genum, _ = strconv.ParseInt(os.Args[3], 0, 32)
			name = os.Args[4]
		}
		r, _ := c.Create(context.Background(),
			&pb.CreateArgs{
				Dirfh: &pb.FileHandle{Inode: uint64(inode), Genum: uint64(genum)},
				Name:  name,
				Attr:  &pb.Attribute{}})

		log.Printf("create response: %v\n", r)

	} else if call == "remove" {
		if len(os.Args) < 5 {
			log.Printf("remove dir_inode dir_genum filename\n")
			os.Exit(1)
		}
		inode, _ := strconv.ParseInt(os.Args[2], 0, 32)
		genum, _ := strconv.ParseInt(os.Args[3], 0, 32)
		name := os.Args[4]

		r, _ := c.Remove(context.Background(),
			&pb.RemoveArgs{
				Dirfh: &pb.FileHandle{Inode: uint64(inode), Genum: uint64(genum)},
				Name:  name})

		log.Printf("remove response: %v", r)

	} else if call == "read" {
		if len(os.Args) < 6 {
			log.Printf("read inode genum offset count\n")
			os.Exit(1)
		}
		inode, _ := strconv.Atoi(os.Args[2])
		genum, _ := strconv.Atoi(os.Args[3])
		offset, _ := strconv.Atoi(os.Args[4])
		count, _ := strconv.Atoi(os.Args[5])

		r, _ := c.Read(context.Background(),
			&pb.ReadArgs{
				Fh:     &pb.FileHandle{Inode: uint64(inode), Genum: uint64(genum)},
				Offset: int32(offset),
				Count:  int32(count)})

		log.Printf("read response: %v\n", r)

	} else if call == "write" {
		if len(os.Args) < 6 {
			log.Printf("write inode genum offset count\n")
			os.Exit(1)
		}
		inode, _ := strconv.Atoi(os.Args[2])
		genum, _ := strconv.Atoi(os.Args[3])
		offset, _ := strconv.Atoi(os.Args[4])
		count, _ := strconv.Atoi(os.Args[5])
		data := []byte{1, 2, 3}

		r, _ := c.Write(context.Background(),
			&pb.WriteArgs{
				Fh:     &pb.FileHandle{Inode: uint64(inode), Genum: uint64(genum)},
				Offset: int32(offset),
				Count:  int32(count),
				Data:   data})

		log.Printf("write response: %v\n", r)

	} else {
		log.Printf("invalid args\n")
	}

}

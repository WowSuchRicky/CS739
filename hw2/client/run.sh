#!/bin/bash

sudo umount test
go build
./client test

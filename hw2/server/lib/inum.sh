#!/bin/sh


if [ -z "$1" ]
then
    echo "usage: ./inode.sh <inode>"
    exit 1
fi

# print file associated with inode (ignoring hard linked / other files)
sudo debugfs -R "ncheck $1" /dev/sda1 2> /dev/null | head -n 2 | tail -n 1 | awk '{ print $2 }'

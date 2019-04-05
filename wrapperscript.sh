#!/bin/bash

# turn on bash's job control
set -m

# Start the primary process and put it in the background
./main loadbalancer 80 &

# Start the helper process
./main appserver 8000

# the my_helper_process might need to know how to wait on the
# primary process to start before it does its work and returns



#!/usr/bin/python
import matplotlib.pyplot as plt
import numpy as np
import math

def mean(numbers):
    return float(sum(numbers)) / max(len(numbers), 1)

def med(numbers):
    numbers.sort()
    return numbers[len(numbers)/2]

with open("no_opt/grpc/rtt.txt") as f:
    grpc_rtt_times = [int(t) for t in f.read().strip().split('\n')]

with open("no_opt/thrift/rtt.txt") as f:
    thrift_rtt_times = [int(t) for t in f.read().strip().split('\n')]

with open("no_opt/grpc/loc_rtt.txt") as f:
    grpc_loc_rtt_times = [int(t) for t in f.read().strip().split('\n')]

with open("no_opt/thrift/loc_rtt.txt") as f:
    thrift_loc_rtt_times = [int(t) for t in f.read().strip().split('\n')]


f_grpc_rtt = np.median(np.array(grpc_rtt_times[::2]))/1000000
s_grpc_rtt = np.median(np.array(grpc_rtt_times[1::2]))/1000000

f_thr_rtt = np.median(np.array(thrift_rtt_times[::2]))/1000000
s_thr_rtt = np.median(np.array(thrift_rtt_times[1::2]))/1000000

f_grpc_loc = np.median(np.array(grpc_loc_rtt_times[::2]))/1000000
s_grpc_loc = np.median(np.array(grpc_loc_rtt_times[1::2]))/1000000

f_thr_loc = np.median(np.array(thrift_loc_rtt_times[::2]))/1000000
s_thr_loc = np.median(np.array(thrift_loc_rtt_times[1::2]))/1000000

grpc_meds = (f_grpc_rtt, s_grpc_rtt, f_grpc_loc, s_grpc_loc)
thrift_meds = (f_thr_rtt, s_thr_rtt, f_thr_loc, s_thr_loc)

ind = np.arange(4)
width = .35

fig, ax = plt.subplots()
grpc = ax.bar(ind, grpc_meds, width, color='blue')
thrift = ax.bar(ind + width, thrift_meds, width, color='red')

ax.set_title('gRPC vs Thrift Round Trip Times')
ax.set_ylabel('Time (ms)')
ax.set_xticks(ind + width)
ax.set_xticklabels(('Remote first', 'Remote second', 'Local first', 'Local second'))
ax.legend([grpc, thrift], ['grpc', 'thrift'])

plt.savefig("rtt.png")

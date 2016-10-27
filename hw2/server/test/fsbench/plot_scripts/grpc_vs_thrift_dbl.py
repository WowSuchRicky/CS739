#!/usr/bin/python
import matplotlib.pyplot as plt

def mean(numbers):
    return float(sum(numbers)) / max(len(numbers), 1)

def med(numbers):
    numbers.sort()
    return numbers[len(numbers)/2]

dbl_grpc_times = []
for i in range(31):
    dbl_grpc_times.append([])
    with open("no_opt/grpc/dbl_times/" + str(i) + ".txt") as f:
        dbl_grpc_times[i] = [int(t) for t in f.read().strip().split('\n')]

dbl_thrift_times = []
for i in range(31):
    dbl_thrift_times.append([])
    with open("opt/thrift/dbl_times/" + str(i) + ".txt") as f:
        dbl_thrift_times[i] = [int(t) for t in f.read().strip().split('\n')]


avg_dbl_grpc_times = []
for i in dbl_grpc_times:
    avg_dbl_grpc_times.append(med(i[::2]))


print avg_dbl_grpc_times

avg_dbl_thrift_times = []
for i in dbl_thrift_times:
    avg_dbl_thrift_times.append(med(i[::2]))


print avg_dbl_thrift_times

bla_grpc_times = []
for i in range(31):
    bla_grpc_times.append([])
    with open("no_opt/grpc/bla_times/" + str(i) + ".txt") as f:
        bla_grpc_times[i] = [int(t) for t in f.read().strip().split('\n')]

bla_thrift_times = []
for i in range(31):
    bla_thrift_times.append([])
    with open("no_opt/thrift/bla_times/" + str(i) + ".txt") as f:
        bla_thrift_times[i] = [int(t) for t in f.read().strip().split('\n')]


avg_bla_grpc_times = []
for i in bla_grpc_times:
    avg_bla_grpc_times.append(med(i[::2]))


print avg_bla_grpc_times

avg_bla_thrift_times = []
for i in bla_thrift_times:
    avg_bla_thrift_times.append(med(i[::2]))


print avg_bla_thrift_times



dg = plt.scatter(range(31), avg_dbl_grpc_times, color='blue')
dt = plt.scatter(range(31), avg_dbl_thrift_times, color='orange')
bg = plt.scatter(range(31), avg_bla_grpc_times, color='red')
bt = plt.scatter(range(31), avg_bla_thrift_times, color='green')

plt.xlim((-1, 37))

plt.title('Grpc vs Thrift double and struct marshalling')

plt.xlabel('log2(Argument value)')
plt.ylabel('Time (ns)')

plt.legend([dg, bg, dt, bt], ['grpc double', 'grpc struct', 'thrift double',
'thrift struct'])
plt.savefig("grpc_vs_thrift_dbl_and_bla.png")

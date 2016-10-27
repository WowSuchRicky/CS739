#!/usr/bin/python
import matplotlib.pyplot as plt
import math

def mean(numbers):
    return float(sum(numbers)) / max(len(numbers), 1)

def med(numbers):
    numbers.sort()
    return numbers[len(numbers)/2]

grpc_times = []
for i in range(21):
    grpc_times.append([])
    with open("no_opt/grpc/str_times/" + str(i) + ".txt") as f:
        grpc_times[i] = [math.log(int(t),2) for t in f.read().strip().split('\n')]

avg_first_grpc_times = []
avg_second_grpc_times = []
for i in grpc_times:
    avg_first_grpc_times.append(med(i[::2]))
    avg_second_grpc_times.append(med(i[1::2]))


print avg_first_grpc_times
print avg_second_grpc_times


thrift_times = []
for i in range(21):
    thrift_times.append([])
    with open("no_opt/thrift/str_times/" + str(i) + ".txt") as f:
        thrift_times[i] = [math.log(int(t), 2) for t in f.read().strip().split('\n')]

avg_first_thrift_times = []
avg_second_thrift_times = []
for i in thrift_times:
    avg_first_thrift_times.append(med(i[::2]))
    avg_second_thrift_times.append(med(i[1::2]))


print avg_first_thrift_times
print avg_second_thrift_times

fg = plt.scatter(range(21), avg_first_grpc_times, color='blue')
sg = plt.scatter(range(21), avg_second_grpc_times, color='red')
ft = plt.scatter(range(21), avg_first_thrift_times, color='orange')
st = plt.scatter(range(21), avg_second_thrift_times, color='green')

plt.xlim((-1, 30))

plt.title('gRPC vs Thrift String Marshalling')
plt.xlabel('log2(Argument size)')
plt.ylabel('log2(Time (ns))')
plt.legend([fg, sg, ft, st], ['grpc 1st', 'grpc 2nd', 'thrift 1st', 'thrift 2nd'])
plt.savefig("str_marshal.png")

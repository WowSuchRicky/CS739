#!/usr/bin/python
import matplotlib.pyplot as plt

def mean(numbers):
    return float(sum(numbers)) / max(len(numbers), 1)

def med(numbers):
    numbers.sort()
    return numbers[len(numbers)/2]

grpc_times = []
for i in range(31):
    grpc_times.append([])
    with open("no_opt/grpc/int_times/" + str(i) + ".txt") as f:
        grpc_times[i] = [int(t) for t in f.read().strip().split('\n')]

thrift_times = []
for i in range(31):
    thrift_times.append([])
    with open("no_opt/thrift/int_times/" + str(i) + ".txt") as f:
        thrift_times[i] = [int(t) for t in f.read().strip().split('\n')]


avg_first_grpc_times = []
avg_second_grpc_times = []
for i in grpc_times:
    avg_first_grpc_times.append(med(i[::2]))
    avg_second_grpc_times.append(med(i[1::2]))


print avg_first_grpc_times
print avg_second_grpc_times

avg_first_thrift_times = []
avg_second_thrift_times = []
for i in thrift_times:
    avg_first_thrift_times.append(med(i[::2]))
    avg_second_thrift_times.append(med(i[1::2]))


print avg_first_thrift_times
print avg_second_thrift_times


fg = plt.scatter(range(31), avg_first_grpc_times, color='blue')
sg = plt.scatter(range(31), avg_second_grpc_times, color='red')

ft = plt.scatter(range(31), avg_first_thrift_times, color='orange')
st = plt.scatter(range(31), avg_second_thrift_times, color='green')

plt.xlim((-1, 45))

plt.title('Grpc vs Thrift int marshalling')

plt.xlabel('log2(Argument value)')
plt.ylabel('Time (ns)')

plt.legend([fg, sg, ft, st], ['grpc 1st', 'grpc 2nd', 'thrift 1st', 'thrift 2nd'])
plt.savefig("grpc_vs_thrift_int.png")

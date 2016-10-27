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

opt_times = []
for i in range(31):
    opt_times.append([])
    with open("opt/grpc/int_times/" + str(i) + ".txt") as f:
        opt_times[i] = [int(t) for t in f.read().strip().split('\n')]


avg_first_grpc_times = []
avg_second_grpc_times = []
for i in grpc_times:
    avg_first_grpc_times.append(med(i[::2]))
    avg_second_grpc_times.append(med(i[1::2]))


print avg_first_grpc_times
print avg_second_grpc_times

avg_first_opt_times = []
avg_second_opt_times = []
for i in opt_times:
    avg_first_opt_times.append(med(i[::2]))
    avg_second_opt_times.append(med(i[1::2]))


print avg_first_opt_times
print avg_second_opt_times


fn = plt.scatter(range(31), avg_first_grpc_times, color='blue')
sn = plt.scatter(range(31), avg_second_grpc_times, color='red')

fo = plt.scatter(range(31), avg_first_opt_times, color='orange')
so = plt.scatter(range(31), avg_second_opt_times, color='green')

plt.xlim((-1, 45))

plt.title('gRPC Optimized vs Non-Optimized int marshalling')

plt.xlabel('log2(Argument value)')
plt.ylabel('Time (ns)')

plt.legend([fo, so, fn, sn], ['1st opt', '2nd opt', '1st nonopt', '2nd nonopt'])
plt.savefig("opt_vs_nonopt_int.png")

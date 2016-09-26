#!/usr/bin/python
import matplotlib.pyplot as plt

def mean(numbers):
    return float(sum(numbers)) / max(len(numbers), 1)

times = []
for i in range(32):
    times.append([])
    with open("int_times/" + str(i) + ".txt") as f:
        times[i] = [int(t) for t in f.read().strip().split('\n')]

avg_first_times = []
avg_second_times = []
for i in times:
    avg_first_times.append(mean(i[::2]))
    avg_second_times.append(mean(i[1::2]))


print avg_first_times
print avg_second_times

plt.scatter(range(32), avg_first_times)
plt.scatter(range(32), avg_second_times)
plt.show()

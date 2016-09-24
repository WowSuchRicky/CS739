#ifndef _BENCH_H_
#define _BENCH_H_

typedef unsigned long long u64;
typedef unsigned int uint;
typedef unsigned long long ull;


uint low;
uint high;

#define RDTSCP(result) \
        asm volatile("rdtscp" : "=a" (low), "=d" (high)); \
        result = ((ull)low) | ((ull)high << 32);

#endif

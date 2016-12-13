#include <sys/mman.h>
#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>

/*
 * We want a way to know heap addresses, from the beginning to the end at time of call.
 *
 * We can use sbrk(0) to see top of heap segment (heap grows up).
 * It actually returns the UPPER LIMIT on the heap at the time of call,
 * Future calls to malloc() should allocate between the top of the heap and that upper
 * limit returned by sbrk(0), meaning that this is only a valid way of knowing the start
 * and end address of the heap if malloc() has fully allocated up to the top of the current
 * break.
 * 
 * NOTE: this may not be for all allocations, and might be implementation dependent.
 * Good resource: http://stackoverflow.com/questions/6988487/what-does-brk-system-call-do
 *
 */

void printHeapValue() {
  printf("Heap address is: %d\n", sbrk(0));
}

void allocateInteger() { // int is 4 bytes here
  int *x = malloc(sizeof(int));
  printf("Allocated integer at address: %d\n", x);
}

int main() {
  printHeapValue();
  int i;
  for (i = 0; i < 100000; i++) {
    allocateInteger();
    printHeapValue();
  }
  return 0;
}


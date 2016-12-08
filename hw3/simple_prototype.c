#include <sys/mman.h>
#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include <errno.h>
#include <string.h>
#include <signal.h>

// addr given to mprotect must be at page boundary (mask lower bits)
#define GET_PAGE_BOUNDARY(addr, pgsize) (addr & ~(pgsize-1))

/*
 * Overall task:
 * - take note of beginning of heap, allocate some stuff in parent function, take note of
 *   end of heap (so we know what we need to protect)
 * - have wrapper function that is called; this wrapper protects the heap, calls into the
 *   child function, and then unprotects the memory (after child function is complete)
 * - have child function that tries to access some variable in the parent, and make sure it
 *   seg faults
 * - next step is then to implement copy on write functionality, which should involve:
 *    1) allocating segment of memory (side pool) that is the size of the parent's heap
 *    2) change the SIGSEGV handler to, on fault in parent's heap range, fulfill the 
 *       the request in that side pool and return what's at that location
 */

void *getHeapBound() {
  return malloc(4 * sizeof(int));
}

int *allocateInteger() { // int is 4 bytes here
  int *x = (int *)malloc(sizeof(int));
  *x = 4;
  printf("Allocated integer at address: %d\n", x);
  return x;
}

void protectAndCall(void *heapBegin, void *heapEnd, int *intPtr, void (*functionPtr)(int *)) {
  int s;

  heapBegin = (void *)GET_PAGE_BOUNDARY((long int)heapBegin, sysconf(_SC_PAGESIZE));
  int heapSize = (int)(heapEnd - heapBegin);
  
  // allocate additional memory to fill up from heapBegin + heapSize to end of that page
  int addrLastPage = GET_PAGE_BOUNDARY((long int)heapBegin + heapSize, sysconf(_SC_PAGESIZE));
  int extraBytesToAllocate = addrLastPage + sysconf(_SC_PAGESIZE) - (long int)heapEnd;
  void *extraAllocation = malloc(extraBytesToAllocate);
  
  // make the memory read only
  s = mprotect(heapBegin, heapSize, PROT_NONE);
  s = mprotect(heapBegin, heapSize, PROT_READ);
  if (s != 0) {
    printf("Protecting memory failed, status %d, errno: %s\n", s, strerror(errno));
    exit(-1);
  }

  //  call our 'lambda'
  return (*functionPtr)(intPtr);

  // make memory available for anything again
  s = mprotect(heapBegin, heapSize, PROT_READ|PROT_WRITE|PROT_EXEC);
  if (s != 0) {
    printf("Unprotecting memory failed, status %d, errno: %s\n", s, strerror(errno));
    exit(-1);
  }

  free(extraAllocation);

}

/*
 * test functions that will be executed after being protected
 */ 
void writeNumber(int *num) {
  *num = 3;
  printf("Number value is %d\n", *num);
}

void readNumber(int *num) {
  int x = *num;
  printf("Number value is %d\n", *num);
}

void allocateStuff(int *existingNum) { 
  int *newNum = (int *)malloc(sizeof(int));
  *newNum = 3 + *existingNum;
}


/* Begin sig handler stuff */
static void
handler(int sig, siginfo_t *si, void *unused) {
  fprintf(stderr, "Got SIGSEGV at address: 0x%lx\n",
    (long) si->si_addr);
    fflush(stdout);
    exit(EXIT_FAILURE);
}


int main() {
  // set up sig handler:
   struct sigaction sa;
   sa.sa_flags = SA_SIGINFO;
   sigemptyset(&sa.sa_mask);
   sa.sa_sigaction = handler;
   if (sigaction(SIGSEGV, &sa, NULL) == -1)
     exit(1);
  
  void *startHeap = getHeapBound();
  int *intPtr = allocateInteger(); // a protected value in parent
  void *endHeap = getHeapBound();
  
  printf("Heap start and end: %d, %d\n", startHeap, endHeap);
  printf("Size of heap in bytes: %d\n", (endHeap - startHeap) / 8);

  protectAndCall(startHeap, endHeap, intPtr, readNumber);
  printf("SUCCESS: reading number in protected parent\n");

  protectAndCall(startHeap, endHeap, intPtr, allocateStuff);
  printf("SUCCESS: allocating and writing to new memory in child\n");

  protectAndCall(startHeap, endHeap, intPtr, writeNumber);
  printf("SUCCESS: child writing to value in protected parent (this SHOULD fail!)\n");
  
  return 0;
}


#include <stdio.h>
#include <pthread.h>
#include "bench.h"

pthread_t tid;
pthread_mutex_t lock;
ull start, end;

void *doNothing(void *arg){
  pthread_mutex_init(&lock, NULL);
  
  RDTSCP(start);
  pthread_mutex_lock(&lock);
  RDTSCP(end);

  printf("%llu\n", end-start);
  
}

int main(int argc, char *argv[]){
  int err;
  
  err = pthread_create(&tid, NULL, &doNothing, NULL);

  pthread_join(tid, NULL);
  
}

int main(){
  int size; size = 30;
  int array[30]; int mem[30];

  setInputArray(array, size);

  printArray(array, size);

  qsort(array, 0, size, mem);

  printArray(array, size);

  return 0;
}

int qsort(int* array, int head, int tail, int* mem){
  int size; size = tail - head;
  if(size < 2){ return 0; }

  int i;
  for(i=0; i<size; i=i+1){ mem[head+i]=array[head+i]; }

  int cmp; cmp = array[head];
  int leftTail; leftTail = head;
  int rightHead; rightHead = tail;

  for(i=1; i<size; i=i+1){
    int val; val = mem[head+i];
    if(val<cmp) {
      array[leftTail] = val;
      leftTail = leftTail+1;
    }
    if(cmp<val+1) {
      array[rightHead-1] = val;
      rightHead = rightHead-1;
    }
  }

  array[leftTail] = cmp;
  qsort(array, head, leftTail, mem);
  qsort(array, rightHead, tail, mem);

  return 0;
}

int setInputArray(int* array, int size){
  int x; x = 1; int i;
  for(i = 0; i<size; i=i+1){
    x = x * 20021 + 1; x = x - (x/65536)*65536;
    array[i] = x - (x/size)*size;
  }
  return 0;
}

int printArray(int* array, int size){
  int i; int v;
  for(i=0;i<size;i=i+1){
    v = *(array + i);
    printInt(v); put(32);
  }
  put(10);
  return 0;
}

int printInt(int n){
  int div; int rem;
  div = n / 10; rem = n - div*10;
  if(div != 0){ printInt(div); }
  put(48 + rem);
  return 0;
}

int put(int c) {
  syscall 1 1 &c 1;
  return 0;
}


format:
.string "value: 0x%x  %lld\n"


.globl main
main:
    push    %rbp
    mov     %rsp, %rbp


    # call RDTSC to get cycle counter val, prepare result 
    mov     $0x0, %rax
    rdtscp	
    shl $0x2, %rdx
    or  %rdx, %rax

    #save RDTSC val for comparison
    mov %rax, %rdx
	
    # printf the RDTSC val, (rax | rdx<<32)
    mov $format, %rdi
    mov %rax, %rsi
    call printf
    
    	
    # call RDTSC again (turn to function later)
    mov $0x0, %rax
    rdtscp
    shl $0x2, %rdx
    or %rdx, %rax

    # get difference of RDTSC vals start and end
    sub %rdx, %rax 

    # print difference
	mov $format, %rdi
	mov %rax, %rsi
	call printf
   

    # clean, return
    mov     $0x00, %rax
    mov     %rbp, %rsp
    pop     %rbp
    ret


format:
.string "value: 0x%x  %d\n"

.globl main
main:
    push    %rbp
    mov     %rsp, %rbp

    mov     $0x2353254, %rax
	
    
    mov     format, %rdi
    mov     %rax, %rsi
    mov     %rax, %rdx
    call    printf

    mov     $0x00, %rax
    mov     %rbp, %rsp
    pop     %rbp
    ret


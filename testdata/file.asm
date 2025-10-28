main:
    # Inicializa registradores
    addi x1, x0, 5
    addi x2, x0, 3
    addi x3, x0, 7
    addi x4, x0, 2
    addi x5, x0, 4
    addi x6, x0, 6

    # Operações ALU
    add  x7, x1, x2
    sub  x8, x3, x1
    add  x9, x4, x5
    sub  x10, x6, x2
    addi x11, x7, 10
    addi x12, x8, 5
    add  x13, x11, x12
    sub  x14, x13, x9
    addi x15, x14, 3
    addi x16, x15, 2
    add  x17, x16, x10
    sub  x18, x17, x1

    # Branches
    beq  x1, x2, branch1
    addi x19, x0, 1
branch1:
    addi x20, x0, 2
    bne  x3, x4, branch2
    addi x21, x0, 3
branch2:
    addi x22, x0, 4
    beq  x5, x6, branch3
    addi x23, x0, 5
branch3:
    addi x24, x0, 6

    # Mais ALU
    add  x25, x20, x22
    sub  x26, x24, x21
    addi x27, x25, 7
    addi x28, x26, 8
    add  x29, x27, x28
    sub  x30, x29, x19

    # Jump
    jal  x0, do_after_jump
    addi x31, x0, 99   # não executa

do_after_jump:
    addi x1, x0, 0

    # Loop final
end:
    nop
    j end
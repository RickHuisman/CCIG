section .text
global _main
_main
  push 5
  push 10
  pop rdi
  pop rax
  cqo
  idiv rdi
  push rax
  pop rax
  ret
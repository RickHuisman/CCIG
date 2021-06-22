all:
	nasm -f macho64 hello.asm
	gcc -arch x86_64 -o hello hello.o

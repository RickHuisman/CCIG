all:
	go run .
	nasm -f macho64 temp.asm
	gcc -arch x86_64 -o temp temp.o

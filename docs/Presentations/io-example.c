#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <unistd.h>
#include <stdio.h>

#define SIZE_TO_READ 100

char buffer[SIZE_TO_READ];

int main()
{
	int file = open("io-example.c", O_RDONLY, 0);
	
	if (file == -1)
	{
		printf("Unable to open file!\n");
		return 1;
	}
    
    // Kernel is free to deschedule us here.
    if(read(file, buffer, SIZE_TO_READ) != -1)
        printf("Success!\n");
    else
        printf("Error!\n");
}

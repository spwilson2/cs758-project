#include <aio.h>
#include <sys/types.h>
#include <fcntl.h>
#include <errno.h>
#include <unistd.h>
#include <string.h>
#include <stdio.h>

#define SIZE_TO_READ 100

char buffer[SIZE_TO_READ];

int main()
{
	int file = open("aio-example.c", O_RDONLY, 0);
	
	if (file == -1)
	{
		printf("Unable to open file!\n");
		return 1;
	}
	
    // Create and fill out AIO control buffer struct
	struct aiocb cb;

	memset(&cb, 0, sizeof(struct aiocb));
	cb.aio_nbytes = SIZE_TO_READ;
	cb.aio_fildes = file;
	cb.aio_offset = 0;
	cb.aio_buf = buffer;
	
	// read!
	if (aio_read(&cb) == -1)
	{
		printf("Unable to create request!\n");
		close(file);
	}
	
	// spin until the request finished
	while(aio_error(&cb) == EINPROGRESS)
	{
		printf("Working...\n");
	}
	
	// success?
	if (aio_return(&cb) != -1)
		printf("Success!\n");
	else
		printf("Error!\n");
}

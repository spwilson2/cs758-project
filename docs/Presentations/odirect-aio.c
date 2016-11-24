#define _GNU_SOURCE
#include <libaio.h>
#include <sys/types.h>
#include <fcntl.h>
#include <errno.h>
#include <unistd.h>
#include <string.h>
#include <stdio.h>
#include <stdlib.h>
#include <errno.h>

#define SIZE_TO_READ 10000
#define ROUNDED_READ (((SIZE_TO_READ/BLKSIZE)*BLKSIZE) + BLKSIZE)
#define BLKSIZE 512

int main()
{
    int file = open("odirect-aio.c", O_RDONLY|O_DIRECT, 0);

    if (file <= 0) {
        printf("%s\n", strerror(errno));
        printf("can't open file \n");
        return 1;
    }

    char *buffer = aligned_alloc(BLKSIZE, ROUNDED_READ);

    if (buffer <= 0) { 
        printf("%s\n", strerror(errno));
        printf("aligned_alloc error\n");
        return 1;
    }

    // Create io_context for the kernel
    io_context_t ctx;
    memset(&ctx, 0, sizeof(ctx));

    int retval = io_setup(ROUNDED_READ, &ctx);

    // Init with max size of request
    if (retval) { 
        printf("%s\n", strerror(-retval));
        printf("io_setup error\n");
        return 1;
    }

    // Allocate a struct containing our AIO request info
    struct iocb *iocb_request;
    iocb_request = calloc(1, sizeof(struct iocb));

    // Initialize the request information.
    io_prep_pread(iocb_request, file, buffer, ROUNDED_READ, 0);

    // Submit our request.
    retval = io_submit(ctx, 1, &iocb_request);
    if (retval != 1) {
        io_destroy(ctx);
        printf("%s\n", strerror(-retval));
        printf("io_submit errorn");
        return 1;
    }

    struct io_event e;
    struct timespec timeout;

    // Poll waiting for a response to our AIO request
    while (1) {
        timeout.tv_nsec=500000000;//0.5s

        int poll_val = io_getevents(ctx, 1, 1, &e, &timeout);

        if (poll_val == 1) {
            printf("%s\n", buffer);
            printf("Success!\n");
            break;
        } else if (poll_val < 0) {
            printf("Error!\n");
            break;
        } else
            sleep(1);
    } 
    io_destroy(ctx);
}

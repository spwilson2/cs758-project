#include <libaio.h>
#include <sys/types.h>
#include <fcntl.h>
#include <errno.h>
#include <unistd.h>
#include <string.h>
#include <stdio.h>
#include <stdlib.h>

#define SIZE_TO_READ 100

char buffer[SIZE_TO_READ];

int main()
{
    int file = open("aio-example.c", O_RDONLY, 0);

    // Create io_context for the kernel
    io_context_t ctx;
    memset(&ctx, 0, sizeof(ctx));

    // Init with max size of request
    if (io_setup(SIZE_TO_READ, &ctx)) { 
        printf("io_setup errorn");
        return 1;
    }

    // Allocate a struct containing our AIO request info
    struct iocb *iocb_request;
    iocb_request = calloc(1, sizeof(struct iocb));

    // Initialize the request information.
    io_prep_pread(iocb_request, file, &buffer, SIZE_TO_READ, 0);

    // Submit our request.
    if (io_submit(ctx, 1, &iocb_request) !=1 ) {
        io_destroy(ctx);
        printf("io_submit errorn");
        return 1;
    }

    struct io_event e;
    struct timespec timeout;

    // Poll waiting for a response to our AIO request
    while (1) {
        timeout.tv_nsec=500000000;//0.5s

        int poll_val = io_getevents(ctx, 0, 1, &e, &timeout);

        if (poll_val == 1) {
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

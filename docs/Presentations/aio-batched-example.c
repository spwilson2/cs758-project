#define _GNU_SOURCE
#include <libaio.h>
#include <sys/types.h>
#include <fcntl.h>
#include <errno.h>
#include <unistd.h>
#include <string.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <sys/syscall.h>

#define SIZE_TO_READ 100
#define CONCURRENT 10

char buffer[SIZE_TO_READ];

int main()
{
    int i = 0;
    int file = open("aio-batched-example.c", O_RDONLY, 0);

    // Create io_context for the kernel
    io_context_t ctx;
    memset(&ctx, 0, sizeof(ctx));

    // Init with max size of request
    if (io_setup(CONCURRENT, &ctx)) { 
        printf("io_setup errorn");
        return 1;
    }

    // Allocate a struct containing our AIO request info
    struct iocb *iocb_request;
    struct iocb **iocb_requestp;
    iocb_request = calloc(CONCURRENT, sizeof(struct iocb));
    iocb_requestp = malloc(CONCURRENT * sizeof(struct iocb*));

    for (i=0;i< CONCURRENT; i++)
        iocb_requestp[i] = &iocb_request[i];

    char **buffers = malloc(sizeof(char*)*CONCURRENT);

    for (i=0;i< CONCURRENT; i++)
        buffers[i] = calloc(SIZE_TO_READ, sizeof(char));

    // Initialize the request information.
    for (i=0;i< CONCURRENT; i++){
        io_prep_pread(&iocb_request[i], file, buffers[i], SIZE_TO_READ, i);
    }

    // Submit all at once.
    for (i=0; i<CONCURRENT; i++){
        if (io_submit(ctx, 1, iocb_requestp) !=1) {
            io_destroy(ctx);
            printf("Error val: %s\n", strerror(errno));
            printf("io_submit errorn\n");
            return 1;
        }
    }

    struct io_event *e = calloc(CONCURRENT,sizeof(struct io_event));
    struct timespec timeout;
    memset(&timeout, 0,sizeof(timeout));

    // Poll waiting for a response to our AIO request
    while (1) {
        timeout.tv_nsec=500000000;//0.5s

        int poll_val = io_getevents(ctx, 1, CONCURRENT, e, &timeout);

        if (poll_val >= 1) {
            printf("poll_val:%d\n", poll_val);
            printf("Success!\n");
            break;
        } else if (poll_val < 0) {
            printf("Error!\n");
            break;
        } else {
            printf("poll_val:%d\n", poll_val);
            printf("Sleeping.\n");
            //sleep(1);
        }
    } 
    printf("%s\n", buffers[9]);
    io_destroy(ctx);
    return 0;
}

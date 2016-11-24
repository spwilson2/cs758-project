package aio

import "syscall"
import "unsafe"

const (
	CMD_PREAD  = syscall.IOCB_CMD_PREAD
	CMD_PWRITE = syscall.IOCB_CMD_PWRITE
	CMD_FSYNC  = syscall.IOCB_CMD_FSYNC
	CMD_FDSYNC = syscall.IOCB_CMD_FDSYNC
	/*
	* These two are experimental.
	* CMD_PREAD  = syscall.IOCB_CMD_PREAD
	* CMD_POLL   = syscall.IOCB_CMD_POLL
	 */
	CMD_NOOP    = syscall.IOCB_CMD_NOOP
	CMD_PREADV  = syscall.IOCB_CMD_PREADV
	CMD_PWRITEV = syscall.IOCB_CMD_PWRITEV
)

//func memzero(p unsafe.Pointer, size uintptr) {
//	b := *(*[]byte)(unsafe.Pointer(&p))
//	for i := uintptr(0); i < size; i++ {
//		b[i] = 0
//	}
//}

func PrepPread(iocb *syscall.Iocb, fd int, buf []byte, count int, offset int64) {

	// Clear out the iocb
	//memzero(unsafe.Pointer(iocb), unsafe.Sizeof(iocb))
	*iocb = syscall.Iocb{}

	iocb.Fildes = uint32(fd)
	iocb.Lio_opcode = CMD_PREAD
	iocb.Reqprio = 0
	iocb.Buf = uint64(uintptr(unsafe.Pointer(&buf[0])))
	iocb.Nbytes = uint64(count)
	iocb.Offset = offset

}

func PrepPwrite(iocb *syscall.Iocb, fd int, buf []byte, count uint64, offset int64) {
	// Clear out the iocb
	//memzero(unsafe.Pointer(iocb), unsafe.Sizeof(iocb))
	*iocb = syscall.Iocb{}

	iocb.Fildes = uint32(fd)
	iocb.Lio_opcode = CMD_PWRITE
	iocb.Reqprio = 0
	iocb.Buf = uint64(uintptr(unsafe.Pointer(&buf[0])))
	iocb.Nbytes = count
	iocb.Offset = offset
}

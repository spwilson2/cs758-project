package blocking

// blocking IO Scheduler

import "os"

/* extend os.File so we can actually make our own methods */
type File struct {
	os.File
}

/* opens file  */
func Open(name string) (*os.File, error) {
	return os.Open(name)
}

/* opens file w/ specified perms and flags*/
func OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

/* Writes p to file opened w/  fd */
func (f *File) Write(b []byte) (int, error) {
	return f.Write(b)
}

/* Writes to a specific byte of a file */
func (f *File) WriteAt(b []byte, off int64) (int, error) {
	return f.WriteAt(b, off)
}

/* Reads from file */
func (f *File) Read(b []byte) (int, error) {
	return f.Read(b)
}

/* Read from file, at offset off */
func (f *File) ReadAt(b []byte, off int64) (int, error) {
	return f.ReadAt(b, off)
}

/* Create file */
func Create(name string) (*os.File, error) {
	return os.Create(name)
}

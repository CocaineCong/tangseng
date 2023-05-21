package storage

import (
	"syscall"
)

func Mmap(fd int, offset int64, length int) ([]byte, error) {
	return syscall.Mmap(fd, offset, length, syscall.PROT_READ, syscall.MAP_SHARED)
}

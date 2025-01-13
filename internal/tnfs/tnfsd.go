package tnfs

/*
#cgo darwin CFLAGS: -DUNIX -DBSD -DENABLE_CHROOT -DNEED_ERRTABLE
#cgo windows CFLAGS: -DWIN32 -DNEED_ERRTABLE -DNEED_BSDCOMPAT
#cgo windows LDFLAGS: -lwsock32
#cgo linux CFLAGS: -DUNIX -DNEED_BSDCOMPAT -DENABLE_CHROOT -DNEED_ERRTABLE
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#ifdef WIN32
#include <winsock2.h>
#include <windows.h>
#endif
#include "../../tnfsd/src/auth.c"
#include "../../tnfsd/src/chroot.c"
#include "../../tnfsd/src/datagram.c"
#include "../../tnfsd/src/directory.c"
#include "../../tnfsd/src/endian.c"
#include "../../tnfsd/src/errortable.c"
#include "../../tnfsd/src/event_common.c"
#ifdef BSD
#include "../../tnfsd/src/event_kqueue.c"
#endif
#ifdef WIN32
#include "../../tnfsd/src/event_select.c"
#endif
#ifdef LINUX
#include "../../tnfsd/src/event_epoll.c"
#endif
#include "../../tnfsd/src/fileinfo.c"
#include "../../tnfsd/src/log.c"
#include "../../tnfsd/src/session.c"
#include "../../tnfsd/src/stats.c"
#include "../../tnfsd/src/strlcat.c"
#include "../../tnfsd/src/strlcpy.c"
#include "../../tnfsd/src/tnfs_file.c"
#include "../../tnfsd/src/tnfsd.c"
*/
import "C"
import (
	"os"
	"unsafe"
)

type InvalidDirError struct{}
type SocketError struct{}

func (e InvalidDirError) Error() string {
	return "invalid directory"
}
func (e SocketError) Error() string {
	return "couldn't open socket"
}

func TnfsdInit(log *os.File) {
	C.tnfsd_init_logs(C.int(log.Fd()))
	C.tnfsd_init()
}

func TnfsdStart(rootPath string, port int, read_only bool) error {
	path := C.CString(rootPath)
	defer C.free(unsafe.Pointer(path))

	err := C.tnfsd_start(path, C.int(port), C.bool(read_only))
	if err == -1 {
		return &InvalidDirError{}
	} else if err == -2 {
		return &SocketError{}
	}
	return nil
}

func TnfsdStop() {
	C.tnfsd_stop()
}

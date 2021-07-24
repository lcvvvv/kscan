package rdp

/*
//darwin编译环境配置
#cgo darwin CFLAGS: -DCGO_OS_DARWIN=1
#cgo darwin CFLAGS: -I./lib/mac/freerdp2/include/freerdp2
#cgo darwin CFLAGS: -I./lib/mac/freerdp2/include/winpr2
#cgo darwin LDFLAGS: ${SRCDIR}/lib/mac/freerdp2/lib/libcrypto.a
#cgo darwin LDFLAGS: ${SRCDIR}/lib/mac/freerdp2/lib/libssl.a
#cgo darwin LDFLAGS: ${SRCDIR}/lib/mac/freerdp2/lib/libfreerdp2.a
#cgo darwin LDFLAGS: ${SRCDIR}/lib/mac/freerdp2/lib/libwinpr2.a

//windows编译环境配置
#cgo windows CFLAGS: -DCGO_OS_WINDOWS=1
#cgo windows CFLAGS: -I./lib/windows/freerdp2/include/freerdp2
#cgo windows CFLAGS: -I./lib/windows/freerdp2/include/winpr2
#cgo windows LDFLAGS: ${SRCDIR}/lib/windows/freerdp2/lib/libcrypto.a
#cgo windows LDFLAGS: ${SRCDIR}/lib/windows/freerdp2/lib/libssl.a
#cgo windows LDFLAGS: ${SRCDIR}/lib/windows/freerdp2/lib/libfreerdp2.a
#cgo windows LDFLAGS: ${SRCDIR}/lib/windows/freerdp2/lib/libwinpr2.a

//linux编译环境配置
#cgo linux CFLAGS: -DCGO_OS_LINUX=1
#cgo linux CFLAGS: -I./lib/linux/freerdp2/include/freerdp2
#cgo linux CFLAGS: -I./lib/linux/freerdp2/include/winpr2
#cgo linux LDFLAGS: ${SRCDIR}/lib/linux/freerdp2/lib/libcrypto.a
#cgo linux LDFLAGS: ${SRCDIR}/lib/linux/freerdp2/lib/libssl.a
#cgo linux LDFLAGS: ${SRCDIR}/lib/linux/freerdp2/lib/libfreerdp2.a
#cgo linux LDFLAGS: ${SRCDIR}/lib/linux/freerdp2/lib/libwinpr2.a

#include <freerdp/freerdp.h>

BOOL rdp_connect(char *server, int port, char *domain, char *login, char *password) {
    int err;
    freerdp* instance;
    instance = freerdp_new();
    if (instance == NULL || freerdp_context_new(instance) == FALSE) {
        return -1;
    }
    instance->settings->Username = login;
    instance->settings->Password = password;
    instance->settings->IgnoreCertificate = TRUE;
    instance->settings->AuthenticationOnly = TRUE;
    instance->settings->ServerHostname = server;
    instance->settings->ServerPort = port;
    instance->settings->Domain = domain;
    freerdp_connect(instance);
    err = freerdp_get_last_error(instance->context);
    switch (err) {
        case 0:
        	freerdp_disconnect(instance);
			freerdp_free(instance);
            return err;
        case 0x00020009:
        case 0x00020014:
        case 0x00020015:
			// login failure
        case 0x0002000d:
        case 0x00020006:
        case 0x00020008:
        case 0x0002000c:
			freerdp_free(instance);
            // cannot establish rdp connection, either the port is not opened or it's
            // not rdp
			return err;
    }
	freerdp_free(instance);
    return err;
}

int check_rdp(char *ip, int port, char *domain, char *login, char *password) {
//int check_rdp() {
    int login_result = 0;
    wLog *root = WLog_GetRoot();
    WLog_SetStringLogLevel(root, "OFF");
    login_result = rdp_connect(ip, port, domain, login, password);
    return login_result;
}

*/
import "C"
import (
	"errors"
	"sync"
)

var mtx sync.Mutex

func Check(ip, domain, login, password string, port int) (bool, error) {
	mtx.Lock()
	defer mtx.Unlock()

	var nIp *C.char = C.CString(ip)
	var nDomain *C.char = C.CString(domain)
	var nLogin *C.char = C.CString(login)
	var nPassword *C.char = C.CString(password)
	var nPort C.int = C.int(port)

	//defer func() {
	//	C.free(unsafe.Pointer(nIp))
	//	C.free(unsafe.Pointer(nDomain))
	//	C.free(unsafe.Pointer(nLogin))
	//	C.free(unsafe.Pointer(nPassword))
	//}()

	rInt := int(C.check_rdp(nIp, nPort, nDomain, nLogin, nPassword))
	switch rInt {
	case 0:
		return true, nil
	case -1:
		return false, errors.New("freerdp init failed")
	default:
		return false, nil
	}
}

package rdp

/*
#cgo CFLAGS: -I./freerdp2/include/freerdp2
#cgo CFLAGS: -I./freerdp2/include/winpr2
#cgo LDFLAGS: -L./freerdp2/lib



#include <freerdp/freerdp.h>

BOOL rdp_connect(char *server, int32_t port, char *domain, char *login, char *password) {
    int32_t err;
    freerdp *instance;
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
    if (err == 0) {
        freerdp_disconnect(instance);
        freerdp_free(instance);
    }
    return err;
}

extern int32_t check_rdp(char *ip, int32_t port, char *domain, char *login, char *password) {
//int32_t check_rdp() {
    int32_t login_result = 0;
    wLog *root = WLog_GetRoot();
    WLog_SetStringLogLevel(root, "OFF");
    login_result = rdp_connect(ip, port, domain, login, password);
    switch (login_result) {
        case 0:
            break;
        case 0x00020009:
        case 0x00020014:
        case 0x00020015:
            // login failure
            break;
        case 0x0002000d:
            break;
        case 0x00020006:
        case 0x00020008:
        case 0x0002000c:
            // cannot establish rdp connection, either the port is not opened or it's
            // not rdp
            break;
    }
    return login_result;
}
*/
import "C"
import "errors"

func Check(ip, domain, login, password string, port int) (bool, error) {
	nip := C.CString(ip)
	ndomain := C.CString(domain)
	nlogin := C.CString(login)
	npassword := C.CString(password)
	nport := C.int32_t(port)
	switch int(C.check_rdp(nip, nport, ndomain, nlogin, npassword)) {
	case 0:
		return true, nil
	case -1:
		return false, errors.New("freerdp init failed")
	default:
		return false, nil
	}
}

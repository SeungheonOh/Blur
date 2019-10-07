package main

/*#cgo LDFLAGS: -lxcb
#cgo CFLAGS: -w
#include <xcb/xcb.h>
#include <xcb/xproto.h>
#include <string.h>
#include <stdlib.h>

xcb_atom_t GetAtom(xcb_connection_t *conn, char *name){
	xcb_atom_t atom;
	xcb_intern_atom_cookie_t cookie;

	cookie = xcb_intern_atom(conn, 0, strlen(name), name);

	xcb_intern_atom_reply_t *reply = xcb_intern_atom_reply(conn, cookie, NULL);
	if(reply) {
		atom = reply->atom;
		free(reply);
	}
	return atom;
}

int GetWindows(xcb_connection_t *conn) {
	xcb_atom_t atom = GetAtom(conn, "_NET_CLIENT_LIST");
	xcb_get_property_cookie_t prop_cookie;
	prop_cookie = xcb_get_property(conn, 0, xcb_setup_roots_iterator(xcb_get_setup(conn)).data->root, atom, 0, 0, (1 << 32)-1);
	xcb_get_property_reply_t *prop_reply;
	prop_reply = xcb_get_property_reply(conn, prop_cookie, NULL);

	return prop_reply->value_len;
}
*/
import "C"
import (
	"fmt"
)

func main() {
	conn := C.xcb_connect(nil, nil)
	defer C.xcb_disconnect(conn)
	fmt.Print("\033[s")
	for {
		fmt.Print("\033[u")
		fmt.Print(C.GetWindows(conn))
	}
}

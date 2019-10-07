#include <xcb/xcb.h>
#include <xcb/xproto.h>
#include <unistd.h>
#include <stdio.h>
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

int main (){

	xcb_connection_t *c;
	xcb_screen_t *screen;

	c = xcb_connect(NULL, NULL);
	screen = xcb_setup_roots_iterator(xcb_get_setup(c)).data;

	xcb_atom_t atom = GetAtom(c, "_NET_CLIENT_LIST");

	xcb_get_property_cookie_t prop_cookie;
	prop_cookie = xcb_get_property(c, 0, screen->root, atom, 0, 0, (1 << 32)-1);
	xcb_get_property_reply_t *prop_reply;
	prop_reply = xcb_get_property_reply(c, prop_cookie, NULL);

	printf("len : %d\n", prop_reply->value_len);
	printf("type : %d\n", prop_reply->type);
	uint32_t *values = (uint32_t*)xcb_get_property_value(prop_reply);
	printf("data size: %d\n", xcb_get_property_value_length(prop_reply));

	for(int i = 0; i < prop_reply->value_len; i++){
		printf("%d : %u\n", i, values[i]);
	}

	return 0;
}

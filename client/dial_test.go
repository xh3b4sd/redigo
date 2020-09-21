package client

import (
	"net"
	"testing"
)

func Test_Client_NewDial_Error_Addr(t *testing.T) {
	c := DefaultDialConfig()
	c.Address = "foo"

	d := NewDial(c)

	_, err := d()
	if e, ok := err.(*net.OpError); ok {
		if e.Op != "dial" {
			t.Fatal("expected", "dial", "got", e.Op)
		}
	} else {
		t.Fatal("expected", "*net.OpError", "got", err)
	}
}

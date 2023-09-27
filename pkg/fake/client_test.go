package fake

import (
	"testing"

	"github.com/xh3b4sd/redigo"
)

func Test_Fake_Interface(t *testing.T) {
	var _ redigo.Interface = New()
}

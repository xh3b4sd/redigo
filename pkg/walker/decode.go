package walker

import (
	"strconv"

	"github.com/xh3b4sd/tracer"
)

func decode(res []interface{}) (int64, []string, error) {
	cur, err := strconv.ParseInt(string(res[0].([]uint8)), 10, 64)
	if err != nil {
		return 0, nil, tracer.Mask(err)
	}

	var str []string
	for _, v := range res[1].([]interface{}) {
		str = append(str, string(v.([]uint8)))
	}

	return cur, str, nil
}

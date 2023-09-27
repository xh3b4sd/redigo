package backup

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"
)

func (b *Backup) Create() error {
	var con redis.Conn
	{
		con = b.poo.Get()
		defer con.Close()
	}

	var pre time.Time
	{
		uni, err := redis.Int64(con.Do("LASTSAVE"))
		if err != nil {
			return tracer.Mask(err)
		}

		pre = time.Unix(uni, 0)
	}

	{
		_, err := redis.String(con.Do("BGSAVE"))
		if err != nil {
			return tracer.Mask(err)
		}
	}

	{
		for {
			var cur time.Time
			{
				uni, err := redis.Int64(con.Do("LASTSAVE"))
				if err != nil {
					return tracer.Mask(err)
				}

				cur = time.Unix(uni, 0)
			}

			{
				if cur.After(pre) {
					break
				}
			}

			{
				time.Sleep(1 * time.Second)
			}
		}
	}

	return nil
}

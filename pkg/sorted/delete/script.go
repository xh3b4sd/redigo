package delete

const deleteCleanScript = `
	redis.call("DEL", KEYS[2])
	redis.call("DEL", KEYS[1])

	return 0
`

const deleteIndexScript = `
	for i = 1, #ARGV do
			local sco = redis.call("ZSCORE", KEYS[1], ARGV[i])

			if (sco ~= false) then
					redis.call("ZREMRANGEBYSCORE", KEYS[2], sco, sco)
			end
	end

	redis.call("ZREM", KEYS[1], unpack(ARGV))

	return 0
`

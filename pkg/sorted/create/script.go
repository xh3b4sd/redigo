package create

const createIndexScript = `
	-- Verify if the score does already exist. The first key here might be "ssk"
	-- and the second argument might be "0.8". If we get any value in response
	-- the score is already taken.
	local res = redis.call("ZRANGE", KEYS[1], ARGV[2], ARGV[2], "BYSCORE")
	if (res[1] ~= nil) then
		return 0
	end

	if (ARGV[3] ~= nil) then
		-- Verify if the index does already exist. The second key here might be
		-- "ssk:ind" and the argument might be "name". If we get any value in
		-- response the index is already taken.
		local i = 3
		while ARGV[i] do
			local res = redis.call("ZSCORE", KEYS[2], ARGV[i])
			if (res ~= false) then
				return 1
			end

			i=i+1
		end

		-- Only if we ensured that the score is unique and that all indizes are not
		-- yet recorded, we can then add them to our sorted sets. Note that all
		-- indices for an element are recorded with the same score.
		local j = 3
		while ARGV[j] do
			redis.call("ZADD", KEYS[2], ARGV[2], ARGV[j])

			j=j+1
		end
	end

	redis.call("ZADD", KEYS[1], ARGV[2], ARGV[1])

	return 2
`

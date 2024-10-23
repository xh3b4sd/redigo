package update

const updateIndexScript = `
	-- Verify if the sorted set does even exist. We must not proceed if we would
	-- create a new sorted set and a new element within it.
	local exi = redis.call("EXISTS", KEYS[1])
	if (exi == 0) then
		return 0
	end

	local function upd(key, new, sco)
		-- We actually verified the existence of the element already. Now we
		-- only fetch the old value in order to perform the update.
		local res = redis.call("ZRANGE", key, sco, sco, "BYSCORE")
		local old = res[1]

		-- If the value did not change it might mean the indices did not change.
		-- We are ok with that internally. It is only important that the user
		-- facing element is properly reported to be updated or not.
		if (old == new) then
			return 2
		end

		redis.call("ZADD", key, sco, new)
		redis.call("ZREM", key, old)

		return 3
	end

	local function ver(key, new, sco)
		-- Verify if the score does already exist. If there is no element we
		-- cannot update it.
		local res = redis.call("ZRANGE", key, sco, sco, "BYSCORE")
		local old = res[1]
		if (old == nil) then
			return 1
		end

		-- Verify if the existing value is already what we want to update to. If
		-- the desired state is already reconciled we do not need to proceed
		-- further.
		if (old == new) then
			return 2
		end

		return 3
	end

	-- Verify all scores have associated values. We need to do this upfront for
	-- the given element and the internally managed indices.
	local i = 3
	while ARGV[i] do
		local res = ver(KEYS[2], ARGV[i], ARGV[2])
		if (res == 1) then
			return res
		end

		i=i+1
	end
	local res = ver(KEYS[1], ARGV[1], ARGV[2])
	if (res == 1) then
		return res
	end

	-- Only if all verifications are completed successfully and there is no
	-- reason to fail anymore we can continue to actually process the updates.
	local j = 3
	while ARGV[j] do
		upd(KEYS[2], ARGV[j], ARGV[2])

		j=j+1
	end

	return upd(KEYS[1], ARGV[1], ARGV[2])
`

const updateScoreScript = `
	-- Verify if the score does already exist. If there is no element we
	-- cannot update it.
	local res = redis.call("ZRANGE", KEYS[1], ARGV[2], ARGV[2], "BYSCORE")

	local old = res[1]
	if (old == nil) then
		return 0
	end

	-- Verify if the existing value is already what we want to update to. If
	-- the desired state is already reconciled we do not need to proceed
	-- further.
	if (old == ARGV[1]) then
		return 1
	end

	redis.call("ZADD", KEYS[1], ARGV[2], ARGV[1])
	redis.call("ZREM", KEYS[1], old)

	return 2
`

# redigo

Redis client abstraction for simple use.



### Conformance Tests

```
docker run --rm -p 127.0.0.1:6379:6379 redis
```

```
docker run --rm -e REDIS_MASTER_HOST=127.0.0.1 -p 127.0.0.1:26379:26379 bitnami/redis-sentinel
```

```
go test ./... --tags single,sentinel
```



### Lua Debugging

```
redis.log(redis.LOG_NOTICE, type(v))
redis.log(redis.LOG_NOTICE, cjson.encode(v))
```

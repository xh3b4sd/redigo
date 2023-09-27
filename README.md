# redigo

Opinionated redis client.



### Conformance Tests

```
docker run --rm --name redis-stack -p 6379:6379 -p 8001:8001 redis/redis-stack:latest
```

```
docker run --rm -e REDIS_MASTER_HOST=127.0.0.1 -p 127.0.0.1:26379:26379 bitnami/redis-sentinel
```

```
go test ./... --tags single,sentinel -count 1 -race
```



### Lua Debugging

```
redis.log(redis.LOG_NOTICE, type(v))
redis.log(redis.LOG_NOTICE, cjson.encode(v))
```



### Version Requirement

```
Redis version=6.2.0
```

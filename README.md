# redigo

Redis client abstraction for simple use.



### Conformance Tests

```
docker run --rm -p 127.0.0.1:6379:6379 redis
```

```
go test ./... --tags redis
```



### Lua Debugging

```
redis.log(redis.LOG_WARNING, type(v))
```

Rate Limiter
============

Rate limiting implementation that decorates a handler with limiting capabilities. Returns 429 should limiting rules are exceeded. The rules look like [limits.json](limits.json):
```
[
    {"Limit": 5, "Duration": "10s"},
    {"Limit": 7, "Duration": "20s"},
    {"Limit": 9, "Duration": "30s"}
]
```

Design
------
Package structure reflects components, together with tests:
* cmd/restapi is the REST entrypoint, exposing limit rated handlers
* pkg/model implements json serializable model
* pkg/db implements interface and in-memory database to store `hits`
* pkg/limiter implements the limiting logic
* pkg/utils adds slice processing functionality
* pkg/decorator adds decorator wrapper, which prematurely returns 429 should limit be breached
* docker/resapi provides a dockerfile to run the code as standalone

Build docker container
----------------------
```
docker build -f docker/restapi/Dockerfile -t goratelimiter.restapi .
```

To run docker container
-----------------------
```
docker run --rm -ti -p8080:8080 goratelimiter.restapi
```

To run from cmd line
--------------------
```
go run cmd/restapi/restapi.go
```

To run tests
------------
```
# unit tests
go test ./... -v

# manual end-2-end tests, note, requires running the server either via cmd line or docker container
./manual-test.sh
```

Unesolved issues
----------------
* Current implementation of `leakybucket` utilizes slices, which leak memory as they grow. Attempt has been made to use `CircularQueue`, but that would drop awareness of rejected requests... Need some time based sliding window instead...
Rate Limiter
============

Rate limiting implementation that decorates a handler with limiting capabilities. Returns 429 should limiting rules are exceeded. The rules look like [limits.json](limits.json):
```
[
  {Limit: 5, Duration: 10},
	{Limit: 7, Duration: 20},
	{Limit: 9, Duration: 30}
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

To run manual tests
-------------------
```
# build the docker image and run the container, as above
./manual-test.sh
```
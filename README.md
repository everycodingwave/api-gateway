## api-gateway
Api-gateway is a http gateway server of autopilot api using redis as read through cache.

### Usage
```
cd <project root folder>
make
```
make command will downloading, building and running test of this project.
after running make, it will create an executable file api-gw. run it as:
```
export CACHE_ADDR=127.0.0.1:6379
./api-gw
```

Using env here as input for better containerizing this project although Dockerfile is not provided.
Another thing is cache testing code will only be run after setting this env.

### Further improvements
* more structured test code and more corner/error case testing
* redis client performance testing results

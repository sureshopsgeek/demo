### The Counter service
The counter service exposes two http endpoints namely counter and metrics through two different http servers. The counter endpoint accepts the status code and increment the internal counter for the  status code. The metrics endpoint exposes the counter for each status code which can be scaped by prometheus.

### Command-line arguments
The service binary accepts the following command line arguments
```
   -v  bool  set to enable verbose logging, default: false
   -cp int   counter server port, default: 8080
   -mp int   metric server port, default: 9201

```
### Starting the service
Following command will start the service in verbose mode and the servers will run on non default ports passed through the flags.
```
   counter -v -cp 8082 -mp 9203
```
### Sending status code to counter endpoint
The counter endpoint accepts the status code in a json body through PUT method.
``` 
   curl -XPUT http://demohost:8082/counter -d '{"code": 401 }'
```
### Checking the metrics endpoint
```
   curl -XGET http://demohost:9203/metrics
```
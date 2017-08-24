# fabric8-dependency-wait-service
A small go binary that waits for service dependencies to be up and running that can be used as a kubernetes init-container to simplify the startup of complex applications

# Environment variables
DEPENDENCY_LOG_VERBOSE: true/false, default=false. puts out a little more extra logs.  
DEPENDENCY_POLL_INTERVAL: a positive integer, default=1.  The interval between each poll of the dependency check.

# Usage
This utility is used as part of the init-containers.  An example is given below.
```
spec:
  template:
    metadata:
	  annotations:
          pod.beta.kubernetes.io/init-containers: |-
            [
            {
              "name": "init-dependencyservice1",
              "image": "fabric8io/fabric8-dependency-wait-service:withEnvVars",
              "imagePullPolicy": "IfNotPresent",
              "command": ["sh", "-c", "fabric8-dependency-wait-service-linux-amd64 postgres://wit@wit-db:5432"],
              "env": [{
                "name": "DEPENDENCY_POLL_INTERVAL",
                "value": "9"
                }, {
                "name": "DEPENDENCY_LOG_VERBOSE",
                "value": "true"
                }]
            },
            {
              "name": "init-dependencyservice2",
              "image": "fabric8io/fabric8-dependency-wait-service:withEnvVars",
              "imagePullPolicy": "IfNotPresent",
              "command": ["sh", "-c", "fabric8-dependency-wait-service-linux-amd64 http://keycloak:80"],
              "env": [{
                "name": "DEPENDENCY_POLL_INTERVAL",
                "value": "11"
                }, {
                "name": "DEPENDENCY_LOG_VERBOSE",
                "value": "true"
                }]
            },
```

* the protocols can be only http or postgres (as of now)
* http: for this a get request is made and if the response is 200 then it is ok.
* postgres: for this the pg_isready utility is used.  The given postgres url is parsed for username, host, and port and passed to the utility.  The db is considered up when a "accepting connections" response is received.

# witch
RESTful + MQTT process supervisor

# Install
```
go get github.com/Eagle-X/witch
```

# Usage
```
Usage of witch:
  -c string
       	Config file (default "witch.yaml")
```

# Config file
```
# Listen address, default: :5671.
listen: :5671
# Specify the process control system, available controls buildin, supervisor and systemd.
# Default: buildin
control: buildin
# Only if control is supervisor or systemd, service MUST be given.
service:
# Only if control is buidin, command MUST be given.
command: sleep 3600
# The pid file of the process to be supervised, MUST change different one.
# Only if control is buildin, pid_file MUST be given.
pid_file: witch.pid
# Connection authentication username and password,
# the format is {username: password, ...}. default: {noadmin: noADMIN}.
auth: {noadmin: noADMIN}
## Specify the MQTT Client configurations
mqtt:
  enable: true
  broker: tls://iot.eclipse.org:8883
  client_id: witch-1
  keepalive: 60
  username:
  password:
  actions_message:
    topic: 'go-mqtt/sample'
    qos: 1
```

# Exmaple
start witch
```
witch -c witch.ymal
```
To control with HTTP REST calls:
```
curl -u noadmin:noADMIN -XPUT -d '{"name":"is_alive"}' http://127.0.0.1:5671/api/app/actions
curl -u noadmin:noADMIN -XPUT -d '{"name":"start"}' http://127.0.0.1:5671/api/app/actions
curl -u noadmin:noADMIN -XPUT -d '{"name":"stop"}' http://127.0.0.1:5671/api/app/actions
curl -u noadmin:noADMIN -XPUT -d '{"name":"restart"}' http://127.0.0.1:5671/api/app/actions
```

To control with MQTT, publish actions message to the choosen topic:
```
* publish /topic -> {"name":"start"}
* publish /topic -> {"name":"stop"}
* publish /topic -> {"name":"restart"}
```



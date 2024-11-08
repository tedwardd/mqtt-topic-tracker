# MQTT Topic Tracker

Subscribes and listens to an MQTT broker and logs what and how many times topics are seen

# Build

## Local

Use this if you want to build it directly on your system

```
$ make build
```

## Container

Optionally, you can build inside a container environment

### Docker

```
$ make docker
```

### Podman

```
make podman
```

# Run

## Local

### No auth

```
$ build/mqtt-topic-tracker
```

### USER/PASS Required

There are two ways to pass username and password; command line arguments or environment variables. You can mix and match them however you like but command line args will always take precedence over environment variables.

#### Args
```
$ build/mqtt-topic-tracker --username=MY_USER --password=MY_PASS
```

#### Env Var
```
$ MQTT_USERNAME=MY_USER MQTT_PASSWORD=MY_PASS build/mqtt-topic-tracker
```

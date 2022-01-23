# mfrp

[![Build Status](https://travis-ci.org/fujiawei-dev/mfrp.svg)](https://travis-ci.org/fujiawei-dev/mfrp)

A mini fast reverse proxy. Just for study.

## Passive Mode

1. mfrpc has a public ip
2. mfrpc starts, then listens on a public port
3. client -> mfrpc -> server(nat) => client <-> mfrpc <-> server(nat)

## Active Mode

1. mfrps has a public ip
2. mfrps starts, then listens on a public port
3. mfrps waits for mfrpc's connection
4. mfrpc connects to mfrps
5. mfrpc > (control req) > mfrps (check: exist? password? idle?)
6. not idle > fail message > mfrps (greet with mfrpc conn in use), closes connection
7. idle > start the proxy server(listen on another public port -> working -> waits for user's connections) > success message
8. client > mfrps -> (work req) > mfrpc (start another work conn) > (work req) > idle?(bad req, close) working?(handle)

## Configuration

### mfrpc

```yaml
Common:
    ServerHost: localhost
    ServerPort: 9527
    LogLevel: debug
ProxyServers:
    - LocalPort: 8080
      PassiveMode: true
      BindAddr: 0.0.0.0
      BindPort: 18124
    - Name: mfrp
      Password: mfrp
      LocalPort: 8080
```

### mfrps

```yaml
Common:
  BindAddr: 0.0.0.0
  BindPort: 9527
  LogLevel: debug
ProxyServers:
  - Name: mfrp
    Password: mfrp
    BindAddr: 0.0.0.0
    ListenPort: 18123
```

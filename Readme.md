# GoRedis

GoRedis is supposed to be a project made up of two parts:

- a high performance, stable, and light-weight Redis server
- a fast, high-performance, and easy-to-use Go client library

## Quick Start

```bash
  ________      __________           .___.__        
 /  _____/  ____\______   \ ____   __| _/|__| ______
/   \  ___ /  _ \|       _// __ \ / __ | |  |/  ___/
\    \_\  (  <_> )    |   \  ___// /_/ | |  |\___ \ 
 \______  /\____/|____|_  /\___  >____ | |__/____  >
        \/              \/     \/     \/         \/ 
```

## Architecture

```bash
.
├─aof
├─cluster
├─config
├─database
├─data_struct
│  ├─bitmap
│  ├─dict
│  ├─list
│  ├─lock
│  ├─set
│  └─sorted_set
├─interface
│  ├─database
│  ├─redis
│  └─tcp
├─lib
│  ├─log
│  ├─sync
│  └─util
├─pubsub
├─redis
│  ├─client
│  ├─conn
│  ├─parser
│  ├─protocol
│  └─server
└─tcp
```

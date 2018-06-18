# SWS - Simple Web Server

By default SWS will serve files from the local directory to port 8080. A number of options are available to modify that behaviour if required.

Binary at: https://gitlab.com/bestmethod/go-sws/-/jobs

```
Usage:
  sws [OPTIONS] [directoryToServe]

Application Options:
  -b, --bind-address=    address to bind to. E.g. ':80' or '127.0.0.1:8000' (default: :8080)
  -l, --print-accesslog  print each file access log as it's requested
  -u, --auth-user=       if set, will enable basic HTTP auth, requiring this user
  -p, --auth-pass=       password to set for the basic HTTP auth user

Help Options:
  -h, --help             Show this help message
```

### If you need https:
###### Use our proxy project: https://gitlab.com/bestmethod/goproxy

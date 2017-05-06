# concourse-flake-hunter


```
concourse-flake-hunter 

  -c, --concourse-url= Concourse URL to authenticate with
  -u, --username=      Username for basic auth
  -p, --password=      Password for basic auth
  -n, --team-name=     Team to authenticate with (default: main)

```

##Example 

`concourse-flake-hunter -c https://my.concourse.com -u <username> -p <password> -n <team-name> search "connection reset"`


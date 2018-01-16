# concourse-flake-hunter


```
concourse-flake-hunter 

  -c, --concourse-url= Concourse URL to authenticate with
  -u, --username=      Username for basic auth
  -n, --team-name=     Team to authenticate with (default: main)

```

## Example 


The following command will search for "connection reset" in the last 150 builds run:

`CONCOURSE_BASIC_AUTH_PW=s3cr3t concourse-flake-hunter -c https://my.concourse.com -u <username> -n <team-name> search --limit 150 "connection reset"`


The outlook will look like the following:


```
+----------------------------+-----------------------------------+
|       PIPELINE/JOB         |            BUILD URL              |
+----------------------------+-----------------------------------+
|  product/unit-test         |  https://www.example.org/build/1  |
+----------------------------+-----------------------------------+
```

# concourse-flake-hunter


```
concourse-flake-hunter 

  -c, --concourse-url= Concourse URL to authenticate with
  -u, --username=      Username for basic auth
  -p, --password=      Password for basic auth
  -n, --team-name=     Team to authenticate with (default: main)

```

## Example 


The following command will search for "connection reset":

`concourse-flake-hunter -c https://my.concourse.com -u <username> -p <password> -n <team-name> search "connection reset"`
`CONCOURSE_BASIC_AUTH_PW=s3cr3t concourse-flake-hunter -c https://my.concourse.com -u <username> -n <team-name> search "connection reset"`

The outlook will look like the following:


```
+----------------------------+-----------------------------------+
|       PIPELINE/JOB         |            BUILD URL              |
+----------------------------+-----------------------------------+
|  product/unit-test         |  https://www.example.org/build/1  |
+----------------------------+-----------------------------------+
```


It's possible to aggregate all the pipelines by using `aggregate` command:

`concourse-flake-hunter -c https://my.concourse.com -u <username> -p <password> -n <team-name> aggregate "connection reset"`



The outlook will look like the following:

```
/tmp/build/8c72b58e/go/src/github.com/albertoleal/concourse-flake-hunter /tmp/build/8c72b58e
[Fail] [BeforeEach] My test 
	Count: 2
	LastOccurance: 2018-01-11 10:10:24 +0000 UTC
		JobName: test-upgrade-product-same-minor
		Date: 2018-01-10 09:43:13 +0000 UTC
		URL: https://my-awesome-concourse/teams/main/pipelines/my-pipeline/jobs/test-upgrade-product-same-minor/builds/36

		JobName: test-upgrade-product
		Date: 2018-01-10 05:16:12 +0000 UTC
		URL: https://my-awesome-concourse/teams/main/pipelines/my-pipeline/jobs/test-upgrade-product-previous-minor/builds/36

```

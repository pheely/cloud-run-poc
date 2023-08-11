# Traffic Management

## Usecase

To explore Cloud Run's support on revisions, rollback, gradual rollout, and traffic-split.

## Employee API's `/help` Endpoint

Employee Rest API has another GET endpoint `/help` which returns the version number of the application. For example, when accessed, the version 1 of the application will be returned the following information:
```bash
Employee API v1
```

And version 2 of the application will be returned 
```bash
Employee API v2
```

## Environment Variables

```bash
GOOGLE_CLOUD_PROJECT=ibcwe-event-layer-f3ccf6d9
REPOSITORY=us-central1-docker.pkg.dev/$GOOGLE_CLOUD_PROJECT/cloud-run-try
```

## Deploy the Service of Version 1

```bash
gcloud run deploy employee-api --image $REPOSITORY/employee:v1 --tag ver1 
```

Each new Cloud Run service revision can be assigned a tag. Doing this will create a revision specific URL which is unknown to the client. This can be useful to handle the traffic profile across multiple revisions. 

When deploying the service, you can specify the tag name using `--tag`. Here we use `ver1` to associate with version 1.

The revision specific URL follows the following convention:
```text
https://[tag-name]---[service-url]
```

For example, the URL of the employee api is 

https://employee-api-oy6beuif2a-uc.a.run.app, 

and the URL for revision 1 is 

https://ver1---employee-api-oy6beuif2a-uc.a.run.app. 

To test the service, run the following command:
```bash
VER1_URL=$(gcloud run services describe employee-api \
--format="value(status.address.url)"); echo $VER1_URL

curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
$VER1_URL/api/help 
```
The output should be:
```text
Employee API v1
```

## Deploy the Version 2 of the Service without Serving Traffic

We will deploy the version 2 of the service without assigning real traffic. We want to test it in production first. This can be done using the `--no-traffic` flag.

```bash
gcloud run deploy employee-api --image $REPOSITORY/employee:v2 --tag ver2 \
--no-traffic
```

The revision specific URL should be https://ver2---employee-api-oy6beuif2a-uc.a.run.app.

To verify revision 2 is deployed successful, run this command:
```bash
VER2_URL=$(gcloud run services describe employee-api \
--format="value(status.traffic[1].url)"); echo $VER2_URL

curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
$VER2_URL/api/help
```

The output should be as follows:
```text
Employee API v2 
```

## Gradual Rolling out Revison 2

Let's migrate 50% of the traffic to the revision tagged `ver2` using the `update-traffic` command:
```bash
gcloud run services update-traffic employee-api \
--to-tags ver2=50
```

Confirm that original service URL now is distributing traffic:
```bash
for i in {1..10}; \
do curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
$VER1_URL/api/help; \
done
```

The output should be similar to the following:
```text
Employee API v1 
Employee API v1 
Employee API v2 
Employee API v2 
Employee API v1 
Employee API v2 
Employee API v1 
Employee API v1 
Employee API v2 
Employee API v1 
```

## Rolling back a Tagged Version

In the event an issue is found, the traffic migration can be rolled back by resetting the percentage.

```bash
gcloud run services update-traffic employee-api \
  --to-tags ver2=0
```

Test the endpoint is distributing traffic to version 1 (tagged `ver1`) only:
```bash
for i in {1..10}; \
do curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
$VER1_URL/api/help; \
done
```

## Traffic Splitting

Let's deploy two more revisions of the service without routing traffic to them.

```bash
gcloud run deploy employee-api --image $REPOSITORY/employee:v3 \
--tag ver3 --no-traffic

gcloud run deploy employee-api --image $REPOSITORY/employee:v4 \
--tag ver4 --no-traffic
```

Create an environment variable for available revisions:
```bash
LIST=$(gcloud run services describe employee-api \
--format='value[delimiter="=25,"](status.traffic.revisionName)')"=25"; echo $LIST
```
```text
employee-api-00001-jex=25,employee-api-00002-nov=25,employee-api-00005-vur=25,employee-api-00006-net=25
```

Now split traffic among the four revisions using the environment variable:
```bash
gcloud run services update-traffic employee-api \
--to-revisions $LIST
```

Test the endpoint is distributing traffic:
```bash
for i in {1..10}; \
do curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
$VER1_URL/api/help; \
done
```

The output should be something like this:
```text
Employee API v1 
Employee API v2 
Employee API v4 
Employee API v2 
Employee API v3 
Employee API v2 
Employee API v4 
Employee API v2 
Employee API v4 
Employee API v1 
```

## Update Traffic to the Latest Version only

Reset the service traffic profile to use the latest deployment:
```bash
gcloud run services update-traffic employee-api --to-latest
```

Verify that the latest revision is able to receive traffic:
```bash
for i in {1..10}; \
do curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
$VER1_URL/api/help; \
done
```

<style>
    h1 {
        color: DarkRed;
        text-align: center;
    }
    h2 {
        color: DarkBlue;
    }
    h3 {
        color: DarkGreen;
    }
    h4 {
        color: DarkMagenta;
    }
    strong {
        color: Maroon;
    }
    em {
        color: Maroon;
    }
    img {
        display: block;
        margin-left: auto;
        margin-right: auto
    }
    code {
        color: SlateBlue;
    }
    mark {
        background-color:GoldenRod;
    }
</style>
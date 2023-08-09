# Traffic Management

## Usecase

To explore Cloud Run's support on revisions, rollback, gradual rollout, and traffic-split.

## About the Sample App

The sample app is a simple API that return the version number of the application. For example, when accessed, the version 1 of the application will be returned the following information:
```bash
API Microservice example: v1
```

And version 2 of the application will be returned 
```bash
API Microservice example: v2
```

## Deploy the Service of Version 1

```bash
gcloud run deploy product-service \
   --image gcr.io/qwiklabs-resources/product-status:0.0.1 \
   --tag test1 \
   --allow-unauthenticated
```

Each new Cloud Run service revision can be assigned a tag. Doing this will create a revision specific URL which is unknown to the client. This can be useful to handle the traffic profile across multiple revisions. 

When deploying the service, you can specify the tag name using `--tag`. Here we use `test1` to associate with version 1.

The revision specific URL follows the following convention:
```text
https://[tag-name]---[service-url]
```

For example, our product service URL is 

https://product-service-oy6beuif2a-uc.a.run.app, 

and the URL for revision 1 is 

https://test1---product-service-oy6beuif2a-uc.a.run.app. 

To test the service, run the following command:
```bash
TEST1_PRODUCT_SERVICE_URL=$(gcloud run services describe product-service --platform managed --format="value(status.address.url)")

curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
$TEST1_PRODUCT_SERVICE_URL/help -w "\n"
```
The output should be:
```text
.
API Microservice example: v1
.
```

## Deploy the Version 2 of the Service without Serving Traffic

We will deploy the version 2 of the service without assigning real traffic. We want to test it in production first. This can be done using the `--no-traffic` flag.

```bash
gcloud run deploy product-service \
  --image gcr.io/qwiklabs-resources/product-status:0.0.2 \
  --no-traffic \
  --tag test2 \
  --allow-unauthenticated
```

The revision specific URL should be https://test2---product-service-oy6beuif2a-uc.a.run.app.

To verify revision 2 is deployed successful, run this command:
```bash
TEST2_PRODUCT_STATUS_URL=$(gcloud run services describe product-service --platform managed --format="value(status.traffic[1].url)")

curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
$TEST2_PRODUCT_STATUS_URL/help -w "\n"
```

The output should be as follows:
```text
API Microservice example: v2
```

## Gradual Rolling out Revison 2

Let's migrate 50% of the traffic to the revision tagged `test2` using the `update-traffic` command:
```bash
gcloud run services update-traffic product-service \
  --to-tags test2=50
```

Confirm that original service URL now is distributing traffic:
```bash
for i in {1..10}; \
do curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
$TEST1_PRODUCT_SERVICE_URL/help -w "\n"; \
done
```

The output should be similar to the following:
```text
API Microservice example: v2
API Microservice example: v1
API Microservice example: v2
API Microservice example: v1
API Microservice example: v2
API Microservice example: v1
API Microservice example: v2
API Microservice example: v1
API Microservice example: v1
API Microservice example: v2
```

## Rolling back a Tagged Version

In the event an issue is found, the traffic migration can be rolled back by resetting the percentage.

```bash
gcloud run services update-traffic product-service \
  --to-tags test2=0
```

Test the endpoint is distributing traffic to version 1 (tagged `test1`) only:
```bash
for i in {1..10}; \
do curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
$TEST1_PRODUCT_SERVICE_URL/help -w "\n"; \
done
```

## Traffic Splitting

Let's deploy two more revisions of the service without routing traffic to them.

```bash
gcloud run deploy product-service \
  --image gcr.io/qwiklabs-resources/product-status:0.0.3 \
  --no-traffic \
  --tag test3 \
  --allow-unauthenticated

gcloud run deploy product-service \
  --image gcr.io/qwiklabs-resources/product-status:0.0.4 \
  --no-traffic \
  --tag test4 \
  --allow-unauthenticated
```

Create an environment variable for available revisions:
```bash
LIST=$(gcloud run services describe product-service --platform=managed --region=$LOCATION --format='value[delimiter="=25,"](status.traffic.revisionName)')"=25"; echo $LIST
```
```text
product-service-00001-nuv=25,product-service-00002-gir=25,product-service-00005-dat=25,product-service-00006-tiy=25
```

Now split traffic among the four revisions using the environment variable:
```bash
gcloud run services update-traffic product-service \
  --to-revisions $LIST
```

Test the endpoint is distributing traffic:
```bash
for i in {1..10}; \
do curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
$TEST1_PRODUCT_SERVICE_URL/help -w "\n"; \
done
```

The output should be something like this:
```text
API Microservice example: v4
API Microservice example: v2
API Microservice example: v3
API Microservice example: v2
API Microservice example: v1
API Microservice example: v2
API Microservice example: v3
API Microservice example: v3
API Microservice example: v1
API Microservice example: v1
```

## Update Traffic to the Latest Version only

Reset the service traffic profile to use the latest deployment:
```bash
gcloud run services update-traffic product-service --to-latest --platform=managed
```

Verify that the latest revision is able to receive traffic:
```bash
for i in {1..10}; \
do curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
$TEST1_PRODUCT_SERVICE_URL/help -w "\n"; \
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
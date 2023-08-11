# Service-to-Service Authentication

## User Case

A Cloud Run service `client` calls another Cloud Run service `employee-api` with service-to-service authentication.

![img](../img/iam.png)

- Employee API only accepts requests from a specific service account.

- The client must present a valid ODIC token when calling the API. The client does this by calling the local metadata server using GCP client library.

## Environment Variables

```bash
GOOGLE_CLOUD_PROJECT=ibcwe-event-layer-f3ccf6d9
REPOSITORY=us-central1-docker.pkg.dev/$GOOGLE_CLOUD_PROJECT/cloud-run-try
```

## Build

### Employee API

```bash
gcloud builds submit --tag $REPOSITORY/employee
```

### Client

```bash
gcloud builds submit --tag $REPOSITORY/client:bad

gcloud builds submit --tag $REPOSITORY/client
```

## Deployment

### Employee API

```bash
gcloud run deploy employee-api --image $REPOSITORY/employee
```

### Service Account
```bash
gcloud iam service-accounts create emp-api-client
```

Role binding
```bash
gcloud run services add-iam-policy-binding employee-api \
--member=serviceAccount:emp-api-client@ibcwe-event-layer-f3ccf6d9.iam.gserviceaccount.com \
--role=roles/run.invoker
```

### Client (bad version) 
```bash
EMPLOYEE_API=$(gcloud run services describe employee-api \
--format "value(status.url)")

gcloud run deploy client --image $REPOSITORY/client:bad \
--set-env-vars EMPLOYEE_API=$EMPLOYEE_API \
--service-account \
emp-api-client@ibcwe-event-layer-f3ccf6d9.iam.gserviceaccount.com
```

#### Test: call Client which will call the API

```bash
CLIENT_URL=$(gcloud run services describe client \
--format "value(status.url)")

curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" $CLIENT_URL
```

And we get:
```text
<html><head>
<meta http-equiv="content-type" content="text/html;charset=utf-8">
<title>403 Forbidden</title>
</head>
<body text=#000000 bgcolor=#ffffff>
<h1>Error: Forbidden</h1>
<h2>Your client does not have permission to get URL <code>/api/help</code> from this server.</h2>
<h2></h2>
</body></html>
```

### Client (Good Version)
```bash
gcloud run deploy client --image $REPOSITORY/client \
--set-env-vars EMPLOYEE_API=$EMPLOYEE_API \
--service-account \
emp-api-client@ibcwe-event-layer-f3ccf6d9.iam.gserviceaccount.com
```

Test again and the call to the API from the client will succeed.
```bash
CLIENT_URL=$(gcloud run services describe client \
--format "value(status.url)")

curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" $CLIENT_URL
```
```text
Employee API v1
```

## Cleanup

```bash
gcloud run services delete employee-api
gcloud run services delete client
gcloud iam service-accounts delete emp-api-client@ibcwe-event-layer-f3ccf6d9.iam.gserviceaccount.com
gcloud artifacts packages delete --repository=cloud-run-try --location=us-central1 employee
gcloud artifacts packages delete --repository=cloud-run-try --location=us-central1 client
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
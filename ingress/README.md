# Cloud Run Ingress Control

## Use Case

Only allow internal requests to reach the Cloud Run service. Internal traffic includes:
- traffic from VPCs 
- traffic certain Google Cloud services in the same GCP project, 
- traffic from Shared VPC
- traffic from regional internal Application Load Balancers, and 
- traffic allowed by VPC service controls

## Configuring Ingress

```bash
gcloud run services update employee-api --ingress=internal
```

## Test

Run the following command on a VM in the same GCP project:
```bash
EMPLOYEE_API=$(gcloud run services describe employee-api --format "value(status.url)")

curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
$EMPLOYEE_API/api/help
```
We will get `Employee API v1`.

Now if we call the same API from on-premise laptop, 
```bash
curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
https://employee-api-oy6beuif2a-uc.a.run.app/api/help
```

we will get an error:
```text
<html><head>
<meta http-equiv="content-type" content="text/html;charset=utf-8">
<title>404 Page not found</title>
</head>
<body text=#000000 bgcolor=#ffffff>
<h1>Error: Page not found</h1>
<h2>The requested URL was not found on this server.</h2>
<h2></h2>
</body></html>
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
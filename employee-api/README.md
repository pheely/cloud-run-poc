# Employee API - a POC on Google Cloud Run and Terraform

## Purpose

- This POC is to explore how Google Cloud Run can be used as a solution for some simple microservices. 

- Use Terraform to deploy all the required resources to support the idea of infrastructure as code.

![img](../img/poc.png)

## Architecture

![img](../img/api.png)

- The application exposes an employee API endpoint which provides CRUD operations on the Employee resource.
- CloudSQL MySQL is used as the data store
- The application code uses the [Cloud SQL connectors](https://cloud.google.com/sql/docs/mysql/connect-connectors) library to connect to MySQL (other options: Unix Socket, TCP Socket)
- The database password is stored in Secret Manager
- The application is deployed as a Cloud Run service using a service account which has two roles:
	- roles/cloudsql.client
	- roles/secretmanager.secretAccessor

## Local testing

Only Docker engine is required. Just run the following command:
```bash
docker-compose up
```

## Build and Deployment Prerequisites 

1. Enable Artifact Registry in your project
	```bash
	gcloud services enable artifactregistry.googleapis.com
	```

2. Enable CloudSQL APIs
	```bash
	gcloud services enable sqladmin.googleapis.com
	gcloud services enable sql-component.googleapis.com
	```

3. Create a repository named `cloud-run-try` on `us-central1`. 
	```bash
	gcloud artifacts repositories create \
	--location us-central1 \
	--repository-format docker \
	cloud-run-try
	```
## Environment Variables

```bash
GOOGLE_CLOUD_PROJECT=ibcwe-event-layer-f3ccf6d9
REPOSITORY=us-central1-docker.pkg.dev/$GOOGLE_CLOUD_PROJECT/cloud-run-try
```

## Build

### Manual Build

1. Build the docker image.
	```bash
	docker build -t $REPOSITORY/employee .
	```
2. Set up credentials to access the repo
	```bash
	gcloud auth configure-docker us-central1-docker.pkg.dev
	```
3. Push the image
	```bash
	docker push $REPOSITORY/employee
	```
4. Verify the image is built and pushed successfully
	```bash
	gcloud run services list
	```

### Cloud Build

```bash
gcloud builds submit --tag $REPOSITORY/employee
```

## Deployment

### Manual deployment

#### CloudSQL database

1. Create a Cloud SQL instance named `sql-db`
	```bash
	gcloud sql instances create sql-db \
	--tier db-f1-micro \
	--database-version MYSQL_8_0 \
	--region us-central1
	```
2. Create the `hr` database
	```bash
	gcloud sql databases create hr --instance sql-db
	```

#### Secret

Create a secret named "DB_PASS":
```bash
echo -n "changeit" | gcloud secrets create DB_PASS --replication-policy automatic --data-file=-
```
To access the contents of the version 1 of the secret:
```bash
gcloud secrets versions access 1 --secret DB_PASS
```

#### Cloud Run service

1. Deploy
	```bash
	gcloud run deploy employee-api --image $REPOSITORY/employee \
	--service-account=gyre-dataflow@ibcwe-event-layer-f3ccf6d9.iam.gserviceaccount.com \
	--add-cloudsql-instances ibcwe-event-layer-f3ccf6d9:us-central1:sql-db \
	--set-env-vars DB_USER=root \
	--set-secrets DB_PASS=DB_PASS:1 \
	--set-env-vars DB_NAME=hr \
	--set-env-vars DB_PRIVATE_IP= \
	--set-env-vars INSTANCE_CONNECTION_NAME=ibcwe-event-layer-f3ccf6d9:us-central1:sql-db
	```
2. Check the service is deployed successfully
	```bash
	gcloud run service list
	```

### Deploying with Terraform

1. Initialization
	```bash
	cd td
	terraform init
	```
2. Creating a plan
	```bash
	terraform plan -out tfplan
	```
3. Applying the plan
	```bash
	terraform apply
	```

## Testing

### Table Preparation

1. Install cloud sql proxy
	```bash
	curl -o cloud-sql-proxy https://storage.googleapis.com/cloud-sql-connectors/cloud-sql-proxy/v2.4.0/cloud-sql-proxy.linux.amd64
	chmod +x cloud-sql-proxy
	sudo mv cloud-sql-proxy /usr/local/bin 
	```
2. Install MySQL Client
	```bash
	sudo apt-get install mysql-client
	```
3. Get the instance connection name
	```bash
	gcloud sql instances describe sql-db|grep -i connection | awk '{print $2}'
	```
4. Start the cloud sql proxy to connect to `sql-db`
	```bash
	cloud-sql-proxy --port 3306 ibcwe-event-layer-f3ccf6d9:us-central1:sql-db
	```
5. Create and populate the `employees` table
	```bash
	mysql -u root -p --host 127.0.0.1 hr < schema.sql
	```

### Cloud Run service endpoint

Get the service endpoint using the following command:

```bash
$ gcloud run services list
   SERVICE       REGION       URL   
✔  employee-api  us-central1  https://employee-api-oy6beuif2a-uc.a.run.app  
```

### Hit the endpoint

1. Display API version information
	```bash
	EMPLOYEE_API=$(gcloud run services describe employee-api \
	--format "value(status.url)")

	curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
	$EMPLOYEE_API/api/help
	```
	Output:
	```text
	Employee API v4
	```
2. List all employees
	```bash
	curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
	$EMPLOYEE_API/api/employee | jq
	```
	Output:
	```yaml
	[
		{
			"id": "1",
			"first_name": "John",
			"last_name": "Doe",
			"department": "Products",
			"salary": 200000,
			"age": 25
		},
		{
			"id": "2",
			"first_name": "Jane",
			"last_name": "Doe",
			"department": "Sales",
			"salary": 100000,
			"age": 22
		}
	]
	```
3. Create a new employee
	```bash
	curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
	-d '{"first_name":"Shrek","last_name":"Unknown","department":"Royal","salary":200000,"age":25}' \
	-X POST $EMPLOYEE_API/api/employee | jq
	```
	Output:
	```yaml
	{
		"id": "3",
		"first_name": "Shrek",
		"last_name": "Unknown",
		"department": "Royal",
		"salary": 200000,
		"age": 25
	}
	```

## Cleanup

### Manual cleanup

```bash
gcloud run services delete employee-api

gcloud sql instances delete sql-db	

gcloud artifacts packages delete employee --repository=cloud-run-try \
--location=us-central1
```

### Cleanup with Terraform

1. Delete all the resources
	```bash
	cd tf
	terraform destroy
	```

2. Delete docker images (packages) in the Artifact Registry repository: same as above
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
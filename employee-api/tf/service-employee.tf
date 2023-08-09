resource "google_cloud_run_v2_service" "employee" {
  name     = "employee-api"
  location = var.region_id

  template {
    volumes {
      name = "cloudsql"
      cloud_sql_instance {
        instances = [google_sql_database_instance.cloudsql-instance.connection_name]
      }
    }

    containers {
      image = "us-central1-docker.pkg.dev/ibcwe-event-layer-f3ccf6d9/cloud-run-try/employee"

      env {
        name = "DB"
        value = "mysql://employee-api:changeit@unix(/cloudsql/ibcwe-event-layer-f3ccf6d9:us-central1:sql-db)/hr"
      }
      volume_mounts {
        name = "cloudsql"
        mount_path = "/cloudsql"
      }
    }
  }
}


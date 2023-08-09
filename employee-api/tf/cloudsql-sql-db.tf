resource "google_sql_database_instance" "cloudsql-instance" {
  name = "sql-db"
  database_version = "MYSQL_8_0"
  region = var.region_id

  settings {
    tier = "db-f1-micro"
  }

  deletion_protection = "false"
}

resource "google_sql_database" "database" {
  name = "hr"
  instance = google_sql_database_instance.cloudsql-instance.name
}

resource "google_sql_user" "admin" {
  name     = "root"
  instance = google_sql_database_instance.cloudsql-instance.name
  host     = "%"
  password = "changeit"
}

resource "google_sql_user" "app-user" {
  name     = "employee-api"
  instance = google_sql_database_instance.cloudsql-instance.name
  host     = "%"
  password = "changeit"
}

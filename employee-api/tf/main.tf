provider google {
    project = var.project_id
    region = var.region_id
}

variable "project_id" {
  type        = string
  default     = "ibcwe-event-layer-f3ccf6d9"
  description = "The Google Cloud Project ID to use"
}

variable "region_id" {
    type        = string
    default     = "us-central1"
    description = "The Google Cloud compute region to use"
}

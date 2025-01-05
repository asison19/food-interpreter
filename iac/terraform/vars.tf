variable "GCP_PROJECT_ID" {
    type = string
    default = ""
    description = "GCP Project ID of where to create the infrastructure."
}

variable "GCP_STATE_BUCKET" {
    type = string
    default = ""
    description = "GCP Terraform state GCS bucket name."
}

module "state" {
    source = "github.com/asison19/terraform-state-gcs"
    GCP_PROJECT_ID = var.GCP_PROJECT_ID
}
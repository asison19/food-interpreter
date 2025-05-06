variable "GCP_PROJECT_ID" {
    type = string
    default = ""
    description = "GCP Project ID of where to create the infrastructure."
}

variable "GCP_PROJECT_REGION" {
    type = string
    default = ""
    description = "GCP Project region of where to create the infrastructure."
}

variable "SHA_SHORT" {
    type = string
    default = ""
    description = "Git SHA shortened used to use the latest image for Cloud Run."
}

variable "FOOD_INTERPRETER_IMAGE_VERSION" {
    type = string
    default = ""
    description = "Version of the food interpreter docker image"
}

variable "ENABLE_KUBERNETES" {
    type = string
    default = "false"
    description = "If enabled, GKE is used."
}

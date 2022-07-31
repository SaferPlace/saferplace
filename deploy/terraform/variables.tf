variable "region" {
  description = "Deployment region"
  type        = string
}

variable "datacenters" {
  description = "List of datacenters to deploy"
  type        = list(string)
}

variable "namespace" {
  description = "Deployment namespace"
  type        = string
}

variable "image" {
  description = "Deployment image"
  type        = string
}

variable "docker_password" {
  description = "Read only token for the container registry"
  type        = string
  sensitive   = true
}

variable "eircode_token" {
  description = "Token used to authenticate to the eircode service"
  type        = string
  sensitive   = true
}

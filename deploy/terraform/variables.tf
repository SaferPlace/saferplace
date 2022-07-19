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

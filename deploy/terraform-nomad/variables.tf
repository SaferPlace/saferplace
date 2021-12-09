variable "region" {
  description = "Deployment region"
  type        = string
}

variable "datacenters" {
  description = "List of datacenters to deploy"
  type        = list(string)
}

variable "image" {
  description = "Deployment image"
  type        = string
}

variable "tags" {
  description = "List of service tags"
  type        = list(string)
}

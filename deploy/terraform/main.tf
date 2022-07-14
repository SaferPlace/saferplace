resource "nomad_namespace" "saferplace" {
  name        = var.namespace
  description = "Safer Place namespace"
}

module "saferplace" {
  source  = "VoyTechnology/nomad/generic"
  version = "0.0.2"

  job_name    = "saferplace"
  namespace   = var.namespace
  region      = var.region
  datacenters = var.datacenters

  image = var.image
  ports = {
    "http" = {
      to     = 8080
      static = 0
    }
  }

  service_port = "http"
  service_tags = [
    "traefik.http.routers.saferplace.rule=Host(`safer.place`)"
  ]
}

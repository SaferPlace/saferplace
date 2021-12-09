resource "nomad_job" "saferplace" {
  jobspec = file("${path.module}/../nomad/saferplace.nomad")

  hcl2 {
    enabled = true

    #Use a reduced variable set until needed
    vars = {
      "region"       = var.region,
      "datacenters"  = var.datacenters,
      "image"        = var.image,
      "service_tags" = var.tags,
    }
  }
}
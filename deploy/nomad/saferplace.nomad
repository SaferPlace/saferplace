variable "region" {
  description = "Region for the deployment"
  type        = string
  default     = "global"
}

variable "datacenters" {
  description = "List of datacenters in the region to deploy to"
  type        = list(string)
  default     = ["dub1"]
}

variable "namespace" {
  description = "Namespace of the deployment"
  type        = string
  default     = "default"
}

variable "image" {
  description = "Full saferplace image to deploy"
  type        = string
  default     = "ghcr.io/saferplace/saferplace"
}

variable "reschedule_attempts" {
  description = "Number of reschedule attempts"
  type        = number
  default     = 5
}

variable "reschedule_interval" {
  description = "Interval between reschedules"
  type        = string
  default     = "1h"
}

variable "reschedule_delay_duration" {
  description = "How long to wait until we reschedule"
  type        = string
  default     = "1m"
}

variable "reschedule_delay_function" {
  description = "Function to use for rescheduling: constant or exponential"
  type        = string
  default     = "exponential"
}

variable "reschedule_delay_max" {
  description = "How long to wait until we reschedule"
  type        = string
  default     = "1h"
}

variable "reschedule_unlimited" {
  description = "Enable unlimited reschedules"
  type        = bool
  default     = false
}

variable "restart_interval" {
  description = "Interval between restarting"
  type        = number
  default     = 5
}

variable "restart_delay_duration" {
  description = "How long to wait until we restart"
  type        = string
  default     = "10s"
}

variable "restart_attempts" {
  description = "Number of allowed attempts at restarting"
  type        = number
  default     = 5
}

variable "restart_mode" {
  description = "Restart mode, either delay of fail"
  type        = string
  default     = "fail"
}

variable "resources_cpu" {
  description = "CPU Mhz that should be allocated"
  type        = number
  default     = 50
}

variable "resources_memory" {
  description = "Memory that should be allocated, in MiB"
  type        = number
  default     = 50
}

variable "service_tags" {
  description = "List of tags for the service"
  type        = list(string)
  default     = []
}

variable "service_check_interval" {
  description = "Interval at which the service health should be checked"
  type        = string
  default     = "10s"
}

variable "service_check_timeout" {
  description = "Timeout for the service status health"
  type        = string
  default     = "1s"
}

job "saferplace" {
  region      = var.region
  datacenters = var.datacenters
  namespace   = var.namespace

  type = "service"

  # For now we only have linux images, but we might change it in the future.
  constraint {
    attribute = "${attr.kernel.name}"
    value     = "linux"
  }

  reschedule {
    attempts       = var.reschedule_attempts
    interval       = var.reschedule_interval
    delay          = var.reschedule_delay_duration
    delay_function = var.reschedule_delay_function
    unlimited      = "${var.reschedule_unlimited}"
  }

  group "saferplace" {
    network {
      port "http" {}
    }

    task "saferplace" {
      driver = "docker"
      config {
        image = var.image
        ports = ["http"]
      }

      env {
        PORT = "${NOMAD_PORT_http}"
      }

      resources {
        cpu    = var.resources_cpu
        memory = var.resources_memory
      }

      service {
        tags = var.service_tags
        port = "http"

        check {
          type     = "http"
          port     = "http"
          path     = "/"
          interval = var.service_check_interval
          timeout  = var.service_check_timeout
        }
      }
    }
  }
}
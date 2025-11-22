variable "project_name" {
  type = string
  description = "the current project name"
}

variable "organization_id" {
  type = string
}

variable "roles" {
  description = "map of defined roles."
  type = map(object({
    key = string
    display_name = string
  }))
}

variable "services" {
  description = "map of registered applications."
  type = map(object({
    name = string
  }))
}

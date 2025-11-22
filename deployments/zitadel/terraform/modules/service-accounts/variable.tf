variable "services" {
  description = "map of registred applications."
  type = map(object({
    service_key = string
    roles = map(string)
  }))
}

variable "project_id" {
  description = "Project ID from IAM module"
  type        = string
}

variable "service_user_ids" {
  description = "Map of service names to their user IDs from IAM module"
  type        = map(string)
}

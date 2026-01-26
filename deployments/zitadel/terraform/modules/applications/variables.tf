variable "applications" {
  description = "map of registered application"
  type = map(object({
    type =  string
    name = string
    response_types = list(string)
    grant_types = list(string)
    auth_method_type = string
    access_token_type = string
    version = string
  }))
}

variable "redirect_uris" {
  description = "list of redirect URIs"
  type = list(string)
}

variable "project_id" {
  description = "project ID from applications module."
  type = string
}

variable "zitadel_domain" {
  type = string
  description = "zitadel_domain represent the zitadel domain name. ex: localhost"
}

variable "zitadel_secure_mode" {
  type = bool
  description = "zitadel_secure_mode represent the insecure toggle for zitadel, true for insecure and false for secure mode."
}

variable "zitadel_port" {
  type = string
  description = "zitadel_port (actually for v1/ is the current zitadel port running."
}
variable "zitadel_token_path" {
  type = string
  description = "path to the token profil file"
}

variable "project_name" {
  type = string
  description = "streamly default project name"
  default = "local"
}

variable "organization_id" {
  type = string
  description = "streamly zitadel organization id" // load from tfvars.
}

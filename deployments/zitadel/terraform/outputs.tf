output "project_id" {
  description = "The ID of the Zitadel project"
  value       = module.iam.project_id
}

output "service_credentials" {
  description = "Service account credentials for all services"
  sensitive   = true
  value       = module.iam.service_credentials
}

output "root_admin_credentials" {
  description = "Root admin user credentials (user_id and PAT)"
  sensitive   = true
  value       = module.iam.root_admin_credentials
}

output "applications" {
  description = "Application client IDs for frontend configuration"
  sensitive = true
  value       = module.applications.applications
}

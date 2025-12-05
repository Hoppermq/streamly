output "project_id" {
  description = "The ID of the Zitadel project"
  value       = module.iam.project_id
}

output "service_credentials" {
  description = "Service account credentials for all services"
  sensitive   = true
  value       = module.iam.service_credentials
}

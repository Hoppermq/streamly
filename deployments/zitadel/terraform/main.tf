terraform {
  required_providers {
    zitadel = {
      source  = "zitadel/zitadel"
      version = "2.7.0"
    }
  }
}

provider "zitadel" {
  domain = var.zitadel_domain
  insecure = var.zitadel_secure_mode
  port = var.zitadel_port
  jwt_profile_json = file(var.zitadel_jwt_profile_file)
}

module "common" {
  source = "./modules/shared"
}

module "iam" {
  source = "./modules/iam"
  providers = {
    zitadel = zitadel
  }

  organization_id = var.organization_id
  project_name    = var.project_name

  roles    = module.common.roles
  services = module.common.service_definitions
}

module "service-accounts" {
  source = "./modules/service-accounts"
  providers = {
    zitadel = zitadel
  }
  depends_on = [module.iam]

  project_id = module.iam.project_id
  service_user_ids = module.iam.service_user_ids
  services = module.common.service_role_mappings
}

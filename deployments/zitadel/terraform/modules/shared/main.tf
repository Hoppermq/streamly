locals {
  roles = {
    super_admin = {
      key = "super_admin"
      display_name = "Super Admin"
    },
    admin = {
      key = "admin"
      display_name = "Admin"
    },
    technical = {
      key = "technical"
      display_name = "Technical"
    }
  }

  services = {
    auth = {
      scopes =["openid", "streamly:events:read", "streamly:events:write"],
      roles = ["technical"],
    },
    ingestor = {
      scopes =["openid", "streamly:events:read", "streamly:events:write"],
      roles = ["technical"]
    },
    query ={
      scopes = ["openid", "streamly:events:read"],
      roles = ["technical"],
    } ,
    processor = {
      scopes = ["openid", "streamly:events:read", "streamly:events:write"],
      roles = ["technical"]
    }
  }

  # Derive service definitions for IAM module (machine user creation)
  service_definitions = {
    for key, config in local.services : key => {
      name = key
    }
  }

  # Derive service role mappings for service-accounts module
  service_role_mappings = {
    for key, config in local.services : key => {
      service_key = key
      roles = {
        for role in config.roles : role => role
      }
    }
  }
}

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
    ingestor = {
      name = "ingestor-service"
    },
    query = {
      name = "query-service"
    },
    processor = {
      name = "processor-service"// should be imported to kep consistancy across other comps
    }
  }

  service_role_mappings = {
    ingestor = {
      service_key = "ingestor"
      roles = {
        "technical" = "technical"
      }
    },
    query = {
      service_key = "query" 
      roles = {
        "technical" = "technical"
      }
    },
    processor = {
      service_key = "processor"
      roles = {
        "technical" = "technical"
      }
    }
  }
}

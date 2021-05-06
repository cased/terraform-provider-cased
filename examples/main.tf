provider "cased" {
  # Grab your Workflow API key from the API keys page:
  # https://app.cased.com/settings/api-keys
  workflows_api_key = ""
}

resource "cased_workflow" "database-provisioning" {
  name = "database-provisioning"

  controls {
    reason = true

    approval {
      count         = 2
      self_approval = false

      # Restrict who can respond to the request
      responders {
        # User
        responder {
          name = "ceo@company.com"
        }

        # Someone from managers group must respond
        responder {
          name     = "managers"
          required = true
        }
      }

      sources {
        slack {
          channel = "#cased-workflows"
        }
      }
    }
  }
}

# Terraform Provider Cased

A Terraform provider for managing resources on your Cased account.

Currently, this provider supports:

- Managing workflows
- Managing webhooks

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
-	[Go](https://golang.org/doc/install) >= 1.15

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```sh
$ go install
```

## Using the provider

### Creating workflows workflow

```tf
terraform {
  required_providers {
    cased = {
      source = "cased/cased"
      version = "0.0.1"
    }
  }
}

provider "cased" {
  # Grab your Workflow API key from the API keys page:
  # https://app.cased.com/settings/api-keys
  workflows_api_key = ""
}

resource "cased_webhooks_endpoint" "webhook-endpoint" {
  url = "https://app.cased.com/webhooks/endpoint"

  // event_types is optional, subscribed to all event types of left blank
  event_types = ["event.created"]
}

// The webhook secret is accessible so you can set environment variables
output "cased_webhooks_endpoint_secret" {
  value = cased_webhooks_endpoint.webhook-endpoint.secret
}

resource "cased_workflow" "database-provisioning" {
  name = "database-provisioning"

  # At least one control is required: reason, authentication, or approval.
  controls {
    reason = true
    authentication = true

    approval {
      # Optional: How many approvals are required to approve the workflow. (Default: 1)
      count = 2

      # Optional: Can the authenticated user self approve the request? (Default: false)
      self_approval = false

      # Optional: Restrict who can respond to the request
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

      # Specify where the approval request will get delivered to.
      sources {
        # You must connect your Slack workspace with your Cased account on your
        # integrations page: https://app.cased.com/settings/integrations
        slack {
          channel = "#cased-workflows"
        }
      }
    }
  }
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

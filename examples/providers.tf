# Copyright (c) HashiCorp, Inc.

terraform {
  required_providers {
    rightbrain = {
      source = "RightbrainAI/tasks"
    }
  }
}

provider "rightbrain" {
  client_id     = "<client-id>"
  client_secret = "<client-secret>"
  org_id        = "<org-id>"
  project_id    = "<project-id>"
}

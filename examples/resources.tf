# Copyright (c) HashiCorp, Inc.

resource "rightbrain_task" "tell-me-a-joke" {
  name        = "Tell me a Joke!"
  description = "Tells a joke about a given subject :)"
  enabled     = true

  llm_model_id  = data.rightbrain_model.gpt-4o-mini.id
  system_prompt = "You can tell good jokes about anything"
  user_prompt   = "Tell me a joke about {{subject}}"
  output_format = {
    "joke" : "str"
  }
}
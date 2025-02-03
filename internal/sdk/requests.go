package sdk

type CreateTaskRequest struct {
	Description   string                  `json:"description"`
	Enabled       bool                    `json:"enabled"`
	ImageRequired bool                    `json:"image_required"`
	LLMModelID    string                  `json:"llm_model_id"`
	Name          string                  `json:"name"`
	OutputFormat  map[string]OutputFormat `json:"output_format"`
	Public        bool                    `json:"public"`
	SystemPrompt  string                  `json:"system_prompt"`
	UserPrompt    string                  `json:"user_prompt"`
}

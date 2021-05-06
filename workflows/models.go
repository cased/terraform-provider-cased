package workflows

type Workflow struct {
	ID         string      `json:"id,omitempty"`
	Name       string      `json:"name,omitempty"`
	ApiURL     string      `json:"api_url,omitempty"`
	Conditions []Condition `json:"conditions,omitempty"`
	Controls   Controls    `json:"controls,omitempty"`
	UpdatedAt  string      `json:"updated_at,omitempty"`
	CreatedAt  string      `json:"created_at,omitempty"`
}

type Controls struct {
	Reason         bool             `json:"reason,omitempty"`
	Authentication bool             `json:"authentication,omitempty"`
	Approval       *ApprovalControl `json:"approval,omitempty"`
}

type ApprovalControl struct {
	Count        int                     `json:"count,omitempty"`
	SelfApproval bool                    `json:"self_approval,omitempty"`
	Duration     int                     `json:"duration,omitempty"`
	Timeout      int                     `json:"timeout,omitempty"`
	Responders   *Responders             `json:"responders,omitempty"`
	Sources      *ApprovalControlSources `json:"sources,omitempty"`
}

type Responders map[string]string

type ApprovalControlSources struct {
	Email bool                        `json:"email,omitempty"`
	Slack *ApprovalControlSourceSlack `json:"slack,omitempty"`
}

type ApprovalControlSourceSlack struct {
	Channel string `json:"channel,omitempty"`
}

type Condition struct {
	Field    string `json:"field"`
	Value    string `json:"value"`
	Operator string `json:"operator"`
}

package workflows

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetWorkflows - Returns list of workflows (auth required)
func (c *Client) GetWorkflows() ([]Workflow, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/workflows", c.URL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	workflows := []Workflow{}
	err = json.Unmarshal(body, &workflows)
	if err != nil {
		return nil, err
	}

	return workflows, nil
}

// CreateWorkflow - Creates a workflow (auth required)
func (c *Client) CreateWorkflow(workflow Workflow) (*Workflow, error) {
	postBody, err := json.Marshal(workflow)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/workflows", c.URL), bytes.NewBuffer(postBody))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	resp := &Workflow{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// UpdateWorkflow - Updates a workflow (auth required)
func (c *Client) UpdateWorkflow(workflow Workflow) (*Workflow, error) {
	postBody, err := json.Marshal(workflow)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/workflows/%s", c.URL, workflow.ID), bytes.NewBuffer(postBody))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	resp := &Workflow{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetWorkflow - Get a workflow (auth required)
func (c *Client) GetWorkflow(workflowID string) (*Workflow, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/workflows/%s", c.URL, workflowID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	resp := &Workflow{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// DeleteWorkflow - Delete a workflow (auth required)
func (c *Client) DeleteWorkflow(workflowID string) (bool, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/workflows/%s", c.URL, workflowID), nil)
	if err != nil {
		return false, err
	}

	if _, err = c.doRequest(req); err != nil {
		return false, err
	}

	return true, nil
}

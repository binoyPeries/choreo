// Copyright 2025 The OpenChoreo Authors
// SPDX-License-Identifier: Apache-2.0

package authz

// Subject represents the entity requesting authorization
type SubjectType string

const (
	SubjectTypeToken SubjectType = "token"
)

type Subject struct {
	Type       SubjectType            `json:"type"`
	Id         string                 `json:"id"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// ResourceRef represents a reference to a resource
type ResourceRef struct {
	Type string `json:"type"`
	ID   string `json:"id"` // can be string, number, or object
}

// Resource represents a resource being accessed
type Resource struct {
	Type string `json:"type"`
	ID   string `json:"id,omitempty"`
	// Path represents an ordered list of resource path items forming a hierarchy
	Path       []ResourceRef          `json:"path"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// Action represents the action being performed on a resource
type Action struct {
	Name       string                 `json:"name"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// Context provides additional context for the authorization request
type Context struct {
	RequestID  string                 `json:"request_id,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// Decision represents the authorization decision
type Decision struct {
	Decision bool             `json:"decision"`
	Context  *DecisionContext `json:"context,omitempty"`
}

// DecisionContext provides additional context for the decision
type DecisionContext struct {
	Reason string `json:"reason,omitempty"`
}

// EvaluateRequest represents an authorization access evaluation request
type EvaluateRequest struct {
	Subject  Subject  `json:"subject"`
	Resource Resource `json:"resource"`
	Action   Action   `json:"action"`
	Context  *Context `json:"context,omitempty"`
}

// EvaluateResponse represents an authorization access evaluation response
type EvaluateResponse struct {
	Decision Decision `json:"decision"`
}

// EvaluationSemantic defines the evaluation semantic options
type EvaluationSemantic string

const (
	EvaluationSemanticExecuteAll          EvaluationSemantic = "execute_all"
	EvaluationSemanticDenyOnFirstDeny     EvaluationSemantic = "deny_on_first_deny"
	EvaluationSemanticPermitOnFirstPermit EvaluationSemantic = "permit_on_first_permit"
)

// EvaluationItem represents a single evaluation item that can inherit from top-level fields
type EvaluationItem struct {
	Subject  *Subject  `json:"subject,omitempty"`
	Resource *Resource `json:"resource,omitempty"`
	Action   *Action   `json:"action,omitempty"`
	Context  *Context  `json:"context,omitempty"`
}

// EvaluationOptions contains options for batch evaluations
type EvaluationOption struct {
	EvaluationSemantic EvaluationSemantic    `json:"evaluation_semantic,omitempty"`
	Properties         map[string]interface{} `json:",inline"`
}

// EvaluatesRequest represents a batch authorization access evaluation request
type EvaluatesRequest struct {
	Subject     *Subject           `json:"subject,omitempty"`
	Resource    *Resource          `json:"resource,omitempty"`
	Action      *Action            `json:"action,omitempty"`
	Context     *Context           `json:"context,omitempty"`
	Evaluations []EvaluationItem   `json:"evaluations,omitempty"`
	Options     EvaluationOption `json:"option"`
}

// EvaluatesResponse represents a batch authorization access evaluation response
// It can be either a single decision (compatible with /evaluate) or a batch of decisions
type EvaluatesResponse struct {
	// Single evaluation response (compatible with /evaluate)
	Decision *Decision `json:"decision,omitempty"`
	// Batch response: one Decision per item, in order
	Evaluations []Decision `json:"evaluations,omitempty"`
}

// SearchResourceSpec represents a resource specification for search
type SearchResourceSpec struct {
	Type       string                 `json:"type"`
	Path       []ResourceRef          `json:"path"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// SearchResourcesRequest represents a resource search API request
type SearchResourcesRequest struct {
	Subject  Subject            `json:"subject"`
	Action   Action             `json:"action"`
	Resource SearchResourceSpec `json:"resource"`
	Context  *Context           `json:"context,omitempty"`
}

// SearchResourcesResponse represents a resource search API response
type SearchResourcesResponse struct {
	Results []Resource `json:"results"`
}

// SearchActionsRequest represents an action search API request
type SearchActionsRequest struct {
	Subject  Subject  `json:"subject"`
	Resource Resource `json:"resource"`
	Context  *Context `json:"context,omitempty"`
}

// SearchActionsResponse represents an action search API response
type SearchActionsResponse struct {
	Results []Action `json:"results"`
}

// SubjectFilter represents a filter for the kind of subjects to return
type SubjectFilter struct {
	Type string `json:"type"`
}

// SubjectResult represents a subject returned in search results
type SubjectResult struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// SearchSubjectsRequest represents a subject search API request
type SearchSubjectsRequest struct {
	Subject  SubjectFilter `json:"subject"`
	Resource Resource      `json:"resource"`
	Action   Action        `json:"action"`
	Context  *Context      `json:"context,omitempty"`
}

// SearchSubjectsResponse represents a subject search API response
type SearchSubjectsResponse struct {
	Results []SubjectResult `json:"results"`
}

type AuthzPDP interface {
	Evaluate(req *EvaluateRequest) (EvaluateResponse, error)
	Evaluates(req *EvaluatesRequest) (EvaluatesResponse, error)
	SearchResources(req *SearchResourcesRequest) (SearchResourcesResponse, error)
	SearchActions(req *SearchActionsRequest) (SearchActionsResponse, error)
	SearchSubjects(req *SearchSubjectsRequest) (SearchSubjectsResponse, error)
}

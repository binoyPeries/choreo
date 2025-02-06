/*
 * Copyright (c) 2025, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package api

import (
	choreov1 "github.com/wso2-enterprise/choreo-cp-declarative-api/api/v1"
)

// ListParams defines common parameters for listing resources
type ListParams struct {
	OutputFormat string
	Name         string
}

// ListProjectParams defines parameters for listing projects
type ListProjectParams struct {
	Organization string
	OutputFormat string
	Name         string
}

// ListComponentParams defines parameters for listing components
type ListComponentParams struct {
	Organization string
	Project      string
	OutputFormat string
	Name         string
}

// CreateOrganizationParams defines parameters for creating organizations
type CreateOrganizationParams struct {
	Name        string
	DisplayName string
	Description string
}

// CreateProjectParams defines parameters for creating projects
type CreateProjectParams struct {
	Organization string
	Name         string
	DisplayName  string
	Description  string
}

// CreateComponentParams contains parameters for component creation
type CreateComponentParams struct {
	Name             string
	DisplayName      string
	Type             choreov1.ComponentType
	Organization     string
	Project          string
	Description      string
	GitRepositoryURL string
	Branch           string
	Context          string
	DockerFile       string
	BuildConfig      string
	Image            string
	Tag              string
	Port             int
	Endpoint         string
}

// ApplyParams defines parameters for applying configuration files
type ApplyParams struct {
	FilePath string
}

// LoginParams defines parameters for login
type LoginParams struct {
	KubeconfigPath string
	Kubecontext    string
}

type LogParams struct {
	Organization string
	Project      string
	Component    string
	Build        string
	Type         string
	Environment  string
	Follow       bool
	TailLines    int64
}

// CreateBuildParams contains parameters for build creation
type CreateBuildParams struct {
	// Basic metadata
	Name         string
	Organization string
	Project      string
	Component    string

	// Build configuration
	Docker    *choreov1.DockerConfiguration
	Buildpack *choreov1.BuildpackConfiguration
}

// ListBuildParams defines parameters for listing builds
type ListBuildParams struct {
	Organization string
	Project      string
	Component    string
	OutputFormat string
	Name         string
}

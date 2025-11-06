// Copyright 2025 The OpenChoreo Authors
// SPDX-License-Identifier: Apache-2.0

package authz 

// ConstructTokenSubject constructs a Subject of type token
func ConstructTokenSubject(token string) Subject {
	return Subject{
		Type: SubjectTypeToken,
		Id:   token,
	}
}

// ConstructComponentResource constructs a Resource of type component with the given hierarchy
func ConstructComponentResource(componentName, projectName, orgName string) Resource {
	return Resource{
		Type: "component",
		ID:   componentName,
		Path: []ResourceRef{
			{Type: "organization", ID: orgName},
			{Type: "project", ID: projectName},
		},
	}
}

// ConstructProjectResource constructs a Resource of type project with the given hierarchy
func ConstructProjectResource(projectName, orgName string) Resource {
	return Resource{
		Type: "project",
		ID:   projectName,
		Path: []ResourceRef{
			{Type: "organization", ID: orgName},
		},
	}
}

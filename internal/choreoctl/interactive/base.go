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

package interactive

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/wso2-enterprise/choreo-cp-declarative-api/internal/choreoctl/util"
)

const (
	StateOrgSelect = iota
	StateProjSelect
	StateCompSelect
)

// BaseModel holds the shared state for interactive models.
type BaseModel struct {
	Organizations            []string
	OrgCursor                int
	Projects                 []string
	ProjCursor               int
	Components               []string
	CompCursor               int
	Environments             []string
	EnvCursor                int
	DeploymentTracks         []string
	DeploymentTrackCursor    int
	DeployableArtifacts      []string
	DeployableArtifactCursor int
	ErrorMsg                 string
	State                    int
}

func NewBaseModel() (*BaseModel, error) {
	orgs, err := util.GetOrganizationNames()
	if err != nil {
		return nil, fmt.Errorf("failed to get organizations: %w", err)
	}

	if len(orgs) == 0 {
		return nil, fmt.Errorf("no organizations found")
	}

	return &BaseModel{
		Organizations: orgs,
	}, nil
}

// RunInteractiveModel starts a Bubble Tea program with the given model
// and returns the final model state after program completion.
func RunInteractiveModel(model tea.Model) (tea.Model, error) {
	p := tea.NewProgram(model)
	finalModel, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to run interactive mode: %w", err)
	}
	return finalModel, nil
}

func IsQuitKey(msg tea.KeyMsg) bool {
	switch msg.String() {
	case "q", "ctrl+c", "esc":
		return true
	default:
		return false
	}
}

func IsEnterKey(msg tea.KeyMsg) bool {
	return msg.String() == "enter"
}

func RenderInputPrompt(prompt, defaultText, currentText, errorMsg string) string {
	var view string
	if defaultText != "" {
		view = fmt.Sprintf("%s (default: %s)\n", prompt, defaultText)
	} else {
		view = prompt + "\n"
	}
	view += currentText
	if errorMsg != "" {
		view += "\nError: " + errorMsg
	}
	return view
}

func RenderListPrompt(header string, items []string, cursor int) string {
	s := header + "\n\n"
	for i, item := range items {
		cursorSymbol := " "
		if i == cursor {
			cursorSymbol = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursorSymbol, item)
	}
	return s
}

func EditTextInputField(msg tea.KeyMsg, input string, cursor int) (string, int) {
	switch msg.Type {
	case tea.KeyBackspace, tea.KeyDelete:
		if len(input) > 0 {
			runes := []rune(input)
			input = string(runes[:len(runes)-1])
			cursor = len(input)
		}
	case tea.KeyRunes:
		input += string(msg.Runes)
		cursor = len(input)
	case tea.KeySpace:
		input += " "
		cursor = len(input)
	}
	return input, cursor
}

func ProcessListCursor(msg tea.KeyMsg, cursor, listLength int) int {
	switch msg.String() {
	case "up", "k":
		if cursor > 0 {
			cursor--
		}
	case "down", "j":
		if cursor < listLength-1 {
			cursor++
		}
	}
	return cursor
}

// UpdateOrgSelect handles organization selection.
// It fetches projects when Enter is pressed.
func (b *BaseModel) UpdateOrgSelect(keyMsg tea.KeyMsg) tea.Cmd {
	if IsEnterKey(keyMsg) {
		projects, err := util.GetProjectNames(b.Organizations[b.OrgCursor])
		if err != nil {
			b.ErrorMsg = fmt.Sprintf("failed to get projects: %v", err)
			return nil
		}
		if len(projects) == 0 {
			b.ErrorMsg = fmt.Sprintf("no projects found in organization '%s'", b.Organizations[b.OrgCursor])
			return nil
		}
		b.Projects = projects
		b.State = StateProjSelect
		b.ErrorMsg = ""
		return nil
	}
	b.OrgCursor = ProcessListCursor(keyMsg, b.OrgCursor, len(b.Organizations))
	return nil
}

// UpdateProjSelect is a helper to handle project selection update.
// It returns a command to load components if needed.
func (b *BaseModel) UpdateProjSelect(keyMsg tea.KeyMsg) (tea.Cmd, error) {
	b.ProjCursor = ProcessListCursor(keyMsg, b.ProjCursor, len(b.Projects))
	if IsEnterKey(keyMsg) {
		// First fetch components for the selected project
		components, err := util.GetComponentNames(b.Organizations[b.OrgCursor], b.Projects[b.ProjCursor])
		if err != nil {
			return nil, fmt.Errorf("failed to get components: %w", err)
		}
		// Store components but don't set state - let the caller handle state transition
		b.Components = components
	}
	return nil, nil
}

// FetchDeploymentTracks retrieves deployment track names based on the current selections.
func (b *BaseModel) FetchDeploymentTracks() ([]string, error) {
	if b.OrgCursor >= len(b.Organizations) ||
		b.ProjCursor >= len(b.Projects) ||
		b.CompCursor >= len(b.Components) {
		return nil, fmt.Errorf("invalid selection indices for deployment tracks")
	}
	return util.GetDeploymentTrackNames(
		b.Organizations[b.OrgCursor],
		b.Projects[b.ProjCursor],
		b.Components[b.CompCursor],
	)
}

// FetchBuildNames retrieves build names based on the current selections.
func (b *BaseModel) FetchBuildNames() ([]string, error) {
	if b.OrgCursor >= len(b.Organizations) ||
		b.ProjCursor >= len(b.Projects) ||
		b.CompCursor >= len(b.Components) {
		return nil, fmt.Errorf("invalid selection indices for build names")
	}
	return util.GetBuildNames(
		b.Organizations[b.OrgCursor],
		b.Projects[b.ProjCursor],
		b.Components[b.CompCursor],
	)
}

// FetchDeployableArtifacts retrieves deployable artifact names based on the current selections.
func (b *BaseModel) FetchDeployableArtifacts() ([]string, error) {
	if b.OrgCursor >= len(b.Organizations) ||
		b.ProjCursor >= len(b.Projects) ||
		b.CompCursor >= len(b.Components) {
		return nil, fmt.Errorf("invalid selection indices for deployable artifacts")
	}
	return util.GetDeployableArtifactNames(
		b.Organizations[b.OrgCursor],
		b.Projects[b.ProjCursor],
		b.Components[b.CompCursor],
	)
}

// RenderProgress renders the selections made so far.
func (b BaseModel) RenderProgress() string {
	var progress strings.Builder
	progress.WriteString("Selected resources:\n")

	if len(b.Organizations) > 0 {
		progress.WriteString(fmt.Sprintf("- organization: %s\n", b.Organizations[b.OrgCursor]))
	}
	return progress.String()
}

// RenderOrgSelection returns a prompt for organization selection.
func (b BaseModel) RenderOrgSelection() string {
	return RenderListPrompt("Select organization:", b.Organizations, b.OrgCursor)
}

// RenderProjSelection returns a prompt for project selection.
func (b BaseModel) RenderProjSelection() string {
	return RenderListPrompt("Select project:", b.Projects, b.ProjCursor)
}

// RenderComponentSelection returns a prompt for component selection.
func (b BaseModel) RenderComponentSelection() string {
	return RenderListPrompt("Select component:", b.Components, b.CompCursor)
}

// RenderEnvironmentSelection returns a prompt for environment selection.
func (b BaseModel) RenderEnvironmentSelection() string {
	return RenderListPrompt("Select environment:", b.Environments, b.EnvCursor)
}

// RenderDeploymentTrackSelection returns a prompt for deployment track selection.
func (b BaseModel) RenderDeploymentTrackSelection() string {
	return RenderListPrompt("Select deployment track:", b.DeploymentTracks, b.DeploymentTrackCursor)
}

// RenderDeployableArtifactSelection returns a prompt for deployable artifact selection.
func (b BaseModel) RenderDeployableArtifactSelection() string {
	return RenderListPrompt("Select deployable artifact:", b.DeployableArtifacts, b.DeployableArtifactCursor)
}

// Reusable fetch functions

// FetchProjects retrieves project names for the currently selected organization.
func (b *BaseModel) FetchProjects() ([]string, error) {
	if b.OrgCursor >= len(b.Organizations) {
		return nil, fmt.Errorf("invalid organization index")
	}
	return util.GetProjectNames(b.Organizations[b.OrgCursor])
}

// FetchComponents retrieves component names for the currently selected organization and project.
func (b *BaseModel) FetchComponents() ([]string, error) {
	if b.OrgCursor >= len(b.Organizations) || b.ProjCursor >= len(b.Projects) {
		return nil, fmt.Errorf("invalid selection indices for components")
	}
	return util.GetComponentNames(b.Organizations[b.OrgCursor], b.Projects[b.ProjCursor])
}

// FetchEnvironments retrieves environment names for the currently selected organization.
// (Assumes environments depend only on the organization; adjust as needed.)
func (b *BaseModel) FetchEnvironments() ([]string, error) {
	if b.OrgCursor >= len(b.Organizations) {
		return nil, fmt.Errorf("invalid organization index")
	}
	return util.GetEnvironmentNames(b.Organizations[b.OrgCursor])
}

func (b *BaseModel) FetchDataPlanes() ([]string, error) {
	if b.OrgCursor >= len(b.Organizations) {
		return nil, fmt.Errorf("invalid organization index")
	}
	return util.GetDataPlaneNames(b.Organizations[b.OrgCursor])
}

func (b *BaseModel) FetchDeployments() ([]string, error) {
	if b.OrgCursor >= len(b.Organizations) ||
		b.ProjCursor >= len(b.Projects) ||
		b.CompCursor >= len(b.Components) {
		return nil, fmt.Errorf("invalid selection indices for deployments")
	}
	return util.GetDeploymentNames(
		b.Organizations[b.OrgCursor],
		b.Projects[b.ProjCursor],
		b.Components[b.CompCursor],
	)
}

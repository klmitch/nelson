// Copyright (c) 2021 Kevin L. Mitchell
//
// Licensed under the Apache License, Version 2.0 (the "License"); you
// may not use this file except in compliance with the License.  You
// may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied.  See the License for the specific language governing
// permissions and limitations under the License.

package nelson

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockICommand struct {
	mock.Mock
}

func (m *mockICommand) GetSummary() string {
	args := m.MethodCalled("GetSummary")

	return args.String(0)
}

func (m *mockICommand) GetDescription() string {
	args := m.MethodCalled("GetDescription")

	return args.String(0)
}

func (m *mockICommand) GetGroup() string {
	args := m.MethodCalled("GetGroup")

	return args.String(0)
}

func (m *mockICommand) GetSubcommands() map[string]ICommand {
	args := m.MethodCalled("GetSubcommands")

	if tmp := args.Get(0); tmp != nil {
		return tmp.(map[string]ICommand)
	}

	return nil
}

func (m *mockICommand) GetDefaults() interface{} {
	args := m.MethodCalled("GetDefaults")

	return args.Get(0)
}

func TestCommandImplementsICommand(t *testing.T) {
	assert.Implements(t, (*ICommand)(nil), &Command{})
}

func TestCommandGetSummary(t *testing.T) {
	obj := &Command{
		Summary: "some text",
	}

	result := obj.GetSummary()

	assert.Equal(t, "some text", result)
}

func TestCommandGetDescription(t *testing.T) {
	obj := &Command{
		Description: "some text",
	}

	result := obj.GetDescription()

	assert.Equal(t, "some text", result)
}

func TestCommandGetGroup(t *testing.T) {
	obj := &Command{
		Group: "some text",
	}

	result := obj.GetGroup()

	assert.Equal(t, "some text", result)
}

func TestCommandGetSubcommands(t *testing.T) {
	sub := &mockICommand{}
	obj := &Command{
		Subcommands: map[string]ICommand{
			"sub": sub,
		},
	}

	result := obj.GetSubcommands()

	assert.Equal(t, map[string]ICommand{
		"sub": sub,
	}, result)
}

func TestCommandGetDefaults(t *testing.T) {
	obj := &Command{
		Defaults: "defaults",
	}

	result := obj.GetDefaults()

	assert.Equal(t, "defaults", result)
}

func TestHiddenCommandImplementsICommand(t *testing.T) {
	assert.Implements(t, (*ICommand)(nil), &HiddenCommand{})
}

func TestHiddenCommandImplementsIWrapped(t *testing.T) {
	assert.Implements(t, (*IWrapped)(nil), &HiddenCommand{})
}

func TestHidden(t *testing.T) {
	cmd := &mockICommand{}

	result := Hidden(cmd)

	assert.Equal(t, &HiddenCommand{
		Wrapped: cmd,
	}, result)
}

func TestHiddenCommandGetSummary(t *testing.T) {
	cmd := &mockICommand{}
	cmd.On("GetSummary").Return("some text")
	obj := &HiddenCommand{
		Wrapped: cmd,
	}

	result := obj.GetSummary()

	assert.Equal(t, "some text", result)
	cmd.AssertExpectations(t)
}

func TestHiddenCommandGetDescription(t *testing.T) {
	cmd := &mockICommand{}
	cmd.On("GetDescription").Return("some text")
	obj := &HiddenCommand{
		Wrapped: cmd,
	}

	result := obj.GetDescription()

	assert.Equal(t, "some text", result)
	cmd.AssertExpectations(t)
}

func TestHiddenCommandGetGroup(t *testing.T) {
	cmd := &mockICommand{}
	cmd.On("GetGroup").Return("some text")
	obj := &HiddenCommand{
		Wrapped: cmd,
	}

	result := obj.GetGroup()

	assert.Equal(t, "some text", result)
	cmd.AssertExpectations(t)
}

func TestHiddenCommandGetSubcommands(t *testing.T) {
	subs := map[string]ICommand{
		"sub": &mockICommand{},
	}
	cmd := &mockICommand{}
	cmd.On("GetSubcommands").Return(subs)
	obj := &HiddenCommand{
		Wrapped: cmd,
	}

	result := obj.GetSubcommands()

	assert.Equal(t, subs, result)
	cmd.AssertExpectations(t)
}

func TestHiddenCommandGetDefaults(t *testing.T) {
	cmd := &mockICommand{}
	cmd.On("GetDefaults").Return("defaults")
	obj := &HiddenCommand{
		Wrapped: cmd,
	}

	result := obj.GetDefaults()

	assert.Equal(t, "defaults", result)
	cmd.AssertExpectations(t)
}

func TestHiddenUnwrap(t *testing.T) {
	cmd := &mockICommand{}
	obj := &HiddenCommand{
		Wrapped: cmd,
	}

	result := obj.Unwrap()

	assert.Same(t, cmd, result)
}

func TestDeprecatedCommandImplementsICommand(t *testing.T) {
	assert.Implements(t, (*ICommand)(nil), &DeprecatedCommand{})
}

func TestDeprecatedCommandImplementsIWrapped(t *testing.T) {
	assert.Implements(t, (*IWrapped)(nil), &DeprecatedCommand{})
}

func TestDeprecated(t *testing.T) {
	cmd := &mockICommand{}

	result := Deprecated(cmd, "alt")

	assert.Equal(t, &DeprecatedCommand{
		Wrapped:     cmd,
		Alternative: "alt",
	}, result)
}

func TestDeprecatedCommandGetSummary(t *testing.T) {
	cmd := &mockICommand{}
	cmd.On("GetSummary").Return("some text")
	obj := &DeprecatedCommand{
		Wrapped: cmd,
	}

	result := obj.GetSummary()

	assert.Equal(t, "some text", result)
	cmd.AssertExpectations(t)
}

func TestDeprecatedCommandGetDescription(t *testing.T) {
	cmd := &mockICommand{}
	cmd.On("GetDescription").Return("some text")
	obj := &DeprecatedCommand{
		Wrapped: cmd,
	}

	result := obj.GetDescription()

	assert.Equal(t, "some text", result)
	cmd.AssertExpectations(t)
}

func TestDeprecatedCommandGetGroup(t *testing.T) {
	cmd := &mockICommand{}
	cmd.On("GetGroup").Return("some text")
	obj := &DeprecatedCommand{
		Wrapped: cmd,
	}

	result := obj.GetGroup()

	assert.Equal(t, "some text", result)
	cmd.AssertExpectations(t)
}

func TestDeprecatedCommandGetSubcommands(t *testing.T) {
	subs := map[string]ICommand{
		"sub": &mockICommand{},
	}
	cmd := &mockICommand{}
	cmd.On("GetSubcommands").Return(subs)
	obj := &DeprecatedCommand{
		Wrapped: cmd,
	}

	result := obj.GetSubcommands()

	assert.Equal(t, subs, result)
	cmd.AssertExpectations(t)
}

func TestDeprecatedCommandGetDefaults(t *testing.T) {
	cmd := &mockICommand{}
	cmd.On("GetDefaults").Return("defaults")
	obj := &DeprecatedCommand{
		Wrapped: cmd,
	}

	result := obj.GetDefaults()

	assert.Equal(t, "defaults", result)
	cmd.AssertExpectations(t)
}

func TestDeprecatedUnwrap(t *testing.T) {
	cmd := &mockICommand{}
	obj := &DeprecatedCommand{
		Wrapped: cmd,
	}

	result := obj.Unwrap()

	assert.Same(t, cmd, result)
}

func TestAliasCommandImplementsICommand(t *testing.T) {
	assert.Implements(t, (*ICommand)(nil), &AliasCommand{})
}

func TestAliasCommandImplementsIWrapped(t *testing.T) {
	assert.Implements(t, (*IWrapped)(nil), &AliasCommand{})
}

func TestAlias(t *testing.T) {
	cmd := &mockICommand{}

	result := Alias(cmd)

	assert.Equal(t, &AliasCommand{
		Wrapped: cmd,
	}, result)
}

func TestAliasCommandGetSummary(t *testing.T) {
	cmd := &mockICommand{}
	cmd.On("GetSummary").Return("some text")
	obj := &AliasCommand{
		Wrapped: cmd,
	}

	result := obj.GetSummary()

	assert.Equal(t, "some text", result)
	cmd.AssertExpectations(t)
}

func TestAliasCommandGetDescription(t *testing.T) {
	cmd := &mockICommand{}
	cmd.On("GetDescription").Return("some text")
	obj := &AliasCommand{
		Wrapped: cmd,
	}

	result := obj.GetDescription()

	assert.Equal(t, "some text", result)
	cmd.AssertExpectations(t)
}

func TestAliasCommandGetGroup(t *testing.T) {
	cmd := &mockICommand{}
	cmd.On("GetGroup").Return("some text")
	obj := &AliasCommand{
		Wrapped: cmd,
	}

	result := obj.GetGroup()

	assert.Equal(t, "some text", result)
	cmd.AssertExpectations(t)
}

func TestAliasCommandGetSubcommands(t *testing.T) {
	subs := map[string]ICommand{
		"sub": &mockICommand{},
	}
	cmd := &mockICommand{}
	cmd.On("GetSubcommands").Return(subs)
	obj := &AliasCommand{
		Wrapped: cmd,
	}

	result := obj.GetSubcommands()

	assert.Equal(t, subs, result)
	cmd.AssertExpectations(t)
}

func TestAliasCommandGetDefaults(t *testing.T) {
	cmd := &mockICommand{}
	cmd.On("GetDefaults").Return("defaults")
	obj := &AliasCommand{
		Wrapped: cmd,
	}

	result := obj.GetDefaults()

	assert.Equal(t, "defaults", result)
	cmd.AssertExpectations(t)
}

func TestAliasUnwrap(t *testing.T) {
	cmd := &mockICommand{}
	obj := &AliasCommand{
		Wrapped: cmd,
	}

	result := obj.Unwrap()

	assert.Same(t, cmd, result)
}

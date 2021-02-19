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

// ICommand is an interface for a command type.  A declared command
// must implement this interface.
type ICommand interface {
	// GetSummary retrieves the command summary.
	GetSummary() string

	// GetDescription retrieves the command's full description.
	GetDescription() string

	// GetGroup retrieves the group name of the command.
	GetGroup() string

	// GetSubcommands retrieves subcommands for this command.
	GetSubcommands() map[string]ICommand

	// GetDefaults retrieves the defaults for arguments for this
	// command.
	GetDefaults() interface{}
}

// Command describes a command.  This type is intended for embedding,
// and implements the ICommand interface.
type Command struct {
	Summary     string              // The summary of the command
	Description string              // The full description of the command
	Group       string              // An optional group name for grouping related subcommands
	Subcommands map[string]ICommand // Subcommands of the command
	Defaults    interface{}         // Defaults for arguments
}

// GetSummary retrieves the command summary.
func (c *Command) GetSummary() string {
	return c.Summary
}

// GetDescription retrieves the command's full description.
func (c *Command) GetDescription() string {
	return c.Description
}

// GetGroup retrieves the group name of the command.
func (c *Command) GetGroup() string {
	return c.Group
}

// GetSubcommands retrieves subcommands for this command.
func (c *Command) GetSubcommands() map[string]ICommand {
	return c.Subcommands
}

// GetDefaults retrieves the defaults for arguments for this command.
func (c *Command) GetDefaults() interface{} {
	return c.Defaults
}

// IWrapped is an interface for commands that wrap other commands.  It
// allows the other commands to be unwrapped.
type IWrapped interface {
	// Unwrap returns the wrapped command.
	Unwrap() ICommand
}

// HiddenCommand wraps a command, causing it to be hidden from the
// usage message.
type HiddenCommand struct {
	Wrapped ICommand // Wrapped command
}

// Hidden wraps a command to indicate that it should be hidden from
// usage messages.
func Hidden(cmd ICommand) *HiddenCommand {
	return &HiddenCommand{
		Wrapped: cmd,
	}
}

// GetSummary retrieves the command summary.
func (c *HiddenCommand) GetSummary() string {
	return c.Wrapped.GetSummary()
}

// GetDescription retrieves the command's full description.
func (c *HiddenCommand) GetDescription() string {
	return c.Wrapped.GetDescription()
}

// GetGroup retrieves the group name of the command.
func (c *HiddenCommand) GetGroup() string {
	return c.Wrapped.GetGroup()
}

// GetSubcommands retrieves subcommands for this command.
func (c *HiddenCommand) GetSubcommands() map[string]ICommand {
	return c.Wrapped.GetSubcommands()
}

// GetDefaults retrieves the defaults for arguments for this command.
func (c *HiddenCommand) GetDefaults() interface{} {
	return c.Wrapped.GetDefaults()
}

// Unwrap returns the wrapped command.
func (c *HiddenCommand) Unwrap() ICommand {
	return c.Wrapped
}

// DeprecatedCommand wraps a command, causing it to be marked
// deprecated.
type DeprecatedCommand struct {
	Wrapped     ICommand // Wrapped command
	Alternative string   // Alternative command to use
}

// Deprecated wraps a command to indicate that it should be marked
// deprecated.
func Deprecated(cmd ICommand, alt string) *DeprecatedCommand {
	return &DeprecatedCommand{
		Wrapped:     cmd,
		Alternative: alt,
	}
}

// GetSummary retrieves the command summary.
func (c *DeprecatedCommand) GetSummary() string {
	return c.Wrapped.GetSummary()
}

// GetDescription retrieves the command's full description.
func (c *DeprecatedCommand) GetDescription() string {
	return c.Wrapped.GetDescription()
}

// GetGroup retrieves the group name of the command.
func (c *DeprecatedCommand) GetGroup() string {
	return c.Wrapped.GetGroup()
}

// GetSubcommands retrieves subcommands for this command.
func (c *DeprecatedCommand) GetSubcommands() map[string]ICommand {
	return c.Wrapped.GetSubcommands()
}

// GetDefaults retrieves the defaults for arguments for this command.
func (c *DeprecatedCommand) GetDefaults() interface{} {
	return c.Wrapped.GetDefaults()
}

// Unwrap returns the wrapped command.
func (c *DeprecatedCommand) Unwrap() ICommand {
	return c.Wrapped
}

// AliasCommand wraps a command and acts as an alias for that command.
type AliasCommand struct {
	Wrapped ICommand // Wrapped command
}

// Alias wraps a command to indicate that it is being used through an
// alias.
func Alias(cmd ICommand) *AliasCommand {
	return &AliasCommand{
		Wrapped: cmd,
	}
}

// GetSummary retrieves the command summary.
func (c *AliasCommand) GetSummary() string {
	return c.Wrapped.GetSummary()
}

// GetDescription retrieves the command's full description.
func (c *AliasCommand) GetDescription() string {
	return c.Wrapped.GetDescription()
}

// GetGroup retrieves the group name of the command.
func (c *AliasCommand) GetGroup() string {
	return c.Wrapped.GetGroup()
}

// GetSubcommands retrieves subcommands for this command.
func (c *AliasCommand) GetSubcommands() map[string]ICommand {
	return c.Wrapped.GetSubcommands()
}

// GetDefaults retrieves the defaults for arguments for this command.
func (c *AliasCommand) GetDefaults() interface{} {
	return c.Wrapped.GetDefaults()
}

// Unwrap returns the wrapped command.
func (c *AliasCommand) Unwrap() ICommand {
	return c.Wrapped
}

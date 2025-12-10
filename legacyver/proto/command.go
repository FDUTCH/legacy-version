package proto

import (
	"github.com/akmalfairuz/legacy-version/internal/typeconf"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"math"
)

// Command holds the data that a command requires to be shown to a player client-side. The command is shown in
// the /help command and auto-completed using this data.
type Command struct {
	// Name is the name of the command. The command may be executed using this name, and will be shown in the
	// /help list with it. It currently seems that the client crashes if the Name contains uppercase letters.
	Name string
	// Description is the description of the command. It is shown in the /help list and when starting to write
	// a command.
	Description string
	// Flags is a combination of flags not currently known. Leaving the Flags field empty appears to work.
	Flags uint16
	// PermissionLevel is the command permission level that the player required to execute this command. The
	// field no longer seems to serve a purpose, as the client does not handle the execution of commands
	// anymore: The permissions should be checked server-side.
	PermissionLevel string
	// AliasesOffset is the offset to a CommandEnum that holds the values that
	// should be used as aliases for this command.
	AliasesOffset uint32
	// ChainedSubcommandOffsets is a slice of offsets that all point to a different ChainedSubcommand from the
	// ChainedSubcommands slice in the AvailableCommands packet.
	ChainedSubcommandOffsets []uint32
	// Overloads is a list of command overloads that specify the ways in which a command may be executed. The
	// overloads may be completely different.
	Overloads []protocol.CommandOverload
}

func (c *Command) Marshal(r protocol.IO) {
	r.String(&c.Name)
	r.String(&c.Description)
	r.Uint16(&c.Flags)
	if IsProtoGTE(r, ID898) {
		r.String(&c.PermissionLevel)
	} else {
		permLevel := uint8(0)
		r.Uint8(&permLevel)
		c.PermissionLevel = "any"
	}
	r.Uint32(&c.AliasesOffset)
	if IsProtoGTE(r, ID898) {
		protocol.FuncSlice(r, &c.ChainedSubcommandOffsets, r.Uint32)
	} else {
		offsets := typeconf.SliceIntToSliceInt[uint32, uint16](c.ChainedSubcommandOffsets)
		protocol.FuncSlice(r, &offsets, r.Uint16)
		c.ChainedSubcommandOffsets = typeconf.SliceIntToSliceInt[uint16, uint32](offsets)
	}
	protocol.Slice(r, &c.Overloads)
}

func (c *Command) ToLatest() protocol.Command {
	return protocol.Command{
		Name:                     c.Name,
		Description:              c.Description,
		Flags:                    c.Flags,
		PermissionLevel:          c.PermissionLevel,
		AliasesOffset:            c.AliasesOffset,
		ChainedSubcommandOffsets: c.ChainedSubcommandOffsets,
		Overloads:                c.Overloads,
	}
}

func (c *Command) FromLatest(latest protocol.Command) Command {
	c.Name = latest.Name
	c.Description = latest.Description
	c.Flags = latest.Flags
	c.PermissionLevel = latest.PermissionLevel
	c.AliasesOffset = latest.AliasesOffset
	c.ChainedSubcommandOffsets = latest.ChainedSubcommandOffsets
	c.Overloads = latest.Overloads
	return *c
}

// CommandEnum represents an enum in a command usage. The enum typically has a type and a set of options that
// are valid. A value that is not one of the options results in a failure during execution.
type CommandEnum struct {
	// Type is the type of the command enum. The type will show up in the command usage as the type of the
	// argument if it has a certain amount of arguments, or when Options is set to true in the
	// command holding the enum.
	Type string
	// ValueIndices holds a list of indices that point to the EnumValues slice in the
	// AvailableCommandsPacket. These represent the options of the enum.
	ValueIndices []uint32
}

// CommandEnumContext holds context required for encoding command enums.
type CommandEnumContext struct {
	EnumValues []string
}

// Marshal encodes/decodes a CommandEnum.
func (ctx CommandEnumContext) Marshal(r protocol.IO, x *CommandEnum) {
	r.String(&x.Type)
	if IsProtoGTE(r, ID898) {
		protocol.FuncSlice(r, &x.ValueIndices, r.Uint32)
	} else {
		protocol.FuncIOSlice(r, &x.ValueIndices, ctx.enumOption)
	}
}

// enumOption writes/reads a command enum option as a byte/uint16/uint32,
// depending on the amount of enum values.
func (ctx CommandEnumContext) enumOption(r protocol.IO, opt *uint32) {
	n := len(ctx.EnumValues)
	switch {
	case n <= math.MaxUint8:
		val := byte(*opt)
		r.Uint8(&val)
		*opt = uint32(val)
	case n <= math.MaxUint16:
		val := uint16(*opt)
		r.Uint16(&val)
		*opt = uint32(val)
	default:
		r.Uint32(opt)
	}
}

// ToLatest ...
func (ce *CommandEnum) ToLatest() protocol.CommandEnum {
	return protocol.CommandEnum{
		Type:         ce.Type,
		ValueIndices: ce.ValueIndices,
	}
}

// FromLatest ...
func (ce *CommandEnum) FromLatest(latest protocol.CommandEnum) CommandEnum {
	ce.Type = latest.Type
	ce.ValueIndices = latest.ValueIndices
	return *ce
}

// ChainedSubcommand represents a subcommand that can have chained commands, such as /execute which allows you to run
// another command as another entity or at a different position etc.
type ChainedSubcommand struct {
	// Name is the name of the chained subcommand and shows up in the list as a regular subcommand enum.
	Name string
	// Values contains the index and parameter type of the chained subcommand.
	Values []ChainedSubcommandValue
}

func (x *ChainedSubcommand) Marshal(r protocol.IO) {
	r.String(&x.Name)
	protocol.Slice(r, &x.Values)
}

// ToLatest ...
func (x *ChainedSubcommand) ToLatest() protocol.ChainedSubcommand {
	values := make([]protocol.ChainedSubcommandValue, len(x.Values))
	for i, v := range x.Values {
		values[i] = v.ToLatest()
	}
	return protocol.ChainedSubcommand{
		Name:   x.Name,
		Values: values,
	}
}

// FromLatest ...
func (x *ChainedSubcommand) FromLatest(latest protocol.ChainedSubcommand) ChainedSubcommand {
	x.Name = latest.Name
	x.Values = make([]ChainedSubcommandValue, len(latest.Values))
	for i, v := range latest.Values {
		x.Values[i] = (&ChainedSubcommandValue{}).FromLatest(v)
	}
	return *x
}

// ChainedSubcommandValue represents the value for a chained subcommand argument.
type ChainedSubcommandValue struct {
	// Index is the index of the argument in the ChainedSubcommandValues slice from the AvailableCommands packet. This is
	// then used to set the type specified by the Value field below.
	Index uint32
	// Value is a combination of the flags above and specified the type of argument. Unlike regular parameter types,
	// this should NOT contain any of the special flags (valid, enum, suffixed or soft enum) but only the basic types.
	Value uint32
}

func (x *ChainedSubcommandValue) Marshal(r protocol.IO) {
	if IsProtoGTE(r, ID898) {
		r.Varuint32(&x.Index)
		r.Varuint32(&x.Value)
	} else {
		v := uint16(x.Index)
		r.Uint16(&v)
		x.Index = uint32(v)
		vv := uint16(x.Value)
		r.Uint16(&vv)
		x.Value = uint32(vv)
	}
}

// ToLatest ...
func (x *ChainedSubcommandValue) ToLatest() protocol.ChainedSubcommandValue {
	return protocol.ChainedSubcommandValue{
		Index: x.Index,
		Value: x.Value,
	}
}

// FromLatest ...
func (x *ChainedSubcommandValue) FromLatest(latest protocol.ChainedSubcommandValue) ChainedSubcommandValue {
	x.Index = latest.Index
	x.Value = latest.Value
	return *x
}

package resolver

import (
	"fmt"
	"slices"

	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/chips"
	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/errors"
	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/lexer"
	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/parser"
)

const MAX_NUMBER_OF_IOS = 80
const MIN_IO_WIDTH = 1
const MAX_IO_WIDTH = 1024

const MAX_NUMBER_OF_PARTS = 100

type Resolver struct {
	chd                       *parser.ParsedChipDefinition
	chipFileName              string
	hdls                      map[string]string
	resolvedChipDef           *ResolvedChipDefinition
	resolvedUsedChipDefs      map[string]*ResolvedChipDefinition
	chipOutputSignalCoverages map[string]signalCoverage
	partInputPinCoverages     map[int]map[string]pinCoverage
}

type connection struct {
	PartName  string
	PartIndex int
	PartIO    chips.IO
	Pin       parser.Pin
	Signal    parser.Signal
}

type signalCoverage map[int]bool
type pinCoverage map[int]bool

func New(chd *parser.ParsedChipDefinition, chipFileName string, hdls map[string]string) *Resolver {
	resolvedChipDef := &ResolvedChipDefinition{
		Inputs:          make(map[string]chips.IO),
		Outputs:         make(map[string]chips.IO),
		InternalSignals: make(map[string]InternalSignal),
	}
	resolvedUsedChipDefs := make(map[string]*ResolvedChipDefinition)
	chipOutputSignalCoverages := make(map[string]signalCoverage)
	partInputPinCoverages := make(map[int]map[string]pinCoverage)
	return &Resolver{
		chd:                       chd,
		chipFileName:              chipFileName,
		resolvedChipDef:           resolvedChipDef,
		resolvedUsedChipDefs:      resolvedUsedChipDefs,
		chipOutputSignalCoverages: chipOutputSignalCoverages,
		partInputPinCoverages:     partInputPinCoverages,
		hdls:                      hdls,
	}
}

func (r *Resolver) Resolve(resolvedChipNames []string, resolvingChipNames []string) (
	*ResolvedChipDefinition,
	map[string]*ResolvedChipDefinition,
	error,
) {
	resolvingChipNames = append(resolvingChipNames, r.chipFileName)

	if err := r.resolveChipName(); err != nil {
		return nil, map[string]*ResolvedChipDefinition{}, err
	}

	if err := r.resolveIO(); err != nil {
		return nil, map[string]*ResolvedChipDefinition{}, err
	}

	if err := r.validateNumberOfParts(); err != nil {
		return nil, map[string]*ResolvedChipDefinition{}, err
	}

	if err := r.resolveUsedChips(resolvedChipNames, resolvingChipNames); err != nil {
		return nil, map[string]*ResolvedChipDefinition{}, err
	}

	for _, part := range r.chd.Parts {
		r.resolvedChipDef.Parts = append(r.resolvedChipDef.Parts, Part{
			Name: part.Name,
		})
	}

	inputConnections, outputConnections, err := r.getInputAndOutputConnections()

	if err != nil {
		return nil, map[string]*ResolvedChipDefinition{}, err
	}

	if err := r.resolveOutputConnections(outputConnections); err != nil {
		return nil, map[string]*ResolvedChipDefinition{}, err
	}

	if err := r.resolveInputConnections(inputConnections); err != nil {
		return nil, map[string]*ResolvedChipDefinition{}, err
	}

	return r.resolvedChipDef, r.resolvedUsedChipDefs, nil
}

func (r *Resolver) resolveChipName() error {
	if r.chd.ChipName.Name != r.chipFileName {
		return r.newResolutionError(
			"File name does not match the chip name",
			r.chd.ChipName.Loc.Line, r.chd.ChipName.Loc.Column,
		)
	}

	r.resolvedChipDef.Name = r.chd.ChipName.Name
	return nil
}

func (r *Resolver) resolveIO() error {
	// validate number of inputs and outputs
	if len(r.chd.Inputs) > MAX_NUMBER_OF_IOS {
		return r.newResolutionError(
			"Number of inputs exceeds maximum allowed",
			r.chd.Inputs[len(r.chd.Inputs)-1].Loc.Line, r.chd.Inputs[len(r.chd.Inputs)-1].Loc.Column,
		)
	}

	if len(r.chd.Outputs) > MAX_NUMBER_OF_IOS {
		return r.newResolutionError(
			"Number of outputs exceeds maximum allowed",
			r.chd.Outputs[len(r.chd.Outputs)-1].Loc.Line, r.chd.Outputs[len(r.chd.Outputs)-1].Loc.Column,
		)
	}

	// validate input widths
	for _, input := range r.chd.Inputs {
		if input.Width < MIN_IO_WIDTH || input.Width > MAX_IO_WIDTH {
			return r.newResolutionError(
				fmt.Sprintf("Input '%s' width out of bounds", input.Name),
				input.Loc.Line, input.Loc.Column,
			)
		}
	}

	// validate output widths
	for _, output := range r.chd.Outputs {
		if output.Width < MIN_IO_WIDTH || output.Width > MAX_IO_WIDTH {
			return r.newResolutionError(
				fmt.Sprintf("Output '%s' width out of bounds", output.Name),
				output.Loc.Line, output.Loc.Column,
			)
		}
	}

	// check for duplicates
	seenIOs := make(map[string]bool)
	for _, input := range r.chd.Inputs {
		if seenIOs[input.Name] {
			return r.newResolutionError(
				fmt.Sprintf("Duplicate input name '%s'", input.Name),
				input.Loc.Line, input.Loc.Column,
			)
		}
		seenIOs[input.Name] = true
	}

	for _, output := range r.chd.Outputs {
		if seenIOs[output.Name] {
			return r.newResolutionError(
				fmt.Sprintf("Duplicate output name '%s'", output.Name),
				output.Loc.Line, output.Loc.Column,
			)
		}
		seenIOs[output.Name] = true
	}

	for _, input := range r.chd.Inputs {
		r.resolvedChipDef.Inputs[input.Name] = chips.IO{Width: input.Width}
	}

	for _, output := range r.chd.Outputs {
		r.resolvedChipDef.Outputs[output.Name] = chips.IO{Width: output.Width}
	}

	return nil
}

func (r *Resolver) validateNumberOfParts() error {
	if len(r.chd.Parts) > MAX_NUMBER_OF_PARTS {
		return r.newResolutionError(
			"Number of parts exceeds maximum allowed",
			r.chd.Parts[len(r.chd.Parts)-1].Loc.Line, r.chd.Parts[len(r.chd.Parts)-1].Loc.Column,
		)
	}
	return nil
}

func (r *Resolver) resolveUsedChips(resolvedChipNames []string, resolvingChipNames []string) error {
	usedChipNames := r.chd.GetUsedChipNames()
	_, usedCustomChipNames, err := r.groupUsedChipNames(usedChipNames)
	if err != nil {
		return err
	}

	for _, chipName := range usedCustomChipNames {
		if slices.Contains(resolvingChipNames, chipName) {
			return r.newResolutionError(
				"Circular dependency detected: "+fmt.Sprintf("%v", append(resolvingChipNames, chipName)),
				0, 0, // TODO: improve with actual location and file name,
			)
		}

		if slices.Contains(resolvedChipNames, chipName) {
			continue
		}

		chd, err := r.lexAndParseHDL(chipName)
		if err != nil {
			return err
		}

		r2 := New(chd, chipName, r.hdls)
		resolvedChipDef, resolvedUsedChipDefs, err := r2.Resolve(resolvedChipNames, resolvingChipNames)
		if err != nil {
			return err
		}

		if !slices.Contains(resolvedChipNames, chipName) {
			resolvedChipNames = append(resolvedChipNames, chipName)
			r.resolvedUsedChipDefs[chipName] = resolvedChipDef
		}

		for name, def := range resolvedUsedChipDefs {
			if !slices.Contains(resolvedChipNames, name) {
				resolvedChipNames = append(resolvedChipNames, name)
				r.resolvedUsedChipDefs[name] = def
			}
		}
	}
	return nil
}

func (r *Resolver) groupUsedChipNames(usedChipNames []string) (builtInChipNames []string, customChipNames []string, err error) {
	for _, name := range usedChipNames {
		_, ok := chips.BuiltInChips[name]
		if ok {
			builtInChipNames = append(builtInChipNames, name)
			continue
		}

		if _, ok := r.hdls[name]; ok {
			customChipNames = append(customChipNames, name)
			continue
		}

		return nil, nil, r.newResolutionError(
			fmt.Sprintf("Used chip '%s' is neither a built-in chip nor a custom chip", name),
			0, 0, // TODO: improve with actual location and file name
		)
	}

	return builtInChipNames, customChipNames, nil
}

func (r *Resolver) lexAndParseHDL(chipName string) (*parser.ParsedChipDefinition, error) {
	l := lexer.New(r.hdls[chipName])
	ts, err := l.Tokenize()
	if err != nil {
		return nil, err
	}

	p := parser.New(ts)
	chd, err := p.ParseChipDefinition()
	if err != nil {
		return nil, err
	}

	return chd, nil
}

func (r *Resolver) getInputAndOutputConnections() (inputConnections []connection, outputConnections []connection, err error) {
	for idx, part := range r.chd.Parts {
		var partInputs map[string]chips.IO
		var partOutputs map[string]chips.IO
		if partDef, isCustomChip := r.resolvedUsedChipDefs[part.Name]; isCustomChip {
			// part is a custom chip
			partInputs = partDef.Inputs
			partOutputs = partDef.Outputs
		} else {
			// part is a built-in chip
			// don't need to validate that it exists, as it was already validated in groupUsedChipNames
			builtInChip := chips.BuiltInChips[part.Name]
			partInputs = builtInChip.Inputs
			partOutputs = builtInChip.Outputs
		}

		for _, conn := range part.Connections {
			if input, isInput := partInputs[conn.Pin.Name]; isInput {
				// the pin of the connection is an input of the used part
				inputConnections = append(inputConnections, connection{
					PartName:  part.Name,
					PartIndex: idx,
					PartIO:    input,
					Pin:       conn.Pin,
					Signal:    conn.Signal,
				})
				continue
			}
			if output, isOutput := partOutputs[conn.Pin.Name]; isOutput {
				// the pin of the connection is an output of the used part
				outputConnections = append(outputConnections, connection{
					PartName:  part.Name,
					PartIndex: idx,
					PartIO:    output,
					Pin:       conn.Pin,
					Signal:    conn.Signal,
				})
				continue
			}
			// else, the pin is not found in the used part's inputs or outputs
			return nil, nil, r.newResolutionError(
				fmt.Sprintf("Pin '%s' not found in part '%s'", conn.Pin.Name, part.Name),
				conn.Pin.Loc.Line, conn.Pin.Loc.Column,
			)
		}
	}

	return inputConnections, outputConnections, nil
}

func (r *Resolver) newResolutionError(message string, line, column int) *errors.ResolutionError {
	return &errors.ResolutionError{
		Message: message,
		Line:    line,
		Column:  column,
		File:    r.chipFileName,
	}
}

func (r *Resolver) resolveOutputConnections(connections []connection) error {
	for _, conn := range connections {
		part := &r.resolvedChipDef.Parts[conn.PartIndex]
		resolvedConn := Connection{}
		err := r.resolveOutputConnectionPin(&resolvedConn, conn)
		if err != nil {
			return err
		}
		err = r.resolveOutputConnectionSignal(&resolvedConn, conn)
		if err != nil {
			return err
		}
		part.OutputConnections = append(part.OutputConnections, resolvedConn)
	}
	return nil
}

func (r *Resolver) resolveOutputConnectionPin(resolvedConnection *Connection, conn connection) error {
	pin := Pin{}
	pin.Name = conn.Pin.Name
	if conn.Pin.Range.IsSpecified {
		if conn.Pin.Range.End >= conn.PartIO.Width || conn.Pin.Range.Start > conn.Pin.Range.End {
			return r.newResolutionError(
				fmt.Sprintf("Pin '%s' range out of bounds for part '%s'", conn.Pin.Name, conn.PartName),
				conn.Pin.Range.Loc.Line, conn.Pin.Range.Loc.Column,
			)
		}
		pin.Range = Range{Start: conn.Pin.Range.Start, End: conn.Pin.Range.End}
	} else {
		pin.Range = Range{Start: 0, End: conn.PartIO.Width - 1}
	}
	resolvedConnection.Pin = pin
	return nil
}

func (r *Resolver) resolveOutputConnectionSignal(resolvedConn *Connection, conn connection) error {
	signal := Signal{}
	signal.Name = conn.Signal.Name
	if conn.Signal.Range.IsSpecified {
		if _, isChipOutput := r.resolvedChipDef.Outputs[signal.Name]; !isChipOutput {
			return r.newResolutionError(
				fmt.Sprintf("Internal output signal '%s' cannot be partially defined", conn.Signal.Name),
				conn.Signal.Loc.Line, conn.Signal.Loc.Column,
			)
		}
		chipOutput := r.resolvedChipDef.Outputs[signal.Name]
		// the signal is an output signal of the chip
		// and here the user defines part of the output signal
		if conn.Signal.Range.Start > conn.Signal.Range.End {
			return r.newResolutionError(
				fmt.Sprintf("Signal '%s' range is invalid", conn.Signal.Name),
				conn.Signal.Range.Loc.Line, conn.Signal.Range.Loc.Column,
			)
		}

		if conn.Signal.Range.End >= chipOutput.Width {
			return r.newResolutionError(
				fmt.Sprintf("Signal '%s' range out of bounds", conn.Signal.Name),
				conn.Signal.Range.Loc.Line, conn.Signal.Range.Loc.Column,
			)
		}

		if conn.Signal.Range.End-conn.Signal.Range.Start != resolvedConn.Pin.Range.End-resolvedConn.Pin.Range.Start {
			return r.newResolutionError(
				fmt.Sprintf("Signal '%s' range width does not match pin '%s' range width", conn.Signal.Name, conn.Pin.Name),
				conn.Signal.Range.Loc.Line, conn.Signal.Range.Loc.Column,
			)
		}

		signal.Range = Range{Start: conn.Signal.Range.Start, End: conn.Signal.Range.End}
		ok := addRangeToSignalCoverages(r.chipOutputSignalCoverages, signal.Name, signal.Range)
		if !ok {
			return r.newResolutionError(
				fmt.Sprintf("Signal '%s' range overlaps with existing ranges", conn.Signal.Name),
				conn.Signal.Range.Loc.Line, conn.Signal.Range.Loc.Column,
			)
		}

		resolvedConn.Signal = signal
		return nil
	}
	if chipOutput, isChipOutput := r.resolvedChipDef.Outputs[signal.Name]; isChipOutput {
		rng := Range{Start: 0, End: chipOutput.Width - 1}
		if resolvedConn.Pin.Range.End-resolvedConn.Pin.Range.Start != rng.End-rng.Start {
			return r.newResolutionError(
				fmt.Sprintf("Signal '%s' width does not match pin '%s' width", conn.Signal.Name, resolvedConn.Pin.Name),
				conn.Signal.Loc.Line, conn.Signal.Loc.Column,
			)
		}
		ok := addRangeToSignalCoverages(r.chipOutputSignalCoverages, signal.Name, rng)
		if !ok {
			return r.newResolutionError(
				fmt.Sprintf("Signal '%s' range overlaps with existing ranges", conn.Signal.Name),
				conn.Signal.Loc.Line, conn.Signal.Loc.Column,
			)
		}
		signal.Range = rng
	} else {
		signal.Range = Range{Start: 0, End: conn.PartIO.Width - 1}
		if _, exists := r.resolvedChipDef.InternalSignals[signal.Name]; !exists {
			r.resolvedChipDef.InternalSignals[signal.Name] = InternalSignal{
				Width: conn.PartIO.Width,
			}
		} else {
			return r.newResolutionError(
				fmt.Sprintf("Internal signal '%s' already defined", signal.Name),
				conn.Signal.Loc.Line, conn.Signal.Loc.Column,
			)
		}
	}

	resolvedConn.Signal = signal
	return nil
}

func (r *Resolver) resolveInputConnections(connections []connection) error {
	for _, conn := range connections {
		part := &r.resolvedChipDef.Parts[conn.PartIndex]
		resolvedConn := Connection{}
		err := r.resolveInputConnectionPin(&resolvedConn, conn)
		if err != nil {
			return err
		}
		err = r.resolveInputConnectionSignal(&resolvedConn, conn)
		if err != nil {
			return err
		}
		part.InputConnections = append(part.InputConnections, resolvedConn)
	}
	return nil
}

func (r *Resolver) resolveInputConnectionPin(resolvedConnection *Connection, conn connection) error {
	pin := Pin{}
	pin.Name = conn.Pin.Name
	if conn.Pin.Range.IsSpecified {
		if conn.Pin.Range.End >= conn.PartIO.Width || conn.Pin.Range.Start > conn.Pin.Range.End {
			return r.newResolutionError(
				fmt.Sprintf("Pin '%s' range out of bounds for part '%s'", conn.Pin.Name, conn.PartName),
				conn.Pin.Range.Loc.Line, conn.Pin.Range.Loc.Column,
			)
		}
		pin.Range = Range{Start: conn.Pin.Range.Start, End: conn.Pin.Range.End}
		ok := addRangeToPinCoverages(r.partInputPinCoverages, conn.PartIndex, conn.Pin.Name, pin.Range)
		if !ok {
			return r.newResolutionError(
				fmt.Sprintf("Pin '%s' range overlaps with existing ranges", conn.Signal.Name),
				conn.Pin.Range.Loc.Line, conn.Pin.Range.Loc.Column,
			)
		}
	} else {
		pin.Range = Range{Start: 0, End: conn.PartIO.Width - 1}
		ok := addRangeToPinCoverages(r.partInputPinCoverages, conn.PartIndex, conn.Pin.Name, pin.Range)
		if !ok {
			return r.newResolutionError(
				fmt.Sprintf("Pin '%s' range overlaps with existing ranges", conn.Signal.Name),
				conn.Pin.Range.Loc.Line, conn.Pin.Range.Loc.Column,
			)
		}
	}
	resolvedConnection.Pin = pin
	return nil
}

func (r *Resolver) resolveInputConnectionSignal(resolvedConn *Connection, conn connection) error {
	signal := Signal{}
	signal.Name = conn.Signal.Name

	internalSignal, isInternalSignal := r.resolvedChipDef.InternalSignals[signal.Name]
	inputIO, isChipInput := r.resolvedChipDef.Inputs[signal.Name]

	_ = internalSignal
	_ = inputIO

	if !isInternalSignal && !isChipInput {
		return r.newResolutionError(
			fmt.Sprintf("Signal '%s' is neither an internal signal nor a chip input", conn.Signal.Name),
			conn.Signal.Loc.Line, conn.Signal.Loc.Column,
		)
	}

	if !conn.Signal.Range.IsSpecified {
		if isInternalSignal {
			rng := Range{Start: 0, End: internalSignal.Width - 1}
			if resolvedConn.Pin.Range.End-resolvedConn.Pin.Range.Start != rng.End-rng.Start {
				return r.newResolutionError(
					fmt.Sprintf("Signal '%s' width does not match pin '%s' width", conn.Signal.Name, resolvedConn.Pin.Name),
					conn.Signal.Loc.Line, conn.Signal.Loc.Column,
				)
			}
			signal.Range = rng
			resolvedConn.Signal = signal
			return nil
		}
		rng := Range{Start: 0, End: inputIO.Width - 1}
		if resolvedConn.Pin.Range.End-resolvedConn.Pin.Range.Start != rng.End-rng.Start {
			return r.newResolutionError(
				fmt.Sprintf("Signal '%s' width does not match pin '%s' width", conn.Signal.Name, resolvedConn.Pin.Name),
				conn.Signal.Loc.Line, conn.Signal.Loc.Column,
			)
		}
		signal.Range = rng
		resolvedConn.Signal = signal
		return nil
	}

	rng := Range{Start: conn.Signal.Range.Start, End: conn.Signal.Range.End}
	if rng.Start > rng.End {
		return r.newResolutionError(
			fmt.Sprintf("Signal '%s' range is invalid", conn.Signal.Name),
			conn.Signal.Range.Loc.Line, conn.Signal.Range.Loc.Column,
		)
	}

	if isInternalSignal {
		if rng.End >= internalSignal.Width {
			return r.newResolutionError(
				fmt.Sprintf("Signal '%s' range out of bounds", conn.Signal.Name),
				conn.Signal.Range.Loc.Line, conn.Signal.Range.Loc.Column,
			)
		}
		if rng.End-rng.Start != resolvedConn.Pin.Range.End-resolvedConn.Pin.Range.Start {
			return r.newResolutionError(
				fmt.Sprintf("Signal '%s' range width does not match pin '%s' range width", conn.Signal.Name, conn.Pin.Name),
				conn.Signal.Range.Loc.Line, conn.Signal.Range.Loc.Column,
			)
		}
		signal.Range = rng
		resolvedConn.Signal = signal
		return nil
	}

	if rng.End >= inputIO.Width {
		return r.newResolutionError(
			fmt.Sprintf("Signal '%s' range out of bounds", conn.Signal.Name),
			conn.Signal.Range.Loc.Line, conn.Signal.Range.Loc.Column,
		)
	}
	if rng.End-rng.Start != resolvedConn.Pin.Range.End-resolvedConn.Pin.Range.Start {
		return r.newResolutionError(
			fmt.Sprintf("Signal '%s' range width does not match pin '%s' range width", conn.Signal.Name, conn.Pin.Name),
			conn.Signal.Range.Loc.Line, conn.Signal.Range.Loc.Column,
		)
	}
	signal.Range = rng
	resolvedConn.Signal = signal
	return nil
}

func addRangeToSignalCoverages(signalCoverages map[string]signalCoverage, signalName string, rng Range) bool {
	if _, ok := signalCoverages[signalName]; !ok {
		signalCoverages[signalName] = make(signalCoverage)
	}

	coverage := signalCoverages[signalName]
	for i := rng.Start; i <= rng.End; i++ {
		if coverage[i] {
			return false
		}
		coverage[i] = true
	}
	signalCoverages[signalName] = coverage
	return true
}

func addRangeToPinCoverages(partInputPinCoverages map[int]map[string]pinCoverage, partIndex int, pinName string, rng Range) bool {
	if _, ok := partInputPinCoverages[partIndex]; !ok {
		partInputPinCoverages[partIndex] = make(map[string]pinCoverage)
	}

	pinCoverages := partInputPinCoverages[partIndex]

	if _, ok := pinCoverages[pinName]; !ok {
		pinCoverages[pinName] = make(pinCoverage)
	}

	coverage := pinCoverages[pinName]
	for i := rng.Start; i <= rng.End; i++ {
		if coverage[i] {
			return false
		}
		coverage[i] = true
	}
	pinCoverages[pinName] = coverage
	return true
}

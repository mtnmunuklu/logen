package modifiers

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strings"
	"unicode/utf16"
)

// syntheticDataGenerator is a global instance of SyntheticDataGenerator.
var syntheticDataGenerator = NewSyntheticDataGenerator()

func GetComparator(modifiers ...string) (ComparatorFunc, error) {
	return getComparator(Comparators, modifiers...)
}

func GetComparatorCaseSensitive(modifiers ...string) (ComparatorFunc, error) {
	return getComparator(ComparatorsCaseSensitive, modifiers...)
}

func getComparator(comparators map[string]Comparator, modifiers ...string) (ComparatorFunc, error) {
	if len(modifiers) == 0 {
		return baseComparator{}.Alters, nil
	}

	// A valid sequence of modifiers is ([ValueModifier]*)[Comparator]?
	// If a comparator is specified, it must be in the last position and cannot be succeeded by any other modifiers
	// If no comparator is specified, the default comparator is used
	var valueModifiers []ValueModifier
	var comparator Comparator
	for i, modifier := range modifiers {
		comparatorModifier := comparators[modifier]
		valueModifier := ValueModifiers[modifier]
		switch {
		// Validate correctness
		case comparatorModifier == nil && valueModifier == nil:
			return nil, fmt.Errorf("unknown modifier %s", modifier)
		case i < len(modifiers)-1 && comparators[modifier] != nil:
			return nil, fmt.Errorf("comparator modifier %s must be the last modifier", modifier)

		// Build up list of modifiers
		case valueModifier != nil:
			valueModifiers = append(valueModifiers, valueModifier)
		case comparatorModifier != nil:
			comparator = comparatorModifier
		}
	}
	if comparator == nil {
		comparator = baseComparator{}
	}

	return func(field, value any) (string, error) {
		var err error
		for _, modifier := range valueModifiers {
			value, err = modifier.Modify(value)
			if err != nil {
				return "", err
			}
		}

		return comparator.Alters(field, value)
	}, nil
}

type Comparator interface {
	Alters(field any, value any) (string, error)
}

type ComparatorFunc func(field, value any) (string, error)

// ValueModifier modifies the expected value before it is passed to the comparator.
// For example, the `base64` modifier converts the expected value to base64.
type ValueModifier interface {
	Modify(value any) (any, error)
}

var Comparators = map[string]Comparator{
	"contains":   contains{generator: syntheticDataGenerator},
	"endswith":   endswith{generator: syntheticDataGenerator},
	"startswith": startswith{generator: syntheticDataGenerator},
	"re":         re{generator: syntheticDataGenerator},
	"cidr":       cidr{generator: syntheticDataGenerator},
	"gt":         gt{generator: syntheticDataGenerator},
	"gte":        gte{generator: syntheticDataGenerator},
	"lt":         lt{generator: syntheticDataGenerator},
	"lte":        lte{generator: syntheticDataGenerator},
}

var ComparatorsCaseSensitive = map[string]Comparator{
	"contains":   containsCS{generator: syntheticDataGenerator},
	"endswith":   endswithCS{generator: syntheticDataGenerator},
	"startswith": startswithCS{generator: syntheticDataGenerator},
	"re":         re{generator: syntheticDataGenerator},
	"cidr":       cidr{generator: syntheticDataGenerator},
	"gt":         gt{generator: syntheticDataGenerator},
	"gte":        gte{generator: syntheticDataGenerator},
	"lt":         lt{generator: syntheticDataGenerator},
	"lte":        lte{generator: syntheticDataGenerator},
}

var ValueModifiers = map[string]ValueModifier{
	"base64": b64{},
	"wide":   wide{},
}

type baseComparator struct{}

func (baseComparator) Alters(field, value any) (string, error) {
	switch {
	case field == nil && value == "null":
		return "", nil
	default:
		// The Sigma spec defines that by default comparisons are case-insensitive
		return fmt.Sprintf("%v equal '%v'", strings.ToLower(coerceString(field)), strings.ToLower(coerceString(value))), nil
	}
}

type contains struct {
	generator *SyntheticDataGenerator
}

func (c contains) Alters(field, value any) (string, error) {
	syntheticValue := c.generator.GenerateSyntheticValue(coerceString(value), "contains")
	return fmt.Sprintf("%v contains '%v'", strings.ToLower(coerceString(field)), strings.ToLower(syntheticValue)), nil
}

type endswith struct {
	generator *SyntheticDataGenerator
}

func (e endswith) Alters(field, value any) (string, error) {
	syntheticValue := e.generator.GenerateSyntheticValue(coerceString(value), "endswith")
	return fmt.Sprintf("%v endswith '%v'", strings.ToLower(coerceString(field)), strings.ToLower(syntheticValue)), nil
}

type startswith struct {
	generator *SyntheticDataGenerator
}

func (s startswith) Alters(field, value any) (string, error) {
	syntheticValue := s.generator.GenerateSyntheticValue(coerceString(value), "startswith")
	return fmt.Sprintf("%v startswith '%v'", strings.ToLower(coerceString(field)), strings.ToLower(syntheticValue)), nil
}

type containsCS struct {
	generator *SyntheticDataGenerator
}

func (c containsCS) Alters(field, value any) (string, error) {
	syntheticValue := c.generator.GenerateSyntheticValue(coerceString(value), "contains")
	return fmt.Sprintf("%v contains '%v'", strings.ToLower(coerceString(field)), syntheticValue), nil
}

type endswithCS struct {
	generator *SyntheticDataGenerator
}

func (e endswithCS) Alters(field, value any) (string, error) {
	syntheticValue := e.generator.GenerateSyntheticValue(coerceString(value), "endswith")
	return fmt.Sprintf("%v endswith '%v'", strings.ToLower(coerceString(field)), syntheticValue), nil
}

type startswithCS struct {
	generator *SyntheticDataGenerator
}

func (s startswithCS) Alters(field, value any) (string, error) {
	syntheticValue := s.generator.GenerateSyntheticValue(coerceString(value), "startswith")
	return fmt.Sprintf("%v startswith '%v'", strings.ToLower(coerceString(field)), syntheticValue), nil
}

type re struct {
	generator *SyntheticDataGenerator
}

func (r re) Alters(field any, value any) (string, error) {
	syntheticValue := r.generator.GenerateSyntheticValue(coerceString(value), "re")
	return fmt.Sprintf("%v re '%v'", strings.ToLower(coerceString(field)), syntheticValue), nil
}

type cidr struct {
	generator *SyntheticDataGenerator
}

func (c cidr) Alters(field any, value any) (string, error) {
	syntheticValue := c.generator.GenerateSyntheticValue(coerceString(value), "cidr")
	return fmt.Sprintf("%v cidr '%v'", strings.ToLower(coerceString(field)), syntheticValue), nil
}

type gt struct {
	generator *SyntheticDataGenerator
}

func (g gt) Alters(field any, value any) (string, error) {
	syntheticValue := g.generator.GenerateSyntheticValue(coerceString(value), "gt")
	return fmt.Sprintf("%v gt '%v'", strings.ToLower(coerceString(field)), syntheticValue), nil
}

type gte struct {
	generator *SyntheticDataGenerator
}

func (g gte) Alters(field any, value any) (string, error) {
	syntheticValue := g.generator.GenerateSyntheticValue(coerceString(value), "gte")
	return fmt.Sprintf("%v gte '%v'", strings.ToLower(coerceString(field)), syntheticValue), nil
}

type lt struct {
	generator *SyntheticDataGenerator
}

func (l lt) Alters(field any, value any) (string, error) {
	syntheticValue := l.generator.GenerateSyntheticValue(coerceString(value), "lt")
	return fmt.Sprintf("%v lt '%v'", strings.ToLower(coerceString(field)), syntheticValue), nil
}

type lte struct {
	generator *SyntheticDataGenerator
}

func (l lte) Alters(field any, value any) (string, error) {
	syntheticValue := l.generator.GenerateSyntheticValue(coerceString(value), "lte")
	return fmt.Sprintf("%v lte '%v'", strings.ToLower(coerceString(field)), syntheticValue), nil
}

type b64 struct{}

func (b64) Modify(value any) (any, error) {
	return base64.StdEncoding.EncodeToString([]byte(coerceString(value))), nil
}

type wide struct{}

func (wide) Modify(value any) (any, error) {
	runes := utf16.Encode([]rune(coerceString(value)))
	bytes := make([]byte, 2*len(runes))
	for i, r := range runes {
		binary.LittleEndian.PutUint16(bytes[i*2:], r)
	}
	return coerceString(bytes), nil
}

func coerceString(v interface{}) string {
	switch vv := v.(type) {
	case string:
		return vv
	case []byte:
		return string(vv)
	default:
		return fmt.Sprint(vv)
	}
}

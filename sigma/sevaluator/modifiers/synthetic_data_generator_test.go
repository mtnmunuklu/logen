package modifiers_test

import (
	"net"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/mtnmunuklu/logen/sigma/sevaluator/modifiers"
)

func TestSyntheticDataGeneratorContains(t *testing.T) {
	generator := modifiers.NewSyntheticDataGenerator()
	expectedValue := "test_value"
	result := generator.GenerateSyntheticValue(expectedValue, "contains")

	if !strings.Contains(result, expectedValue) {
		t.Errorf("Expected result to contain %s, but got: %s", expectedValue, result)
	}
}

func TestSyntheticDataGeneratorStartsWith(t *testing.T) {
	generator := modifiers.NewSyntheticDataGenerator()
	expectedPrefix := "prefix_"
	result := generator.GenerateSyntheticValue(expectedPrefix, "startswith")

	if !strings.HasPrefix(result, expectedPrefix) {
		t.Errorf("Expected result to start with %s, but got: %s", expectedPrefix, result)
	}
}

func TestSyntheticDataGeneratorEndsWith(t *testing.T) {
	generator := modifiers.NewSyntheticDataGenerator()
	expectedSuffix := "_suffix"
	result := generator.GenerateSyntheticValue(expectedSuffix, "endswith")

	if !strings.HasSuffix(result, expectedSuffix) {
		t.Errorf("Expected result to end with %s, but got: %s", expectedSuffix, result)
	}
}

func TestSyntheticDataGeneratorRegex(t *testing.T) {
	generator := modifiers.NewSyntheticDataGenerator()
	expectedRegexPattern := "\\d{2}-BC\\S{4}"
	result := generator.GenerateSyntheticValue(expectedRegexPattern, "re")

	matched, err := regexp.MatchString(expectedRegexPattern, result)
	if err != nil || !matched {
		t.Errorf("Expected result to match regex pattern %s, but got: %s", expectedRegexPattern, result)
	}
}

func TestSyntheticDataGeneratorCIDR(t *testing.T) {
	generator := modifiers.NewSyntheticDataGenerator()
	expectedCIDR := "192.168.1.0/24"
	result := generator.GenerateSyntheticValue(expectedCIDR, "cidr")

	ip, _, err := net.ParseCIDR(expectedCIDR)
	if err != nil || result != ip.String() {
		t.Errorf("Expected result to match CIDR block %s, but got: %s", expectedCIDR, result)
	}
}

func TestSyntheticDataGeneratorGreaterThan(t *testing.T) {
	generator := modifiers.NewSyntheticDataGenerator()

	// Integer case
	expectedInt := 100
	resultInt, _ := strconv.Atoi(generator.GenerateSyntheticValue(strconv.Itoa(expectedInt), "gt"))

	if resultInt <= expectedInt {
		t.Errorf("Expected result to be greater than %d, but got: %d", expectedInt, resultInt)
	}

	// String case
	expectedStr := "test"
	resultStr := generator.GenerateSyntheticValue(expectedStr, "gt")

	if resultStr <= expectedStr {
		t.Errorf("Expected result to be greater than %s, but got: %s", expectedStr, resultStr)
	}
}

func TestSyntheticDataGeneratorGreaterOrEqual(t *testing.T) {
	generator := modifiers.NewSyntheticDataGenerator()

	// Integer case
	expectedInt := 100
	resultInt, _ := strconv.Atoi(generator.GenerateSyntheticValue(strconv.Itoa(expectedInt), "gte"))

	if resultInt < expectedInt {
		t.Errorf("Expected result to be greater than or equal to %d, but got: %d", expectedInt, resultInt)
	}

	// String case
	expectedStr := "test"
	resultStr := generator.GenerateSyntheticValue(expectedStr, "gte")

	if resultStr < expectedStr {
		t.Errorf("Expected result to be greater than or equal to %s, but got: %s", expectedStr, resultStr)
	}
}

func TestSyntheticDataGeneratorLessThan(t *testing.T) {
	generator := modifiers.NewSyntheticDataGenerator()

	// Integer case
	expectedInt := 100
	resultInt, _ := strconv.Atoi(generator.GenerateSyntheticValue(strconv.Itoa(expectedInt), "lt"))

	if resultInt >= expectedInt {
		t.Errorf("Expected result to be lesser than %d, but got: %d", expectedInt, resultInt)
	}

	// String case
	expectedStr := "test"
	resultStr := generator.GenerateSyntheticValue(expectedStr, "lt")

	if resultStr >= expectedStr {
		t.Errorf("Expected result to be lesser than %s, but got: %s", expectedStr, resultStr)
	}
}

func TestSyntheticDataGeneratorLesserOrEqual(t *testing.T) {
	generator := modifiers.NewSyntheticDataGenerator()

	// Integer case
	expectedInt := 100
	resultInt, _ := strconv.Atoi(generator.GenerateSyntheticValue(strconv.Itoa(expectedInt), "lte"))

	if resultInt > expectedInt {
		t.Errorf("Expected result to be lesser than or equal to %d, but got: %d", expectedInt, resultInt)
	}

	// String case
	expectedStr := "test"
	resultStr := generator.GenerateSyntheticValue(expectedStr, "lte")

	if resultStr > expectedStr {
		t.Errorf("Expected result to be lesser than or equal to %s, but got: %s", expectedStr, resultStr)
	}
}

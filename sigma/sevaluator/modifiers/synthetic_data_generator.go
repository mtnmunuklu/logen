package modifiers

import (
	"fmt"
	"math/rand"
	"net"
	"regexp/syntax"
	"strings"
	"time"
)

// SyntheticDataGenerator is used for generating synthetic data.
type SyntheticDataGenerator struct {
	randomGenerator *rand.Rand
}

// NewSyntheticDataGenerator creates a new instance of SyntheticDataGenerator.
func NewSyntheticDataGenerator() *SyntheticDataGenerator {
	source := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(source)
	return &SyntheticDataGenerator{
		randomGenerator: randomGenerator,
	}
}

// GenerateSyntheticValue generates a synthetic value based on a specific operation type.
func (g *SyntheticDataGenerator) GenerateSyntheticValue(value string, operationType string) string {
	syntheticData := g.generateRandomString(10)

	switch operationType {
	case "contains":
		randomIndex := g.randomGenerator.Intn(len(syntheticData) + 1)
		syntheticData = syntheticData[:randomIndex] + value + syntheticData[randomIndex:]
	case "startswith":
		syntheticData = value + syntheticData
	case "endswith":
		syntheticData = syntheticData + value
	case "re":
		syntheticData = g.generateRegexSyntheticData(value)
	case "cidr":
		syntheticData = g.generateCIDRMatch(value)
	default:
		syntheticData = value
	}

	return syntheticData
}

// GenerateRegexSyntheticData generates a synthetic value based on the given regex pattern.
func (g *SyntheticDataGenerator) generateRegexSyntheticData(pattern string) string {
	re, err := syntax.Parse(pattern, syntax.Perl)
	if err != nil {
		// Handle parsing error
		fmt.Println("Regex parsing error:", err)
		return ""
	}

	var syntheticDataBuilder strings.Builder
	g.generateRegexSyntheticDataRecursive(re, &syntheticDataBuilder)
	return syntheticDataBuilder.String()
}

func (g *SyntheticDataGenerator) generateRegexSyntheticDataRecursive(re *syntax.Regexp, builder *strings.Builder) {
	// If re.Sub is empty, perform other operations
	if re.Sub == nil || len(re.Sub) == 0 {
		// Handle other syntax.Op values
		g.handleOtherRegexOperations(re, builder)
		return
	}

	// Assuming there is only one sub-expression in OpRepeat
	sub := re.Sub[0]

	switch re.Op {
	case syntax.OpRepeat:
		// Repeat: build the repeated string and append to the builder
		var repeatedBuilder strings.Builder

		// OpRepeat sub-expression check
		if sub != nil {
			for j := 0; j < re.Min+g.randomGenerator.Intn(re.Max-re.Min+1); j++ {
				// Process sub-expressions within the repeated expression
				g.generateRegexSyntheticDataRecursive(sub, &repeatedBuilder)
			}
		}

		builder.WriteString(repeatedBuilder.String())
	default:
		// Handle other syntax.Op values
		for _, sub := range re.Sub {
			g.generateRegexSyntheticDataRecursive(sub, builder)
		}
	}
}

func (g *SyntheticDataGenerator) handleOtherRegexOperations(sub *syntax.Regexp, builder *strings.Builder) {
	switch sub.Op {
	case syntax.OpConcat:
		// Concatenation, process each sub-expression
		for _, concatSub := range sub.Sub {
			g.generateRegexSyntheticDataRecursive(concatSub, builder)
		}
	case syntax.OpCharClass:
		// Character class, randomly select a character
		class := sub.Rune
		if len(class) > 0 {
			char := class[g.randomGenerator.Intn(len(class))]
			builder.WriteRune(char)
		}
	case syntax.OpLiteral:
		// Append literal character to the builder
		builder.WriteString(sub.String())
	case syntax.OpCapture, syntax.OpNoMatch:
		// Recursively process the sub-expression for capture group or no match
		g.generateRegexSyntheticDataRecursive(sub, builder)
	default:
		// Unsupported syntax.Op value
		fmt.Println("Unsupported syntax.Op:", sub.Op)
	}
}

// GenerateCIDRMatch generates a synthetic value matching a specific CIDR block.
func (g *SyntheticDataGenerator) generateCIDRMatch(cidrValue string) string {
	_, ipNet, err := net.ParseCIDR(cidrValue)
	if err != nil {
		// Handle parsing error
		fmt.Println("CIDR parsing error:", err)
		return ""
	}

	// Calculate the number of available IP addresses in the CIDR block
	ipCount := 0
	for range g.IPRange(ipNet) {
		ipCount++
	}

	// Generate a random index within the range of available IPs
	randomIndex := g.randomGenerator.Intn(ipCount)

	// Get the IP address at the randomly selected index
	randomIP := g.GetNthIP(ipNet, randomIndex)

	return randomIP.String()
}

// IPRange generates all IP addresses in the given CIDR block
func (g *SyntheticDataGenerator) IPRange(ipNet *net.IPNet) []net.IP {
	ip := ipNet.IP.Mask(ipNet.Mask)
	var ips []net.IP
	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
		ips = append(ips, net.IP{ip[0], ip[1], ip[2], ip[3]})
	}
	return ips[1 : len(ips)-1] // Exclude network and broadcast addresses
}

// GetNthIP returns the nth IP address in the given CIDR block
func (g *SyntheticDataGenerator) GetNthIP(ipNet *net.IPNet, n int) net.IP {
	ip := ipNet.IP.Mask(ipNet.Mask)
	for i := 0; i < n; i++ {
		inc(ip)
	}
	return net.IP{ip[0], ip[1], ip[2], ip[3]}
}

// Helper function to increment IP address
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// GenerateRandomString generates a random string of a specific length.
func (g *SyntheticDataGenerator) generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[g.randomGenerator.Intn(len(charset))]
	}
	return string(result)
}

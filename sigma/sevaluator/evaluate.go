package sevaluator

import (
	"context"
	"fmt"
	"strings"

	"github.com/mtnmunuklu/logen/sigma"
)

// RuleEvaluator represents a rule evaluator that is capable of computing the search, condition, and query results of a Sigma rule.
// It holds the rule configuration, search conditions, and field mappings necessary to apply the rule to log events and generate the query results.
type RuleEvaluator struct {
	sigma.Rule
	config          []sigma.Config      // Additional configuration options to use when evaluating the rule
	indexes         []string            // The list of indexes that this rule should be applied to. Computed from the Logsource field in the rule and any config that's supplied.
	indexConditions []sigma.Search      // Any field-value conditions that need to match for this rule to apply to events from []indexes
	fieldmappings   map[string][]string // A compiled mapping from rule fieldnames to possible event fieldnames

	expandPlaceholder func(ctx context.Context, placeholderName string) ([]string, error) // A function to expand placeholders in the Sigma rule template
	caseSensitive     bool
}

// ForRule constructs a new RuleEvaluator with the given Sigma rule and evaluation options.
// It applies any provided options to the new RuleEvaluator and returns it.
func ForRule(rule sigma.Rule, options ...Option) *RuleEvaluator {
	e := &RuleEvaluator{Rule: rule}
	for _, option := range options {
		option(e)
	}
	return e
}

// Result represents the evaluation result of a Sigma rule.
// It contains the search, condition, aggregation, and query results of the rule evaluation.
type Result struct {
	Searches    map[string][]string // The map of search identifiers to their result values
	Conditions  map[int][]string    // The map of condition indices to their result values
	SourceTypes map[int]string      // The map of sourcetype indices to their result values
	Queries     map[int]string      // The map of query indices to their result values
}

// This function returns a Result object containing the evaluation results for the rule's Detection field.
// It uses the evaluateSearch, evaluateSearchExpression and evaluateAggregationExpression functions to compute the results.
func (rule RuleEvaluator) Alters(ctx context.Context) (Result, error) {
	result := Result{
		Searches:    make(map[string][]string),
		Conditions:  make(map[int][]string),
		SourceTypes: make(map[int]string),
		Queries:     make(map[int]string),
	}

	// Evaluate all the search expressions in the Detection field and store the results in the SearchResults map of the result object.
	for identifier, search := range rule.Detection.Searches {
		var err error
		result.Searches[identifier], err = rule.evaluateSearch(ctx, search)
		if err != nil {
			return Result{}, fmt.Errorf("error evaluating search %s: %w", identifier, err)
		}
	}

	// Evaluate all the search expressions in the Detection field's Conditions array and store the results in the ConditionResults map of the result object.
	// If a condition has an Aggregation field, also evaluate it and store the result in the AggregationResults map of the result object.
	for conditionIndex, condition := range rule.Detection.Conditions {
		result.Conditions[conditionIndex] = rule.evaluateSearchExpression(condition.Search, []string{}, true)
	}

	// Combine the search results and condition results to form the final query strings for each condition.
	// The query strings are stored in the QueryResults map of the result object.
	for i, conditionResult := range result.Conditions {
		conditionList := make([]string, 0, len(conditionResult))
		for _, condition := range conditionResult {
			// If the condition matches any search identifier, replace it with the corresponding search results
			if value, ok := result.Searches[condition]; ok {
				if len(conditionResult) > 1 {
					conditionList = append(conditionList, "("+strings.Join(value, " and ")+")")
				} else if len(value) > 1 {
					conditionList = append(conditionList, strings.Join(value, " and "))
				} else {
					conditionList = append(conditionList, strings.Join(value, ""))
				}
			} else {
				// If the condition doesn't match any search identifier, add it as is to the conditionList
				conditionList = append(conditionList, condition)
			}
		}

		// add the conditionList to the final query string
		result.Queries[i] = strings.Join(conditionList, "")

		// Add the sourcetype condition to the final query string, if applicable
		if rule.Logsource.Product != "" && rule.Logsource.Service != "" {
			result.SourceTypes[i] = rule.Logsource.Product + " " + rule.Logsource.Service
		} else if rule.Logsource.Product != "" && rule.Logsource.Service == "" {
			result.SourceTypes[i] = rule.Logsource.Product
		}
	}

	return result, nil
}

package context

import (
	"fmt"
	"strconv"
	"strings"
)

// DefaultPolicy is used when the configuration does not select a policy.
const DefaultPolicy = "k-callsite-2"

// NewSelector builds a context selector from a policy string. Supported
// policies are "ci" (context insensitive) and "k-callsite-<n>"; an empty
// policy defaults to DefaultPolicy. Unknown policies are reported instead of
// leaving a nil selector behind.
func NewSelector(policy string) (ContextSelector, error) {
	policy = strings.TrimSpace(policy)
	if policy == "" {
		policy = DefaultPolicy
	}
	if policy == "ci" {
		return CreateCISelector(), nil
	}
	if strings.HasPrefix(policy, "k-callsite-") {
		limit, err := strconv.Atoi(strings.TrimPrefix(policy, "k-callsite-"))
		if err != nil || limit < 0 {
			return nil, fmt.Errorf("invalid context policy %q: expected k-callsite-<non-negative integer>", policy)
		}
		return CreateKCallSiteSelector(limit), nil
	}
	return nil, fmt.Errorf("unsupported context policy %q (supported: ci, k-callsite-<n>)", policy)
}

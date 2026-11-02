package env0

import (
	"errors"
	"fmt"
	"net/netip"
	"strconv"
	"strings"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func validateIpAllowlistEntry(i any, p cty.Path) diag.Diagnostics {
	entry := i.(string)
	if err := checkIpAllowlistEntry(entry); err != nil {
		return diag.Errorf("invalid ip_allowlist entry '%s': %v. Use an address (198.51.100.7), a CIDR subnet (198.51.100.0/28) or a range (198.51.100.40-198.51.100.60)", entry, err)
	}

	return nil
}

// netip matches the server's checks: it rejects IPv4 octets with leading zeros and any whitespace,
// so an entry that passes is stored and read back unchanged.
func checkIpAllowlistEntry(entry string) error {
	if strings.Contains(entry, "/") {
		_, err := netip.ParsePrefix(entry)

		return err
	}

	if startText, endText, isRange := strings.Cut(entry, "-"); isRange {
		start, err := parseIpAllowlistAddr(startText)
		if err != nil {
			return err
		}

		end, err := parseIpAllowlistAddr(endText)
		if err != nil {
			return err
		}

		if start.BitLen() != end.BitLen() {
			return errors.New("both ends of a range must be in the same IP family")
		}

		if start.Compare(end) > 0 {
			return errors.New("the range start is after its end")
		}

		return nil
	}

	_, err := parseIpAllowlistAddr(entry)

	return err
}

func parseIpAllowlistAddr(text string) (netip.Addr, error) {
	addr, err := netip.ParseAddr(text)
	if err == nil && addr.Zone() != "" {
		return addr, errors.New("IPv6 zones are not supported")
	}

	return addr, err
}

// Same conversion as the env0 UI. Entries are validated already.
func toIpAllowlistRules(entries []any) []client.IpAllowlistRule {
	rules := make([]client.IpAllowlistRule, 0, len(entries))

	for _, e := range entries {
		entry := e.(string)

		family := "ipv4"
		if strings.Contains(entry, ":") {
			family = "ipv6"
		}

		if network, prefixText, isSubnet := strings.Cut(entry, "/"); isSubnet {
			prefix, _ := strconv.Atoi(prefixText)
			rules = append(rules, client.IpAllowlistRule{Type: "subnet", Network: network, Prefix: &prefix, Family: family})
		} else if start, end, isRange := strings.Cut(entry, "-"); isRange {
			rules = append(rules, client.IpAllowlistRule{Type: "range", Start: start, End: end, Family: family})
		} else {
			rules = append(rules, client.IpAllowlistRule{Type: "address", Address: entry, Family: family})
		}
	}

	return rules
}

func fromIpAllowlistRules(rules []client.IpAllowlistRule) []string {
	entries := make([]string, 0, len(rules))

	for _, rule := range rules {
		switch rule.Type {
		case "subnet":
			prefix := 0
			if rule.Prefix != nil {
				prefix = *rule.Prefix
			}

			entries = append(entries, fmt.Sprintf("%s/%d", rule.Network, prefix))
		case "range":
			entries = append(entries, rule.Start+"-"+rule.End)
		default:
			entries = append(entries, rule.Address)
		}
	}

	return entries
}

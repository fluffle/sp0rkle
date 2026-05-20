//go:build integration

package regtest

import (
	"regexp"
	"testing"
)

var reIPv4Netmask = regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`)

func testNetmask(t *testing.T) {
	t.Run("cidr_ipv4", func(t *testing.T) {
		_, err := h.CommandAndExpect("netmask 192.168.1.0/24", h.Regex(reIPv4Netmask))
		if err != nil {
			t.Errorf("expected IPv4 in response, got: %v", err)
		}
	})

	t.Run("cidr_ipv4_16", func(t *testing.T) {
		_, err := h.CommandAndExpect("netmask 10.0.0.0/16", h.Regex(reIPv4Netmask))
		if err != nil {
			t.Errorf("expected IPv4 in response, got: %v", err)
		}
	})

	t.Run("cidr_ipv4_8", func(t *testing.T) {
		_, err := h.CommandAndExpect("netmask 172.16.0.0/12", h.Regex(reIPv4Netmask))
		if err != nil {
			t.Errorf("expected 172.16 in response, got: %v", err)
		}
	})

	t.Run("cidr_ipv4_32", func(t *testing.T) {
		_, err := h.CommandAndExpect("netmask 192.168.1.1/32", h.Regex(reIPv4Netmask))
		if err != nil {
			t.Errorf("expected IPv4 in response, got: %v", err)
		}
	})

	t.Run("cidr_ipv6", func(t *testing.T) {
		_, err := h.CommandAndExpect("netmask 2001:db8::/32", h.Contains("2001:db8"))
		if err != nil {
			t.Errorf("expected 2001:db8 in response, got: %v", err)
		}
	})

	t.Run("cidr_ipv6_128", func(t *testing.T) {
		_, err := h.CommandAndExpect("netmask ::1/128", h.Contains("::1"))
		if err != nil {
			t.Errorf("expected ::1 in response, got: %v", err)
		}
	})

	t.Run("ip_and_mask_ipv4", func(t *testing.T) {
		_, err := h.CommandAndExpect("netmask 192.168.1.100 255.255.255.0", h.Regex(reIPv4Netmask))
		if err != nil {
			t.Errorf("expected IP in response, got: %v", err)
		}
	})

	t.Run("ip_and_mask_ipv6", func(t *testing.T) {
		_, err := h.CommandAndExpect("netmask 2001:db8::1 ffff:ffff:ffff:ffff::", h.Contains("2001:db8"))
		if err != nil {
			t.Errorf("expected 2001:db8 in response, got: %v", err)
		}
	})

	t.Run("ip_and_mask_mixed_error", func(t *testing.T) {
		_, err := h.CommandAndExpect("netmask 192.168.1.1 ffff:ffff::", h.Contains("can't mix"))
		if err != nil {
			t.Errorf("expected mixed version error, got: %v", err)
		}
	})

	t.Run("invalid_ip", func(t *testing.T) {
		_, err := h.CommandAndExpect("netmask notanip 255.255.255.0", h.Contains("couldn't be parsed"))
		if err != nil {
			t.Errorf("expected parse error for invalid IP, got: %v", err)
		}
	})

	t.Run("bad_args", func(t *testing.T) {
		_, err := h.CommandAndExpect("netmask", h.Contains("bad netmask args"))
		if err != nil {
			t.Errorf("expected bad args response, got: %v", err)
		}
	})

	t.Run("bad_args_single_no_cidr", func(t *testing.T) {
		_, err := h.CommandAndExpect("netmask justonearg", h.Contains("bad netmask args"))
		if err != nil {
			t.Errorf("expected bad args for single non-cidr arg, got: %v", err)
		}
	})

	t.Run("invalid_cidr", func(t *testing.T) {
		_, err := h.CommandAndExpect("netmask 999.999.999.999/32", h.Regex(regexp.MustCompile(`(?i)error|couldn't`)))
		if err != nil {
			t.Errorf("expected error for invalid CIDR, got: %v", err)
		}
	})

	t.Run("cidr_range_display", func(t *testing.T) {
		_, err := h.CommandAndExpect("netmask 192.168.1.0/24", h.Regex(regexp.MustCompile(`range \d+\.\d+\.\d+\.\d+-\d+\.\d+\.\d+\.\d+`)))
		if err != nil {
			t.Errorf("expected range in response, got: %v", err)
		}
	})

	t.Run("netmask_ipv4_invalid", func(t *testing.T) {
		_, err := h.CommandAndExpect("netmask 192.168.1.1 999.0.0.1", h.Contains("couldn't be parsed"))
		if err != nil {
			t.Errorf("expected parse error for invalid mask, got: %v", err)
		}
	})

	t.Run("case_insensitive", func(t *testing.T) {
		_, err := h.CommandAndExpect("NETMASK 10.0.0.0/8", h.Regex(reIPv4Netmask))
		if err != nil {
			t.Errorf("expected case insensitive command, got: %v", err)
		}
	})
}

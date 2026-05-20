//go:build integration

package regtest

import (
	"regexp"
	"testing"
)

func testCalc(t *testing.T) {
	t.Run("basic_addition", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc 2 + 2", h.Contains("4"))
		if err != nil {
			t.Errorf("expected '4' in response, got: %v", err)
		}
	})

	t.Run("basic_subtraction", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc 10 - 3", h.Contains("7"))
		if err != nil {
			t.Errorf("expected '7' in response, got: %v", err)
		}
	})

	t.Run("basic_multiplication", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc 6 * 7", h.Contains("42"))
		if err != nil {
			t.Errorf("expected '42' in response, got: %v", err)
		}
	})

	t.Run("basic_division", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc 100 / 4", h.Contains("25"))
		if err != nil {
			t.Errorf("expected '25' in response, got: %v", err)
		}
	})

	t.Run("modulo", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc 17 % 5", h.Contains("2"))
		if err != nil {
			t.Errorf("expected '2' in response, got: %v", err)
		}
	})

	t.Run("exponentiation_double_star", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc 2 ** 8", h.Contains("256"))
		if err != nil {
			t.Errorf("expected '256' in response, got: %v", err)
		}
	})

	t.Run("exponentiation_caret", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc 3 ^ 3", h.Contains("27"))
		if err != nil {
			t.Errorf("expected '27' in response, got: %v", err)
		}
	})

	t.Run("operator_precedence", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc 2 + 3 * 4", h.Contains("14"))
		if err != nil {
			t.Errorf("expected '14' (precedence), got: %v", err)
		}
	})

	t.Run("parentheses", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc (2 + 3) * 4", h.Contains("20"))
		if err != nil {
			t.Errorf("expected '20' (parentheses), got: %v", err)
		}
	})

	t.Run("math_functions", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc sqrt(144)", h.Contains("12"))
		if err != nil {
			t.Errorf("expected '12' in response, got: %v", err)
		}
	})

	t.Run("math_functions_floor", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc floor(3.7)", h.Contains("3"))
		if err != nil {
			t.Errorf("expected '3' in response, got: %v", err)
		}
	})

	t.Run("math_functions_ceil", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc ceil(3.2)", h.Contains("4"))
		if err != nil {
			t.Errorf("expected '4' in response, got: %v", err)
		}
	})

	t.Run("math_functions_abs", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc abs(-42)", h.Contains("42"))
		if err != nil {
			t.Errorf("expected '42' in response, got: %v", err)
		}
	})

	t.Run("math_functions_atan2", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc atan2(1, 1)", h.Contains("0.78"))
		if err != nil {
			t.Errorf("expected atan2(1,1) ~= 0.78, got: %v", err)
		}
	})

	t.Run("math_functions_hypot", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc hypot(3, 4)", h.Contains("5"))
		if err != nil {
			t.Errorf("expected '5' in response, got: %v", err)
		}
	})

	t.Run("math_functions_max_min", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc max(3, 7)", h.Contains("7"))
		if err != nil {
			t.Errorf("expected '7' in response, got: %v", err)
		}
	})

	t.Run("constants_pi", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc pi", h.Contains("3.14"))
		if err != nil {
			t.Errorf("expected pi ~= 3.14, got: %v", err)
		}
	})

	t.Run("constants_e", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc e", h.Contains("2.71"))
		if err != nil {
			t.Errorf("expected e ~= 2.71, got: %v", err)
		}
	})

	t.Run("constants_phi", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc phi", h.Contains("1.61"))
		if err != nil {
			t.Errorf("expected phi ~= 1.61, got: %v", err)
		}
	})

	t.Run("constants_answer", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc answer", h.Contains("42"))
		if err != nil {
			t.Errorf("expected '42' in response, got: %v", err)
		}
	})

	t.Run("result_variable", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc 2 + 3", h.Contains("5"))
		if err != nil {
			t.Errorf("expected '5' for initial response, got: %v", err)
		}
		_, err = h.CommandAndExpect("calc 2 + 3 + result", h.Contains("10"))
		if err != nil {
			t.Errorf("expected '10' for result response, got: %v", err)
		}
	})

	t.Run("error_empty_input", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc", h.Contains("error"))
		if err != nil {
			t.Errorf("expected error response for empty calc, got: %v", err)
		}
	})

	t.Run("error_invalid_expression", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc foo + bar", h.Contains("error"))
		if err != nil {
			t.Errorf("expected error response for invalid expression, got: %v", err)
		}
	})

	t.Run("error_syntax_error", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc 2 ++ 2", h.Contains("error"))
		if err != nil {
			t.Errorf("expected error response for syntax error, got: %v", err)
		}
	})

	t.Run("float_result", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc 10 / 3", h.Contains("3.33"))
		if err != nil {
			t.Errorf("expected float result for 10/3, got: %v", err)
		}
	})

	t.Run("complex_expression", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc (2 + 3) * (4 - 1)", h.Contains("15"))
		if err != nil {
			t.Errorf("expected '15' for complex expression, got: %v", err)
		}
	})

	t.Run("case_insensitive", func(t *testing.T) {
		_, err := h.CommandAndExpect("CALC 1 + 1", h.Contains("2"))
		if err != nil {
			t.Errorf("expected case insensitive command, got: %v", err)
		}
	})

	t.Run("log_functions", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc log2(1024)", h.Contains("10"))
		if err != nil {
			t.Errorf("expected '10' for log2(1024), got: %v", err)
		}
	})

	t.Run("log10_function", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc log10(1000)", h.Contains("3"))
		if err != nil {
			t.Errorf("expected '3' for log10(1000), got: %v", err)
		}
	})

	t.Run("math_functions_cbrt", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc cbrt(27)", h.Contains("3"))
		if err != nil {
			t.Errorf("expected '3' for cbrt(27), got: %v", err)
		}
	})

	t.Run("math_functions_exp", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc exp(0)", h.Contains("1"))
		if err != nil {
			t.Errorf("expected '1' for exp(0), got: %v", err)
		}
	})

	t.Run("math_functions_gamma", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc gamma(1)", h.Contains("1"))
		if err != nil {
			t.Errorf("expected '1' for gamma(1), got: %v", err)
		}
	})

	t.Run("math_functions_sin_cos", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc sin(0)", h.Contains("0"))
		if err != nil {
			t.Errorf("expected '0' for sin(0), got: %v", err)
		}
	})

	t.Run("math_functions_tan", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc tan(0)", h.Contains("0"))
		if err != nil {
			t.Errorf("expected '0' for tan(0), got: %v", err)
		}
	})

	t.Run("math_functions_int_trunc", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc int(3.9)", h.Contains("3"))
		if err != nil {
			t.Errorf("expected '3' for int(3.9), got: %v", err)
		}
	})

	t.Run("math_functions_logb", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc logb(8)", h.Contains("3"))
		if err != nil {
			t.Errorf("expected '3' for logb(8), got: %v", err)
		}
	})

	t.Run("prefix_minus", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc -5 + 3", h.Contains("-2"))
		if err != nil {
			t.Errorf("expected '-2' for -5+3, got: %v", err)
		}
	})

	t.Run("negative_constant", func(t *testing.T) {
		_, err := h.CommandAndExpect("calc -pi", h.Regex(regexp.MustCompile(`-3\.14`)))
		if err != nil {
			t.Errorf("expected -pi, got: %v", err)
		}
	})
}

package enmime

import "testing"

func TestErrorStringConversion(t *testing.T) {
	e := &Error{
		Name:   "WarnName",
		Detail: "Warn Details",
		Severe: false,
	}

	want := "[W] WarnName: Warn Details"
	got := e.String()
	if got != want {
		t.Error("got:", got, "want:", want)
	}

	e = &Error{
		Name:   "ErrorName",
		Detail: "Error Details",
		Severe: true,
	}

	want = "[E] ErrorName: Error Details"
	got = e.String()
	if got != want {
		t.Error("got:", got, "want:", want)
	}
}

func TestErrorAddError(t *testing.T) {
	p := &Part{}
	p.addError(errorMalformedHeader, "1 %v %q", 2, "three")

	if len(p.Errors) != 1 {
		t.Fatal("len(p.Errors) ==", len(p.Errors), ", want: 1")
	}
	e := p.Errors[0]

	if e.Name != string(errorMalformedHeader) {
		t.Errorf("e.Name == %q, want: %q", e.Name, errorMalformedHeader)
	}
	if !e.Severe {
		t.Errorf("e.Severe == %v, want: true", e.Severe)
	}
	want := "1 2 \"three\""
	if e.Detail != want {
		t.Errorf("e.Detail == %q, want: %q", e.Detail, want)
	}
}

func TestErrorAddWarning(t *testing.T) {
	p := &Part{}
	p.addWarning(errorMalformedHeader, "1 %v %q", 2, "three")

	if len(p.Errors) != 1 {
		t.Fatal("len(p.Errors) ==", len(p.Errors), ", want: 1")
	}
	e := p.Errors[0]

	if e.Name != string(errorMalformedHeader) {
		t.Errorf("e.Name == %q, want: %q", e.Name, errorMalformedHeader)
	}
	if e.Severe {
		t.Errorf("e.Severe == %v, want: false", e.Severe)
	}
	want := "1 2 \"three\""
	if e.Detail != want {
		t.Errorf("e.Detail == %q, want: %q", e.Detail, want)
	}
}

func TestErrorEnvelopeWarnings(t *testing.T) {
	// To pass each file below must error one or more times with the specified errorName, and no
	// other errorNames.
	var files = []struct {
		filename string
		perror   errorName
	}{
		{"bad-final-boundary.raw", errorMissingBoundary},
		{"bad-header-wrap.raw", errorMalformedHeader},
		{"html-only-inline.raw", errorPlainTextFromHTML},
		{"missing-content-type.raw", errorMissingContentType},
		{"missing-content-type2.raw", errorMissingContentType},
		{"empty-header.raw", errorMissingContentType},
		{"unk-encoding-part.raw", errorContentEncoding},
		{"unk-charset-html-only.raw", errorCharsetConversion},
		{"unk-charset-part.raw", errorCharsetConversion},
	}

	for _, tt := range files {
		// Mail with disposition attachment
		msg := openTestData("low-quality", tt.filename)
		e, err := ReadEnvelope(msg)
		if err != nil {
			t.Fatal("Failed to parse MIME:", err)
		}

		if len(e.Errors) == 0 {
			t.Error("Got 0 warnings, expected at least one for:", tt.filename)
		}

		satisfied := false
		for _, perr := range e.Errors {
			if perr.Name == string(tt.perror) {
				satisfied = true
				if perr.Severe {
					t.Errorf("Expected Severe to be false, got true for %q", perr.Name)
				}
			}
		}

		if !satisfied {
			var errorList string
			for _, perr := range e.Errors {
				errorList += perr.String()
				errorList += "\n"
			}
			t.Errorf(
				"File %q should have error of type %q, got these instead:\n%s",
				tt.filename,
				tt.perror,
				errorList)
		}
	}
}

package manifest

import "testing"

func TestSource_WithoutRepository(t *testing.T) {
	source := Source{
		Host: "source.com",
		Tag:  "v1.0.0",
	}

	target := Target{
		Host: "target.com",
	}

	targetWithRepository := Target{
		Host:       "target.com",
		Repository: "bar",
	}

	testCases := []struct {
		source         Source
		target         Target
		expectedSource string
		expectedTarget string
	}{
		{
			source,
			target,
			"source.com:v1.0.0",
			"target.com:v1.0.0",
		},
		{
			source,
			targetWithRepository,
			"source.com:v1.0.0",
			"target.com/bar:v1.0.0",
		},
	}

	for _, testCase := range testCases {
		testCase.source.Target = testCase.target

		if testCase.source.Image() != testCase.expectedSource {
			t.Errorf("expected source %s, actual %s", testCase.expectedSource, testCase.source.Image())
		}

		if testCase.source.TargetImage() != testCase.expectedTarget {
			t.Errorf("expected target %s, actual %s", testCase.expectedTarget, testCase.source.TargetImage())
		}
	}
}

func TestSource_WithRepository(t *testing.T) {
	source := Source{
		Host:       "source.com",
		Repository: "repo",
		Tag:        "v1.0.0",
	}

	target := Target{
		Host: "target.com",
	}

	targetWithRepository := Target{
		Host:       "target.com",
		Repository: "bar",
	}

	testCases := []struct {
		source         Source
		target         Target
		expectedSource string
		expectedTarget string
	}{
		{
			source,
			target,
			"source.com/repo:v1.0.0",
			"target.com/repo:v1.0.0",
		},
		{
			source,
			targetWithRepository,
			"source.com/repo:v1.0.0",
			"target.com/bar/repo:v1.0.0",
		},
	}

	for _, testCase := range testCases {
		testCase.source.Target = testCase.target

		if testCase.source.Image() != testCase.expectedSource {
			t.Errorf("expected source %s, actual %s", testCase.expectedSource, testCase.source.Image())
		}

		if testCase.source.TargetImage() != testCase.expectedTarget {
			t.Errorf("expected target %s, actual %s", testCase.expectedTarget, testCase.source.TargetImage())
		}
	}
}

func TestSource_WithNestedRepository(t *testing.T) {
	source := Source{
		Host:       "source.com",
		Repository: "repo/foo",
		Tag:        "v1.0.0",
	}

	target := Target{
		Host: "target.com",
	}

	targetWithRepository := Target{
		Host:       "target.com",
		Repository: "bar",
	}

	testCases := []struct {
		source         Source
		target         Target
		expectedSource string
		expectedTarget string
	}{
		{
			source,
			target,
			"source.com/repo/foo:v1.0.0",
			"target.com/repo/foo:v1.0.0",
		},
		{
			source,
			targetWithRepository,
			"source.com/repo/foo:v1.0.0",
			"target.com/bar/repo/foo:v1.0.0",
		},
	}

	for _, testCase := range testCases {
		testCase.source.Target = testCase.target

		if testCase.source.Image() != testCase.expectedSource {
			t.Errorf("expected source %s, actual %s", testCase.expectedSource, testCase.source.Image())
		}

		if testCase.source.TargetImage() != testCase.expectedTarget {
			t.Errorf("expected target %s, actual %s", testCase.expectedTarget, testCase.source.TargetImage())
		}
	}
}

func TestSource_Digest(t *testing.T) {
	source := Source{
		Host:       "source.com",
		Repository: "repo",
		Digest:     "sha256:123",
	}

	target := Target{
		Host: "target.com",
	}

	source.Target = target

	const expectedSource = "source.com/repo@sha256:123"
	if source.Image() != expectedSource {
		t.Errorf("unexpected source string. expected %s, actual %s", source.Image(), expectedSource)
	}

	const expectedTarget = "target.com/repo:123"
	if source.TargetImage() != expectedTarget {
		t.Errorf("unexpected target string. expected %s, actual %s", source.TargetImage(), expectedTarget)
	}
}

package fixtures

import (
	"path/filepath"
	"testing"

	"github.com/bblfsh/csharp-driver/driver/normalizer"
	"gopkg.in/bblfsh/sdk.v2/driver"
	"gopkg.in/bblfsh/sdk.v2/driver/fixtures"
	"gopkg.in/bblfsh/sdk.v2/driver/native"
)

const projectRoot = "../../"

var Suite = &fixtures.Suite{
	Lang: "csharp",
	Ext:  ".cs",
	Path: filepath.Join(projectRoot, fixtures.Dir),
	NewDriver: func() driver.Native {
		return native.NewDriverAt(filepath.Join(projectRoot, "build/bin/native"), native.UTF8)
	},
	Transforms: normalizer.Transforms,
	BenchName:  "parser_context",
	Semantic: fixtures.SemanticConfig{
		BlacklistTypes: []string{
			// TODO: more types and uncomment
			"IdentifierName",
			"Block",
			"UsingDirective",
			"QualifiedName",
			"StringLiteralExpression",
			"SingleLineCommentTrivia",
			"SingleLineDocumentationCommentTrivia",
			//"MultiLineCommentTrivia",
			"TrueLiteralExpression",
			"FalseLiteralExpression",
			"Parameter",
		},
	},
	Docker: fixtures.DockerConfig{
		//Image:"image:tag", // TODO: specify a docker image with language runtime
	},
}

func TestCsharpDriver(t *testing.T) {
	Suite.RunTests(t)
}

func BenchmarkCsharpDriver(b *testing.B) {
	Suite.RunBenchmarks(b)
}

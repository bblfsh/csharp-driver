package fixtures

import (
	"path/filepath"
	"testing"

	"github.com/bblfsh/csharp-driver/driver/normalizer"
	"gopkg.in/bblfsh/sdk.v2/driver"
	"gopkg.in/bblfsh/sdk.v2/driver/fixtures"
	"gopkg.in/bblfsh/sdk.v2/driver/native"
	"gopkg.in/bblfsh/sdk.v2/uast/transformer/positioner"
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
			"ArgListKeyword",
			"Block",
			"ConstructorDeclaration",
			"DestructorDeclaration",
			"FalseLiteralExpression",
			//"IdentifierName", // FIXME
			"IdentifierToken",
			"MethodDeclaration",
			"MultiLineCommentTrivia",
			"Parameter",
			"QualifiedName",
			"SingleLineCommentTrivia",
			"SingleLineDocumentationCommentTrivia",
			"StringLiteralExpression",
			"TrueLiteralExpression",
			"UsingDirective",
		},
	},
	VerifyTokens: []positioner.VerifyToken{
		{Types: []string{
			"IdentifierToken",
			"ClassKeyword",
			"FalseLiteralExpression",
			"MultiLineCommentTrivia",
			"QualifiedName",
			"SingleLineCommentTrivia",
			"SingleLineDocumentationCommentTrivia",
			"StringLiteralExpression",
			"TrueLiteralExpression",
		}},
	},
}

func TestCsharpDriver(t *testing.T) {
	Suite.RunTests(t)
}

func BenchmarkCsharpDriver(b *testing.B) {
	Suite.RunBenchmarks(b)
}

package normalizer

import (
	"gopkg.in/bblfsh/sdk.v2/uast"
	. "gopkg.in/bblfsh/sdk.v2/uast/transformer"
)

var Preprocess = Transformers([][]Transformer{
	{Mappings(Preprocessors...)},
}...)

var Normalize = Transformers([][]Transformer{

	{Mappings(Normalizers...)},
}...)

// Preprocessors is a block of AST preprocessing rules rules.
var Preprocessors = []Mapping{
	// replace "IsEmpty" with "Length: 0" in spans
	Map(
		Part("_", Obj{
			uast.KeyType: String("TextSpan"),
			"IsEmpty":    AnyNode(nil),
		}),
		Part("_", Obj{
			uast.KeyType: String("TextSpan"),
			"Length":     Int(0),
		}),
	),

	// restore "Start" if it was removed as a default value
	Map(
		Check(
			Not(Has{
				"Start": AnyNode(nil),
			}),
			Obj{
				uast.KeyType: String("TextSpan"),
				"Length":     Var("length"),
				"End":        Var("end"),
			},
		),
		Obj{
			uast.KeyType: String("TextSpan"),
			"Start":      Int(0),
			"Length":     Var("length"),
			"End":        Var("end"),
		},
	),

	// remove SpanStart
	Map(
		// TODO(dennwc): add it as a custom position field?
		Part("_", Obj{
			"SpanStart": AnyNode(nil),
		}),
		Part("_", Obj{}),
	),

	Map(
		Part("_", Obj{
			"FullSpan": Obj{
				uast.KeyType: String("TextSpan"),
				"Length":     AnyNode(nil),
				"Start":      Var("start"),
				"End":        Var("end"),
			},
			// TODO(dennwc): add it as a custom position field?
			"Span": AnyNode(nil),
		}),
		Part("_", Obj{
			// remap to temporary keys and let ObjectToNode to pick them up
			"spanStart": Var("start"),
			"spanEnd":   Var("end"),
		}),
	),

	ObjectToNode{
		OffsetKey:    "spanStart",
		EndOffsetKey: "spanEnd",
	}.Mapping(),
}

// Normalizers is the main block of normalization rules to convert native AST to semantic UAST.
var Normalizers = []Mapping{
	MapSemantic("IdentifierNameSyntax", uast.Identifier{}, MapObj(
		Obj{
			"Identifier": Obj{
				uast.KeyType: String("SyntaxToken"),
				// TODO(dennwc): assert that it's the same as in parent
				uast.KeyPos: AnyNode(nil),

				"LeadingTrivia":  Arr(),
				"TrailingTrivia": Arr(),
				"RawKind":        Int(8508),
				"Text":           Var("name"),
				"Value":          Var("name"),
				"ValueText":      Var("name"),
			},
			"RawKind": Int(8616),
		},
		Obj{
			"Name": Var("name"),
		},
	)),

	MapSemantic("LiteralExpressionSyntax", uast.String{}, MapObj(
		Obj{
			"Token": Obj{
				uast.KeyType: String("SyntaxToken"),
				// TODO(dennwc): assert that it's the same as in parent
				uast.KeyPos: AnyNode(nil),

				"LeadingTrivia":  Arr(),
				"TrailingTrivia": Arr(),
				"RawKind":        Int(8511),
				// contains escaped value, we don't need it in this mode
				"Text":      AnyNode(nil),
				"Value":     Var("val"),
				"ValueText": Var("val"),
			},
			"RawKind": Int(8750),
		},
		Obj{
			"Value": Var("val"),
		},
	)),

	MapSemantic("BlockSyntax", uast.Block{}, MapObj(
		Obj{
			"Statements": Var("stmts"),
			"RawKind":    Int(8792),
			// TODO(dennwc): remap to custom positional fields
			"OpenBraceToken":  AnyNode(nil),
			"CloseBraceToken": AnyNode(nil),
		},
		Obj{
			"Statements": Var("stmts"),
		},
	)),
}

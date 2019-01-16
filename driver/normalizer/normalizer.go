package normalizer

import (
	"gopkg.in/bblfsh/sdk.v2/uast"
	"gopkg.in/bblfsh/sdk.v2/uast/nodes"
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
	// Erase Whitespace and EndOfLine trivias.
	Map(
		Obj{
			uast.KeyType: Check(
				In(nodes.String("WhitespaceTrivia"), nodes.String("EndOfLineTrivia")),
				AnyNode(nil),
			),
			"FullSpan":    AnyNode(nil),
			"Span":        AnyNode(nil),
			"SpanStart":   AnyNode(nil),
			"IsDirective": Bool(false),
		},
		// cannot delete directly, so set to nil
		Is(nil),
	),

	// Now all whitespace "SyntaxTrivia" nodes are nil, we need
	// to cleanup arrays that were hosting those nodes.
	//
	// Find "LeadingTrivia" and "TrailingTrivia" and replace arrays
	// where all nodes are nil with an empty array.
	// TODO(dennwc): this only works if all nodes are nil, but real nodes
	//               may contain whitespace and comment Trivias in the
	//               same field
	Map(
		Part("_", Obj{
			"LeadingTrivia": Check(
				And(All(Is(nil)), OfKind(nodes.KindArray)),
				AnyNode(nil),
			),
		}),
		Part("_", Obj{
			"LeadingTrivia": Arr(),
		}),
	),
	Map(
		Part("_", Obj{
			"TrailingTrivia": Check(
				And(All(Is(nil)), OfKind(nodes.KindArray)),
				AnyNode(nil),
			),
		}),
		Part("_", Obj{
			"TrailingTrivia": Arr(),
		}),
	),

	// Drop "IsEmpty" field from TextSpan.
	// We can detect it with "Length == 0", if necessary.
	Map(
		Part("_", Obj{
			uast.KeyType: String("TextSpan"),
			"IsEmpty":    AnyNode(nil),
		}),
		Part("_", Obj{
			uast.KeyType: String("TextSpan"),
		}),
	),

	// Remove SpanStart from nodes. It duplicates positional info.
	// TODO(dennwc): add it as a custom position field?
	Map(
		Part("_", Obj{
			"SpanStart": AnyNode(nil),
		}),
		Part("_", Obj{}),
	),

	// Positional info is stored in a child node in FullSpan field.
	//
	// This is not supported by ObjectToNode helper, and we are
	// too lazy to create positional node by hand.
	//
	// Instead, we will temporary remap positional info to
	// "spanStart" and "spanEnd" fields of the root node, and
	// ObjectToNode will pick them up later to build a proper
	// positional node.
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

	// Use temporary fields from the previous transform to create positional node.
	ObjectToNode{
		OffsetKey:    "spanStart",
		EndOffsetKey: "spanEnd",
	}.Mapping(),
}

// Normalizers is the main block of normalization rules to convert native AST to semantic UAST.
var Normalizers = []Mapping{

	MapSemantic("IdentifierName", uast.Identifier{}, MapObj(
		Obj{
			"Identifier": Obj{
				uast.KeyType: String("IdentifierToken"),
				// TODO(dennwc): assert that it's the same as in parent
				uast.KeyPos: AnyNode(nil),

				// trivia == whitespace; can safely drop it
				"LeadingTrivia":  AnyNode(nil),
				"TrailingTrivia": AnyNode(nil),

				"IsMissing": Bool(false),

				// all token values are the same
				"Text":      Var("name"),
				"Value":     Var("name"),
				"ValueText": Var("name"),
			},
			"Arity": Int(0),

			// TODO(dennwc): these assertions might not be valid for all cases
			//               and will break this annotation, but at least it will
			//               help us detect the case when it's not valid
			"IsMissing":          Bool(false),
			"IsStructuredTrivia": Bool(false),
			"IsUnmanaged":        Bool(false),

			// TODO(dennwc): might be useful later; drop it for now
			"IsVar": AnyNode(nil),
		},
		Obj{
			"Name": Var("name"),
		},
	)),

	MapSemantic("StringLiteralExpression", uast.String{}, MapObj(
		Obj{
			"Token": Obj{
				uast.KeyType: String("StringLiteralToken"),
				// TODO(dennwc): assert that it's the same as in parent
				uast.KeyPos: AnyNode(nil),

				// trivia == whitespace; can safely drop it
				"LeadingTrivia":  AnyNode(nil),
				"TrailingTrivia": AnyNode(nil),

				"IsMissing": Bool(false),

				// contains escaped value, we don't need it in canonical UAST
				"Text": AnyNode(nil),

				// both values are the same
				"Value":     Var("val"),
				"ValueText": Var("val"),
			},
			"IsMissing":          Bool(false),
			"IsStructuredTrivia": Bool(false),
		},
		Obj{
			"Value": Var("val"),
		},
	)),

	MapSemantic("Block", uast.Block{}, MapObj(
		Obj{
			"Statements": Var("stmts"),
			// TODO(dennwc): remap to custom positional fields
			"OpenBraceToken":  AnyNode(nil),
			"CloseBraceToken": AnyNode(nil),
		},
		Obj{
			"Statements": Var("stmts"),
		},
	)),

	// Import (aka UsingDirectiveSyntax) is trivial.
	//
	// "Name" field is QualifiedIdentifier or Identifier and we remap to
	// "Path" in Import.
	//
	// Also, C# assumes that "using" statement imports all the symbols
	// from that package, so we also set an "All" field on Import.
	MapSemantic("UsingDirective", uast.Import{}, MapObj(
		Obj{
			"Name": Var("path"),
			// TODO(dennwc): remap to custom positional fields
			"SemicolonToken": AnyNode(nil),
			"UsingKeyword":   AnyNode(nil),
		},
		Obj{
			"Path": Var("path"),
			"All":  Bool(true),
		},
	)),

	// QualifiedIdentifier case is interesting in the sense that AST nodes
	// are organized as a linked list.
	//
	// The root QualifiedNameSyntax node will have the "Right" field pointing
	// to an Identifier (it was already converted from IdentifierNameSyntax)
	// and the "Left" field may either point to another Identifier or
	// to the next QualifiedNameSyntax (down to root of the package hierarchy).
	//
	// For the first case, we create a single QualifiedIdentifier node
	// by making a "Names" array from "Left" and "Right" Identifiers.
	//
	// For the second case, we rely on the fact that transforms are
	// using the DFS order. We assert that "Left" is already a
	// QualifiedIdentifier (all children were converted by DFS) and
	// save its "Names". Then we can simply create a new QualifiedIdentifier
	// and append "Right" (Identifier) to the end of the saved "Names" array.
	MapSemantic("QualifiedName", uast.QualifiedIdentifier{}, MapObj(
		CasesObj("case",
			// common
			Obj{
				"Right": Var("right"),
			},
			Objs{
				// the last name = identifier
				{
					"Left": Check(HasType(uast.Identifier{}), Var("left")),
				},
				// linked list
				{
					"Left": UASTType(uast.QualifiedIdentifier{}, Obj{
						// FIXME: start position
						uast.KeyPos: AnyNode(nil),
						"Names":     Var("names"),
					}),
				},
			},
		),
		CasesObj("case", nil,
			Objs{
				// the last name = identifier
				{
					"Names": Arr(Var("left"), Var("right")),
				},
				// linked list
				{
					"Names": Append(Var("names"), Arr(Var("right"))),
				},
			},
		),
	)),
}

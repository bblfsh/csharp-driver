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

func funcDefMap(typ string) Mapping {
	return MapSemantic(typ, uast.FunctionGroup{}, MapObj(
		Obj{
			"Body": Var("body"),
			"Identifier": Obj{
				uast.KeyType: String("IdentifierToken"),
				// TODO: assert that is the same as in parent
				uast.KeyPos:      Var("id_pos"),
				"LeadingTrivia":  Any(),
				"TrailingTrivia": Any(),
				"IsMissing": Bool(false),
				"Text":      Var("name"),
				"Value":     Var("name"),
				"ValueText": Var("name"),
			},
			"ParameterList": Obj{
				uast.KeyType: String("ParameterList"),
				uast.KeyPos: Any(),
				"OpenParenToken": Any(),
				"CloseParenToken": Any(),
				"IsMissing": Bool(false),
				"IsStructuredTrivia": Bool(false),
				"Parameters": Var("params"),
			},
			"ReturnType": Cases("cases_return",
				// Type = PredefinedType
				Obj{
					uast.KeyType: String("PredefinedType"),
					uast.KeyPos: Any(),
					"IsMissing": Bool(false),
					"IsStructuredTrivia": Bool(false),
					"IsUnmanaged": Any(),
					"IsVar": Any(),
					"Keyword": Obj{
						uast.KeyType: Any(),
						uast.KeyPos: Any(),
						"IsMissing": Bool(false),
						"LeadingTrivia": Any(),
						"TrailingTrivia": Any(),
						"Text": Var("rettype_text"),
						"Value": Var("rettype_text"),
						"ValueText": Var("rettype_text"),
					},
				},
				// Type = Identifier (User type)
				Var("rettype_ident"),
			),
		},

		Obj{
			"Nodes": Arr(
				// TODO: add docs and annotations
				UASTType(uast.Alias{}, Obj{
					"Name": UASTType(uast.Identifier{}, Obj{
						"Name": Var("name"),
					}),
					"Node": UASTType(uast.Function{}, Obj{
						"Body": Var("body"),
						"Type": UASTType(uast.FunctionType{}, Obj{
							"Arguments": Var("params"),
							"Returns": Arr(
								Cases("cases_return",
									// Type = PredefinedType
									UASTType(uast.Argument{}, Obj{
										"Type": UASTType(uast.Identifier{}, Obj{
											"Name": Var("rettype_text"),
										}),
									}),
									// Type = Identifier
									UASTType(uast.Argument{}, Obj{
										"Type": Var("rettype_ident"),
									}),
								),
							),
						}),
					}),
				}),
			),
		},
	))
}

// Preprocessors is a block of AST preprocessing rules rules.
var Preprocessors = []Mapping{
	// Erase Whitespace and EndOfLine trivias.
	Map(
		Obj{
			uast.KeyType: Check(
				In(nodes.String("WhitespaceTrivia"), nodes.String("EndOfLineTrivia")),
				Any(),
			),
			"FullSpan":    Any(),
			"Span":        Any(),
			"SpanStart":   Any(),
			"IsDirective": Bool(false),
		},
		// cannot delete directly, so set to nil
		Is(nil),
	),

	// Now all whitespace "SyntaxTrivia" nodes are nil, we need
	// to cleanup arrays that were hosting those nodes.
	//
	// Find "LeadingTrivia" and "TrailingTrivia" and drop nil
	// elements in those arrays.
	Map(
		Part("_", Obj{
			"LeadingTrivia": dropNils{Var("arr")},
		}),
		Part("_", Obj{
			"LeadingTrivia": Var("arr"),
		}),
	),
	Map(
		Part("_", Obj{
			"TrailingTrivia": dropNils{Var("arr")},
		}),
		Part("_", Obj{
			"TrailingTrivia": Var("arr"),
		}),
	),

	// Drop "IsEmpty" field from TextSpan.
	// We can detect it with "Length == 0", if necessary.
	Map(
		Part("_", Obj{
			uast.KeyType: String("TextSpan"),
			"IsEmpty":    Any(),
		}),
		Part("_", Obj{
			uast.KeyType: String("TextSpan"),
		}),
	),

	// Remove SpanStart from nodes. It duplicates positional info.
	// TODO(dennwc): add it as a custom position field?
	Map(
		Part("_", Obj{
			"SpanStart": Any(),
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
				"Length":     Any(),
				"Start":      Var("start"),
				"End":        Var("end"),
			},
			// TODO(dennwc): add it as a custom position field?
			"Span": Any(),
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

	// Add an empty @token field to comment nodes. It's necessary to pass the check
	// in the comment extractor.
	Map(
		Part("_", Obj{
			uast.KeyType: String("SingleLineCommentTrivia"),
		}),
		Part("_", Obj{
			uast.KeyType:  String("SingleLineCommentTrivia"),
			uast.KeyToken: String(""),
		}),
	),
	Map(
		Part("_", Obj{
			uast.KeyType: String("SingleLineDocumentationCommentTrivia"),
		}),
		Part("_", Obj{
			uast.KeyType:  String("SingleLineDocumentationCommentTrivia"),
			uast.KeyToken: String(""),
		}),
	),
}

// Normalizers is the main block of normalization rules to convert native AST to semantic UAST.
var Normalizers = []Mapping{

	// remove empty identifiers
	Map(
		Check(
			Has{
				uast.KeyType: String("IdentifierName"),
				"IsMissing":  Bool(true), // is set for empty identifiers
				"Identifier": Has{
					// make sure it's empty, we don't want to wipe something useful
					"Text":      String(""),
					"Value":     String(""),
					"ValueText": String(""),
				},
			},
			Any(),
		),
		Is(nil),
	),

	MapSemantic("IdentifierName", uast.Identifier{}, MapObj(
		Obj{
			"Identifier": Obj{
				uast.KeyType: String("IdentifierToken"),
				// TODO(dennwc): assert that it's the same as in parent
				uast.KeyPos: Any(),

				// trivia == whitespace; can safely drop it
				"LeadingTrivia":  Any(),
				"TrailingTrivia": Any(),

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

			// TODO(dennwc): this is true for Value == "unmanaged" and it looks
			//				 more like a keyword, probably unrecognized one
			"IsUnmanaged": Any(),

			// TODO(dennwc): might be useful later; drop it for now
			"IsVar": Any(),
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
				uast.KeyPos: Any(),

				// trivia == whitespace; can safely drop it
				"LeadingTrivia":  Any(),
				"TrailingTrivia": Any(),

				"IsMissing": Bool(false),

				// contains escaped value, we don't need it in canonical UAST
				"Text": Any(),

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

	MapSemantic("TrueLiteralExpression", uast.Bool{}, MapObj(
		Obj{
			"Token": Obj{
				uast.KeyType: String("TrueKeyword"),
				uast.KeyPos: Any(),

				// trivia == whitespace; can safely drop it
				"LeadingTrivia": Any(),
				"TrailingTrivia": Any(),

				"Text": String("true"),
				"Value": Bool(true),
				"ValueText": String("true"),

				"IsMissing": Bool(false),
			},
			"IsMissing": Bool(false),
			"IsStructuredTrivia": Bool(false),
		},
		Obj{
			"Value": Bool(true),
		},
	)),

	MapSemantic("FalseLiteralExpression", uast.Bool{}, MapObj(
		Obj{
			"Token": Obj{
				uast.KeyType: String("FalseKeyword"),
				uast.KeyPos: Any(),

				// trivia == whitespace; can safely drop it
				"LeadingTrivia": Any(),
				"TrailingTrivia": Any(),

				"Text": String("false"),
				"Value": Bool(false),
				"ValueText": String("false"),

				"IsMissing": Bool(false),
			},
			"IsMissing": Bool(false),
			"IsStructuredTrivia": Bool(false),
		},
		Obj{
			"Value": Bool(false),
		},
	)),

	MapSemantic("Block", uast.Block{}, MapObj(
		Obj{
			"Statements": Var("stmts"),
			// TODO(dennwc): remap to custom positional fields
			"OpenBraceToken":  Any(),
			"CloseBraceToken": Any(),
		},
		Obj{
			"Statements": Var("stmts"),
		},
	)),

	MapSemantic("SingleLineCommentTrivia", uast.Comment{}, MapObj(
		Obj{
			uast.KeyToken: CommentText([2]string{"//", ""}, "text"),
			"IsDirective": Bool(false),
		},
		CommentNode(false, "text", nil),
	)),
	// TODO(dennwc): differentiate from regular comments
	MapSemantic("SingleLineDocumentationCommentTrivia", uast.Comment{}, MapObj(
		Obj{
			uast.KeyToken: CommentText([2]string{"///", ""}, "text"),
			"IsDirective": Bool(false),
		},
		CommentNode(false, "text", nil),
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
			"SemicolonToken": Any(),
			"UsingKeyword":   Any(),
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
						uast.KeyPos: Any(),
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

	MapSemantic("Parameter", uast.Argument{}, MapObj(
		Obj{
			"Identifier": Obj{
					// Could be identifier or some keywords like __arg, same fields in any case
					uast.KeyType: Any(),
					// TODO: assert that is the same as in parent
					uast.KeyPos:      Var("id_pos"),
					"LeadingTrivia":  Any(),
					"TrailingTrivia": Any(),
					"IsMissing": Bool(false),
					"Text":      Var("name"),
					"Value":     Var("name"),
					"ValueText": Var("name"),
				},
			"AttributeLists":     Any(),
			"Default":            Var("def_init"),
			"IsMissing":          Any(),
			"IsStructuredTrivia": Any(),
			"Modifiers":          Cases("varargs",
				// FIXME: must be something like "Contains", not exact match
				Arr(),
				Arr(Obj{
					uast.KeyType: String("ParamsKeyword"),
					uast.KeyPos: Any(),
					"IsMissing": Bool(false),
					"LeadingTrivia": Any(),
					"TrailingTrivia": Any(),
					"Text": Any(),
					"Value": Any(),
					"ValueText": Any(),
				}),
				Any(),
			),
			"Type": Var("type"),
		},
		Obj{
			"Name": UASTType(uast.Identifier{}, Obj{
					"Name": Var("name"),
			}),
			"Type": Var("type"),
			"Init": Var("def_init"),
			"Variadic": Cases("varargs",
				Bool(false),
				Bool(true),
				Bool(false),
			),
			//"Variadic": Bool(false),
			"MapVariadic": Bool(false),
			"Receiver": Bool(false),
		},
	)),

	funcDefMap("MethodDeclaration"),
	// FIXME: these two are not matching rightly
	funcDefMap("ConstructorDeclaration"),
	funcDefMap("DestructorDeclaration"),
}

// dropNils accepts a array node, removes all nil values from it and passes it to
// a specified suboperation.
// It will not restore nil values when constructing nodes (not reversible).
type dropNils struct {
	op Op
}

func (op dropNils) Kinds() nodes.Kind {
	return nodes.KindArray
}

func (op dropNils) Check(st *State, n nodes.Node) (bool, error) {
	arr, ok := n.(nodes.Array)
	if !ok && n != nil {
		return false, nil
	}
	out := make(nodes.Array, 0, len(arr))
	for _, e := range arr {
		if e != nil {
			out = append(out, e)
		}
	}
	return op.op.Check(st, out)
}

func (op dropNils) Construct(st *State, n nodes.Node) (nodes.Node, error) {
	return op.op.Construct(st, n)
}

package normalizer

import (
	"errors"

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

var _ Op = opArrHasKeyword{}

type opArrHasKeyword struct {
	keyword string
	opHas   Op
	opRest  Op
}

func (op opArrHasKeyword) Kinds() nodes.Kind {
	return nodes.KindArray
}

func (op opArrHasKeyword) Check(st *State, n nodes.Node) (bool, error) {
	arr, ok := n.(nodes.Array)
	if !ok && arr != nil {
		return false, nil
	}
	// find a node with a specified type and drop if from array
	// the boolean flag that we pass to a sub-op will indicate
	// if we found it or not
	for i, n := range arr {
		obj, ok := n.(nodes.Object)
		if !ok {
			continue
		}
		v, ok := obj[uast.KeyType]
		if !ok {
			continue
		}
		typ, ok := v.(nodes.String)
		if !ok || string(typ) != op.keyword {
			continue
		}
		// found the keyword
		if ok, err := op.opHas.Check(st, nodes.Bool(true)); err != nil || !ok {
			return ok, err
		}
		rest := make(nodes.Array, len(arr)-1)
		copy(rest, arr[:i])
		copy(rest[i:], arr[i+1:])
		return op.opRest.Check(st, rest)
	}
	// not found, default to false
	if ok, err := op.opHas.Check(st, nodes.Bool(false)); err != nil || !ok {
		return ok, err
	}
	return op.opRest.Check(st, n)
}

func (op opArrHasKeyword) Construct(st *State, n nodes.Node) (nodes.Node, error) {
	// first, we will need to read the flag from sub-op
	// if it's false, we will just pass and array as-is
	// if it's true, we will synthesize and append a node to it

	v, err := op.opHas.Construct(st, nil)
	if err != nil {
		return nil, err
	}
	has, ok := v.(nodes.Bool)
	if !ok {
		return nil, ErrUnexpectedType.New(nodes.Bool(false), v)
	}
	n, err = op.opRest.Construct(st, n)
	if err != nil {
		return nil, err
	} else if !has {
		// pass as-is
		return n, nil
	}
	// synthesize the node

	// TODO(dennwc): synthesize the node once we care about reverse transform
	return n, nil
}

var _ Op = opArrToChain{}

type opArrToChain struct {
	opMods Op
	opType Op
	// TODO(dennwc): maybe whitelist only known modifiers? seen so far:
	//  			 - RefKeyword
	//				 - OutKeyword (we should move it to Returns)
}

func (op opArrToChain) Kinds() nodes.Kind {
	return nodes.KindObject
}

func (op opArrToChain) Check(st *State, n nodes.Node) (bool, error) {
	// we assert that the passed node is an object and start
	// checking the Type field recursively
	// if there is one, we will remove it from the "Type" field
	// from current node and append it to an array
	// and we repeat it recursively on the value of the "Type" field
	var mods nodes.Array

	// TODO(dennwc): implement when we will need a reversal
	if ok, err := op.opType.Check(st, n); err != nil || !ok {
		return ok, err
	}
	return op.opMods.Check(st, mods)
}

func (op opArrToChain) Construct(st *State, n nodes.Node) (nodes.Node, error) {
	// load two nodes:
	// - the first one is an array of modifiers (objects)
	// - the second one is a type node
	nd, err := op.opMods.Construct(st, nil)
	if err != nil {
		return nil, err
	}
	mods, ok := nd.(nodes.Array)
	if !ok {
		return nil, ErrUnexpectedType.New(nodes.Array{}, nd)
	}
	typ, err := op.opType.Construct(st, n)
	if err != nil {
		return nil, err
	}
	// we will now use each modifier to construct a chain or a linked list of nodes
	// by adding a "Type" field to each modifier, that will point to the current node
	for _, nd := range mods {
		mod, ok := nd.(nodes.Object)
		if !ok {
			return nil, ErrUnexpectedType.New(nodes.Object{}, nd)
		}
		mod = mod.CloneObject()
		if _, ok := mod["Type"]; ok {
			return nil, errors.New("unexpected field in modifier: Type")
		}
		mod["Type"] = typ
		typ = mod
	}
	return typ, nil
}

func funcDefMap(typ string) Mapping {
	return MapSemantic(typ, uast.FunctionGroup{}, MapObj(
		Fields{
			{Name: "Body", Op: Var("body")},
			{Name: "Identifier", Op: Var("name")},
			{Name: "ParameterList", Op: Obj{
				uast.KeyType:         String("ParameterList"),
				uast.KeyPos:          Any(),
				"OpenParenToken":     Any(),
				"CloseParenToken":    Any(),
				"IsMissing":          Bool(false),
				"IsStructuredTrivia": Bool(false),
				"Parameters":         Var("params"),
			}},
			{Name: "ReturnType", Optional: "optReturn", Op: Var("rettype")},
		},

		Obj{
			"Nodes": Arr(
				UASTType(uast.Alias{}, Obj{
					"Name": Var("name"),
					"Node": UASTType(uast.Function{}, Obj{
						"Body": Var("body"),
						"Type": UASTType(uast.FunctionType{}, Fields{
							{Name: "Arguments", Op: Var("params")},
							{Name: "Returns", Optional: "optReturn", Op: Arr(
								UASTType(uast.Argument{}, Obj{
									"Type": Var("rettype"),
								}))},
						}),
					}),
				}),
			),
		},
	))
}

// useFullSpan is a set of node types that use FullSpan for positions instead of Span
var useFullSpan = []nodes.Value{
	nodes.String("SingleLineDocumentationCommentTrivia"),
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

	// Positional info is stored in a child node in Span field.
	//
	// This is not supported by ObjectToNode helper, and we are
	// too lazy to create positional node by hand.
	//
	// Instead, we will temporary remap positional info to
	// "spanStart" and "spanEnd" fields of the root node, and
	// ObjectToNode will pick them up later to build a proper
	// positional node.
	//
	// There is also a FullSpan field that includes leading/trailing
	// whitespaces and sometimes node tokens. We ignore this second
	// position for most nodes, but there are few exceptions where
	// we use FullSpan and ignore Span.
	Map(
		Part("_", CasesObj("case", Obj{}, Objs{
			// exceptions - use FullSpan
			{
				uast.KeyType: Check(
					In(useFullSpan...),
					Var("typ"),
				),
				"FullSpan": Obj{
					uast.KeyType: String("TextSpan"),
					"Length":     Any(),
					"Start":      Var("start"),
					"End":        Var("end"),
				},
				// TODO(dennwc): add it as a custom position field?
				"Span": Any(),
			},
			// other nodes - use Span
			{
				uast.KeyType: Check(
					Not(In(useFullSpan...)),
					Var("typ"),
				),
				"Span": Obj{
					uast.KeyType: String("TextSpan"),
					"Length":     Any(),
					"Start":      Var("start"),
					"End":        Var("end"),
				},
				// TODO(dennwc): add it as a custom position field?
				"FullSpan": Any(),
			},
		})),
		Part("_", CasesObj("case", Obj{
			// remap to temporary keys and let ObjectToNode to pick them up
			"spanStart": Var("start"),
			"spanEnd":   Var("end"),
		}, Objs{
			// exceptions
			{
				uast.KeyType: Check(
					In(useFullSpan...),
					Var("typ"),
				),
			},
			// other nodes
			{
				uast.KeyType: Check(
					Not(In(useFullSpan...)),
					Var("typ"),
				),
			},
		})),
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
	Map(
		Part("_", Obj{
			uast.KeyType: String("MultiLineCommentTrivia"),
		}),
		Part("_", Obj{
			uast.KeyType:  String("MultiLineCommentTrivia"),
			uast.KeyToken: String(""),
		}),
	),
}

// Normalizers is the main block of normalization rules to convert native AST to semantic UAST.
var Normalizers = []Mapping{

	// remove empty identifier tokens
	Map(
		Check(
			Has{
				uast.KeyType: String("IdentifierToken"),
				// make sure it's empty, we don't want to wipe something useful
				"Text":      String(""),
				"Value":     String(""),
				"ValueText": String(""),
			},
			Any(),
		),
		Is(nil),
	),

	MapSemantic("IdentifierToken", uast.Identifier{}, MapObj(
		Obj{
			// trivia == whitespace; can safely drop it
			"LeadingTrivia":  Any(),
			"TrailingTrivia": Any(),
			"IsMissing":      Bool(false),

			// we drop this one, because C# allows to declare
			// a "for" identifier by using "@for" notation
			// and we don't need that token in Semantic mode
			"Text": Any(),
			// all other token values are the same
			"Value":     Var("name"),
			"ValueText": Var("name"),
		}, Obj{
			"Name": Var("name"),
		},
	)),

	// remove empty identifiers
	Map(
		Check(
			Has{
				uast.KeyType: String("IdentifierName"),
				"Identifier": Is(nil),
			},
			Any(),
		),
		Is(nil),
	),

	Map(
		Obj{
			uast.KeyType: String("IdentifierName"),
			uast.KeyPos:  Any(), // TODO(dennwc): assert that it's the same

			"Identifier": Var("ident"),

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
		Var("ident"),
	),

	// Special: is a keyword, but used as an identifier (Parameter name)
	MapSemantic("ArgListKeyword", uast.Identifier{}, MapObj(
		Obj{
			// trivia == whitespace; can safely drop it
			"LeadingTrivia":  Any(),
			"TrailingTrivia": Any(),

			"IsMissing": Bool(false),

			// all token values are the same
			"Text":      String("__arglist"),
			"Value":     String("__arglist"),
			"ValueText": String("__arglist"),
		}, Obj{
			"Name": String("__arglist"),
		},
	)),

	MapSemantic("StringLiteralExpression", uast.String{}, MapObj(
		Obj{
			"Token": Obj{
				uast.KeyType: String("StringLiteralToken"),
				uast.KeyPos:  Any(),

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
				uast.KeyPos:  Any(),

				// trivia == whitespace; can safely drop it
				"LeadingTrivia":  Any(),
				"TrailingTrivia": Any(),

				"Text":      String("true"),
				"Value":     Bool(true),
				"ValueText": String("true"),

				"IsMissing": Bool(false),
			},
			"IsMissing":          Bool(false),
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
				uast.KeyPos:  Any(),

				// trivia == whitespace; can safely drop it
				"LeadingTrivia":  Any(),
				"TrailingTrivia": Any(),

				"Text":      String("false"),
				"Value":     Bool(false),
				"ValueText": String("false"),

				"IsMissing": Bool(false),
			},
			"IsMissing":          Bool(false),
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

	MapSemantic("MultiLineCommentTrivia", uast.Comment{}, MapObj(
		Obj{
			uast.KeyToken: CommentText([2]string{"/*", "*/"}, uast.KeyToken),
			"IsDirective": Bool(false),
		},
		CommentNode(true, uast.KeyToken, nil),
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

	// Old style multiple arguments: argument with the magic name "__arglist"
	MapSemantic("Parameter", uast.Argument{}, MapObj(
		Obj{
			"Identifier": Check(
				Has{
					uast.KeyType: String(uast.TypeOf(uast.Identifier{})),
					"Name":       String("__arglist"),
				}, Var("name")),
			"AttributeLists":     Arr(), // TODO(dennwc): any cases when it's not empty?
			"Default":            Var("def_init"),
			"IsMissing":          Bool(false),
			"IsStructuredTrivia": Bool(false),
			"Modifiers":          Arr(), // TODO(dennwc): any cases when it's not empty?
			"Type":               Var("type"),
		},
		Obj{
			"Name":        Var("name"),
			"Type":        Var("type"),
			"Init":        Var("def_init"),
			"Variadic":    Bool(true),
			"MapVariadic": Bool(false),
			"Receiver":    Bool(false),
		},
	)),

	// Normal parameter, potential multiple args expressed by "params" in modifiers
	MapSemantic("Parameter", uast.Argument{}, MapObj(
		Obj{
			"Identifier": Check(Has{
				uast.KeyType: String(uast.TypeOf(uast.Identifier{})),
			}, Var("name")),
			"AttributeLists":     Arr(), // TODO(dennwc): any cases when it's not empty?
			"Default":            Var("def_init"),
			"IsMissing":          Bool(false),
			"IsStructuredTrivia": Any(),
			"Modifiers": opArrHasKeyword{
				keyword: "ParamsKeyword",
				opHas:   Var("variadic"),
				opRest: opArrHasKeyword{
					keyword: "ThisKeyword",
					opHas:   Var("this"),
					opRest:  Var("rest"),
				},
			},
			"Type": Var("type"),
		},
		Obj{
			"Name": Var("name"),
			"Type": opArrToChain{
				opMods: Var("rest"),
				opType: Var("type"),
			},
			"Init":        Var("def_init"),
			"Variadic":    Var("variadic"),
			"MapVariadic": Bool(false),
			"Receiver":    Var("this"),
		},
	)),

	funcDefMap("MethodDeclaration"),
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

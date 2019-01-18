package normalizer

import (
	"gopkg.in/bblfsh/sdk.v2/uast"
	"gopkg.in/bblfsh/sdk.v2/uast/role"

	. "gopkg.in/bblfsh/sdk.v2/uast/transformer"
	"gopkg.in/bblfsh/sdk.v2/uast/transformer/positioner"
)

// Native is the of list `transformer.Transformer` to apply to a native AST.
// To learn more about the Transformers and the available ones take a look to:
// https://godoc.org/gopkg.in/bblfsh/sdk.v2/uast/transformer
var Native = Transformers([][]Transformer{
	// The main block of transformation rules.
	{Mappings(Annotations...)},
	{
		// RolesDedup is used to remove duplicate roles assigned by multiple
		// transformation rules.
		RolesDedup(),
	},
}...)

// PreprocessCode is a special block of transformations that are applied at the end
// and can access original source code file. It can be used to improve or
// fix positional information.
//
// https://godoc.org/gopkg.in/bblfsh/sdk.v2/uast/transformer/positioner
var PreprocessCode = []CodeTransformer{
	positioner.FromOffset(),
	positioner.TokenFromSource{
		Types: []string{
			"SingleLineCommentTrivia",
			"SingleLineDocumentationCommentTrivia",
		},
	},
}

var Code []CodeTransformer // legacy stage, will be deprecated

// Annotations is a list of individual transformations to annotate a native AST with roles.
var Annotations = []Mapping{
	AnnotateType("internal-type", nil, role.Incomplete),
	AnnotateType("CompilationUnit", nil, role.File, role.Module),
	AnnotateType("Block", nil, role.Block),
	AnnotateType("NamespaceDeclaration", nil, role.Block, role.Scope),
	AnnotateType("ArrayType", nil, role.List, role.Incomplete),
	AnnotateType("ArrayRankSpecifier", nil, role.List, role.Incomplete),
	AnnotateType("BracketedArgumentList", nil, role.List, role.Value, role.Incomplete), // i in someaArray[i]
	AnnotateType("OmittedArraySizeExpression", nil, role.List, role.Expression, role.Incomplete),
	AnnotateType("ArrayCreationExpression", nil, role.List, role.Expression, role.Instance, role.Incomplete),
	AnnotateType("ElementAccessExpression", nil, role.List, role.Value, role.Incomplete),
	AnnotateType("CastExpression", nil, role.Expression, role.Incomplete),
	AnnotateType("PredefinedType", nil, role.Incomplete, role.Declaration, role.Variable),
	AnnotateType("GenericName", nil, role.Identifier, role.Incomplete), // FIXME: get the role from the IdentifierToken child
	AnnotateType("TypeArgumentList", nil, role.Argument, role.List, role.Instance, role.Incomplete), // generic <T,U> types on instantiation
	AnnotateType("TypeParameterList", nil, role.Argument, role.List, role.Incomplete), // generic <T,U> types on specification
	AnnotateType("TypeParameter", nil, role.Argument, role.Incomplete),
	AnnotateType("ObjectCreationExpression", nil, role.Type, role.Instance),
	AnnotateType("CollectionInitializerExpression", nil, role.Incomplete, role.Value),
	AnnotateType("SimpleAssignmentExpression", nil, role.Assignment, role.Expression),

	AnnotateType("NumericLiteralExpression", nil, role.Expression, role.Number, role.Literal),
	AnnotateType("CharacterLiteralExpression", nil, role.Expression, role.Character, role.Literal),
	AnnotateType("StringLiteralExpression", nil, role.Literal, role.String, role.Expression),
	AnnotateType("TrueLiteralExpression", nil, role.Literal, role.Boolean, role.Expression),
	AnnotateType("FalseLiteralExpression", nil, role.Literal, role.Boolean, role.Expression),

	AnnotateType("QueryExpression", nil, role.Expression, role.Incomplete),
	AnnotateType("QueryBody", nil, role.Expression, role.Body, role.Incomplete),
	AnnotateType("SelectClause", nil, role.Expression, role.Incomplete),
	AnnotateType("WhereClause", nil, role.Expression, role.Incomplete),
	AnnotateType("FromClause", nil, role.Expression, role.Incomplete),

	AnnotateType("InvocationExpression", nil, role.Function, role.Call),
	AnnotateType("ArgumentList", nil, role.Function, role.Call, role.Argument, role.List),
	AnnotateType("Argument", nil, role.Function, role.Call, role.Argument),
	AnnotateType("SimpleMemberAccessExpression", nil, role.Qualified),

	AnnotateType("AmpersandAmpersandToken", nil, role.Operator, role.Relational, role.And),
	AnnotateType("AmpersandEqualsToken", nil, role.Operator, role.And, role.Equal),
	AnnotateType("AmpersandToken", nil, role.Operator, role.Bitwise, role.And),
	AnnotateType("AsteriskEqualsToken", nil, role.Operator, role.Arithmetic, role.Multiply, role.Equal),
	AnnotateType("AsteriskToken", nil, role.Operator, role.Arithmetic, role.Multiply),
	AnnotateType("BarBarToken", nil, role.Operator, role.Relational, role.Or),
	AnnotateType("BarEqualsToken", nil, role.Operator, role.Bitwise, role.Or, role.Equal),
	AnnotateType("BarToken", nil, role.Operator, role.Bitwise, role.Or),
	AnnotateType("CaretEqualsToken", nil, role.Operator, role.Bitwise, role.Xor, role.Equal),
	AnnotateType("CaretToken", nil, role.Operator, role.Bitwise, role.Xor),
	AnnotateType("CloseBraceToken", nil, role.Incomplete),
	AnnotateType("CloseBracketToken", nil, role.Incomplete),
	AnnotateType("CloseParenToken", nil, role.Incomplete),
	AnnotateType("ColonToken", nil, role.Incomplete),
	AnnotateType("CommaToken", nil, role.Incomplete),
	AnnotateType("DotToken", nil, role.Incomplete),
	AnnotateType("EndOfFileToken", nil, role.Incomplete),
	AnnotateType("EqualsEqualsToken", nil, role.Operator, role.Relational, role.Equal),
	AnnotateType("EqualsGreaterThanToken", nil, role.Operator, role.Relational, role.GreaterThanOrEqual),
	AnnotateType("EqualsToken", nil, role.Operator, role.Equal),
	AnnotateType("ExclamationEqualsToken", nil, role.Operator, role.Relational, role.Not, role.Equal),
	AnnotateType("ExclamationToken", nil, role.Operator, role.Not),
	AnnotateType("GreaterThanEqualsToken", nil, role.Operator, role.Relational, role.GreaterThanOrEqual),
	AnnotateType("GreaterThanGreaterThanEqualsToken", nil, role.Operator, role.Bitwise, role.RightShift, role.Equal),
	AnnotateType("GreaterThanGreaterThanToken", nil, role.Operator, role.Bitwise, role.RightShift),
	AnnotateType("GreaterThanToken", nil, role.Operator, role.Relational, role.GreaterThan),
	AnnotateType("InterpolatedStringEndToken", nil, role.Incomplete, role.String),
	AnnotateType("InterpolatedStringStartToken", nil, role.Incomplete, role.String),
	AnnotateType("InterpolatedStringTextToken", nil, role.Incomplete, role.String),
	AnnotateType("LessThanEqualsToken", nil, role.Operator, role.Relational, role.LessThanOrEqual),
	AnnotateType("LessThanLessThanEqualsToken", nil, role.Operator, role.Bitwise, role.LeftShift, role.Equal),
	AnnotateType("LessThanLessThanToken", nil, role.Operator, role.Bitwise, role.LeftShift),
	AnnotateType("LessThanToken", nil, role.Operator, role.Relational, role.LessThan),
	AnnotateType("MinusEqualsToken", nil, role.Operator, role.Arithmetic, role.Substract, role.Equal),
	AnnotateType("MinusGreaterThanToken", nil, role.Operator, role.Dereference),
	AnnotateType("MinusMinusToken", nil, role.Operator, role.Unary, role.Decrement),
	AnnotateType("MinusToken", nil, role.Operator, role.Arithmetic, role.Substract),
	AnnotateType("OmittedArraySizeExpressionToken", nil, role.Incomplete),
	AnnotateType("OpenBraceToken", nil, role.Incomplete),
	AnnotateType("OpenBracketToken", nil, role.Incomplete),
	AnnotateType("OpenParenToken", nil, role.Incomplete),
	AnnotateType("PercentEqualsToken", nil, role.Operator, role.Arithmetic, role.Modulo, role.Equal),
	AnnotateType("PercentToken", nil, role.Operator, role.Arithmetic, role.Modulo),
	AnnotateType("PlusEqualsToken", nil, role.Operator, role.Arithmetic, role.Add, role.Equal),
	AnnotateType("PlusPlusToken", nil, role.Unary, role.Arithmetic, role.Increment),
	AnnotateType("PlusToken", nil, role.Operator, role.Arithmetic, role.Substract),
	AnnotateType("QuestionQuestionToken", nil, role.Operator, role.Incomplete),
	AnnotateType("QuestionToken", nil, role.Operator, role.Incomplete),
	AnnotateType("SemicolonToken", nil, role.Incomplete),
	AnnotateType("SlashEqualsToken", nil, role.Operator, role.Arithmetic, role.Divide, role.Equal),
	AnnotateType("SlashToken", nil, role.Operator, role.Arithmetic, role.Divide),
	AnnotateType("TildeToken", nil, role.Operator, role.Unary, role.Bitwise, role.Not),

	// We probably need a role.Keyword for languages like this that add a specific node for them
	AnnotateType("None", nil, role.Incomplete), // e.g. SemiColonField or lines not ended in ;
	AnnotateType("UsingKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Import, role.Incomplete),
	AnnotateType("AbstractKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("AddKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("AsKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("AscendingKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("AssemblyKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("AsyncKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("AwaitKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("BaseKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("BoolKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Boolean, role.Declaration),
	AnnotateType("BreakKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Break),
	AnnotateType("ByteKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Number, role.Declaration),
	AnnotateType("CaseKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Switch, role.Case),
	AnnotateType("CharKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Character, role.Declaration),
	AnnotateType("CheckedKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("ClassKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Type, role.Declaration),
	AnnotateType("ConstKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Declaration, role.Incomplete),
	AnnotateType("DecimalKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Number, role.Declaration),
	AnnotateType("DefaultKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Default),
	AnnotateType("DelegateKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("DoKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("DoubleKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Number, role.Declaration),
	AnnotateType("ElseKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Else),
	AnnotateType("EnumKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("EventKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("ExplicitKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("ExternKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("FixedKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("FloatKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Number, role.Declaration),
	AnnotateType("ForKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.For),
	AnnotateType("FromKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("GetKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("GotoKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Goto),
	AnnotateType("IfKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.If),
	AnnotateType("ImplicitKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("InKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("IntKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Number, role.Declaration),
	AnnotateType("InterfaceKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("InternalKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("IsKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("LockKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("LongKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Number, role.Declaration),
	AnnotateType("NamespaceKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Block),
	AnnotateType("NewKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Instance),
	AnnotateType("NullKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Null, role.Literal),
	AnnotateType("ObjectKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Type, role.Incomplete),
	AnnotateType("OperatorKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("OrderByKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("OverrideKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("ParamsKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("PartialKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("PrivateKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Visibility, role.Instance),
	AnnotateType("ProtectedKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Visibility, role.Subtype),
	AnnotateType("PublicKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Visibility, role.World),
	AnnotateType("ReadOnlyKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("RemoveKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("ReturnKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Return),
	AnnotateType("SByteKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Number, role.Declaration),
	AnnotateType("SealedKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("SelectKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("SetKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("ShortKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Number, role.Declaration),
	AnnotateType("SizeOfKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("StackAllocKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("StaticKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("StringKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.String, role.Declaration),
	AnnotateType("StructKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Type, role.Declaration),
	AnnotateType("SwitchKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Switch),
	AnnotateType("ThisKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("ThrowKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("TypeOfKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("UIntKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Number, role.Declaration),
	AnnotateType("ULongKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Number, role.Declaration),
	AnnotateType("UShortKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Number, role.Declaration),
	AnnotateType("UncheckedKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("UnsafeKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("UsingKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("VirtualKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("VoidKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("VolatileKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("WhereKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Incomplete),
	AnnotateType("WhileKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.While),
	AnnotateType("YieldKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Return, role.Incomplete),

	AnnotateType("IdentifierToken", FieldRoles{"Value": {Rename: uast.KeyToken}},
		role.Identifier, role.Expression),

	AnnotateType("BinaryExpression_BitwiseAndExpression", nil, role.Binary, role.Expression,
		role.Bitwise, role.And),
	AnnotateType("BinaryExpression_BitwiseOrExpression", nil, role.Binary, role.Expression,
		role.Bitwise, role.And),
	AnnotateType("BinaryExpression_BitwiseNotExpression", nil, role.Binary, role.Expression,
		role.Bitwise, role.Not),
	AnnotateType("BinaryExpression_ExclusiveOrExpression", nil, role.Binary, role.Expression,
		role.Bitwise, role.And),
	AnnotateType("BinaryExpression_LeftShiftExpression", nil, role.Binary, role.Expression,
		role.Bitwise, role.LeftShift),
	AnnotateType("BinaryExpression_RightShiftExpression", nil, role.Binary, role.Expression,
		role.Bitwise, role.LeftShift),
	AnnotateType("PrefixUnaryExpression_BitwiseNotExpression", nil, role.Unary, role.Expression,
		role.Bitwise, role.Not),

	AnnotateType("BinaryExpression_LessThanExpression", nil, role.Binary, role.Expression,
		role.Relational, role.LessThan),
	AnnotateType("BinaryExpression_NotEqualsExpression", nil, role.Binary, role.Expression,
		role.Relational, role.Not, role.Equal),
	AnnotateType("BinaryExpression_LessThanOrEqualExpression", nil, role.Binary, role.Expression,
		role.Relational, role.LessThanOrEqual),
	AnnotateType("BinaryExpression_LogicalAndExpression", nil, role.Binary, role.Expression,
		role.Relational, role.Or),
	AnnotateType("BinaryExpression_LogicalOrExpression", nil, role.Binary, role.Expression,
		role.Relational, role.Or),
	AnnotateType("BinaryExpression_EqualsExpression", nil, role.Binary, role.Expression,
		role.Relational, role.Equal),

	AnnotateType("BinaryExpression_AddExpression", nil, role.Binary, role.Expression,
		role.Arithmetic, role.Add),
	AnnotateType("BinaryExpression_SubtractExpression", nil, role.Binary, role.Expression,
		role.Arithmetic, role.Substract),
	AnnotateType("PrefixUnaryExpression_PostIncrementExpression", nil, role.Unary, role.Expression,
		role.Arithmetic, role.Increment),
	AnnotateType("PrefixUnaryExpression_PostDecrementExpression", nil, role.Unary, role.Expression,
		role.Arithmetic, role.Decrement),
	AnnotateType("PostfixUnaryExpression_PostIncrementExpression", nil, role.Unary, role.Expression,
		role.Arithmetic, role.Increment),
	AnnotateType("PostfixUnaryExpression_PostDecrementExpression", nil, role.Unary, role.Expression,
		role.Arithmetic, role.Decrement),

	AnnotateType("ParenthesizedExpression", nil, role.Expression),
	AnnotateType("LocalDeclarationStatement", nil, role.Declaration, role.Expression),
	AnnotateType("VariableDeclaration", nil, role.Declaration, role.Variable, role.Expression),
	AnnotateType("VariableDeclarator", nil, role.Declaration, role.Variable, role.Right),
	AnnotateType("NameEquals", nil, role.Assignment, role.Right),
	AnnotateType("EqualsValueClause", nil, role.Assignment, role.Right),
	AnnotateType("ExpressionStatement", nil, role.Expression, role.Statement),

	AnnotateType("NumericLiteralToken", nil, role.Value, role.Number, role.Literal),
	AnnotateType("CharacterLiteralToken", nil, role.Literal, role.Character),
	AnnotateType("StringLiteralToken", nil, role.Literal, role.String),
	AnnotateType("TrueKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Boolean, role.Literal),
	AnnotateType("FalseKeyword", FieldRoles{"Value": {Rename: uast.KeyToken}}, role.Boolean, role.Literal),

	AnnotateType("ClassDeclaration", nil, role.Type, role.Declaration), // FIXME: incomplete
	AnnotateType("FieldDeclaration", nil, role.Type, role.Declaration, role.Variable),
	AnnotateType("MethodDeclaration", nil, role.Type, role.Function, role.Declaration), // FIXME: incomplete
	AnnotateType("UsingDirective", nil, role.Import, role.Statement),
	AnnotateType("IdentifierName", nil, role.Identifier), // FIXME: get the token from the IdentifierToken child
	AnnotateType("ParameterList", nil, role.Function, role.Declaration, role.Argument, role.List),
	AnnotateType("Parameter", nil, role.Function, role.Declaration, role.Argument),
	AnnotateType("ReturnStatement", nil, role.Statement, role.Return),
	AnnotateType("SimpleLambdaExpression", nil, role.Function, role.Declaration, role.Anonymous, role.Expression),

	// FIXME: incomplete
	AnnotateType("WhileStatement", nil, role.While, role.Statement),
	AnnotateType("ForStatement", nil, role.For, role.Statement),
	AnnotateType("IfStatement", nil, role.If, role.Statement),
}

using System;
using System.Linq;
using System.Collections.Generic;

using Newtonsoft.Json;
using Newtonsoft.Json.Serialization;
using Microsoft.CodeAnalysis;
using Microsoft.CodeAnalysis.CSharp;

namespace native
{
    public class ParseRequest
    {
        public string content;
    }

    public class ParseResponse
    {
        public string status;
        public List<string> errors;
        public Object ast;
    }

    class Program
    {
        static void Main(string[] args)
        {
            var jsonSerializerSettings = new JsonSerializerSettings
            {
                PreserveReferencesHandling = PreserveReferencesHandling.None,
                // ignore loops
                ReferenceLoopHandling = ReferenceLoopHandling.Ignore,
                // controls how individual fields are converted
                ContractResolver = new ASTContractResolver(),
            };
            var jsonSerializer = JsonSerializer.Create(jsonSerializerSettings);
            var jsonWriter = new JsonTextWriter(Console.Out);

            string line;
            while ((line = Console.ReadLine()) != null)
            {
                // TODO(dennwc): handle exceptions and syntax errors
                ParseRequest req = JsonConvert.DeserializeObject<ParseRequest>(line);

                Object ast = Parse(req.content);

                ParseResponse resp = new ParseResponse
                {
                    status = "ok",
                    ast = ast,
                };
                jsonSerializer.Serialize(jsonWriter, resp);
                jsonWriter.WriteWhitespace("\n");
                jsonWriter.Flush();
            }
        }

        static Object Parse(string source)
        {
            SyntaxTree tree = CSharpSyntaxTree.ParseText(source);
            var cstree = (CSharpSyntaxTree)tree;
            return cstree.GetRoot();
        }
    }

    class ASTContractResolver : DefaultContractResolver
    {
        protected override IList<JsonProperty> CreateProperties(Type type, MemberSerialization memberSerialization)
        {
            IList<JsonProperty> properties = base.CreateProperties(type, memberSerialization);

            bool hasRawKind = false;

            // drop properties that we won't use
            properties = properties.Where((p) => {
                switch (p.PropertyName)
                {
                // is set nearly for every node; always "C#"
                case "Language":
                // reference from tokens to the corresponding SyntaxTree root
                // always the same in every node
                case "SyntaxTree":
                // don't need any parent references
                case "ParentTrivia":

                    return false;
                case "RawKind":
                    // RawKind stores a specific node type as an int enum value
                    // we drop this field and merge it with virtual @type field
                    hasRawKind = true;
                    return false;
                default:
                    // don't need those Contains<FieldName> and Has<FieldName> flags
                    // can check the field value directly
                    if (p.PropertyName.StartsWith("Contains") ||
                        p.PropertyName.StartsWith("Has")) {
                        return false;
                    }
                    return true;
                }
            }).ToList();

            // add a virtual @type property
            properties.Add(new JsonProperty()
            {
                PropertyName = "@type",
                PropertyType = typeof(string),
                Readable = true,
                Writable = false,
                ValueProvider = new TypeValueProvider(type, hasRawKind)
            });

            return properties;
        }
    }

    class TypeValueProvider : IValueProvider {
        Type _type;
        bool _hasKind;
        public TypeValueProvider(Type v, bool hasKind = false)
        {
            _type = v;
            _hasKind = hasKind;
        }
        public Object GetValue(Object target)
        {
            string name = _type.Name;
            // most nodes types contain "Syntax" prefix, so we trim it
            if (name.EndsWith("Syntax"))
            {
                name = name.Substring(0, name.Length - 6);
            }
            if (!_hasKind)
            {
                return name;
            }
            // RawKind gives a more specific type name for AST nodes
            int kind = (int)_type.GetProperty("RawKind").GetValue(target, null);
            string skind = Enum.GetName(typeof(SyntaxKind), kind);

            if (skind == name)
            {
                return name;
            }
            // for tokens and trivias the RawKind alone is descriptive enough
            else if (name == "SyntaxToken" || name == "SyntaxTrivia")
            {
                return skind;
            }
            // some RawKinds end with the type name
            if (skind.EndsWith(name))
            {
                return skind;
            }
            return name + "_" + skind;
        }
        public void SetValue(Object target, Object value)
        {
            // do nothing
        }
    }
}
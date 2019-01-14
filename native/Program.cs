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
                string json = JsonConvert.SerializeObject(resp, jsonSerializerSettings);
                Console.WriteLine(json);
            }
        }

        static Object Parse(string source)
        {
            SyntaxTree tree = CSharpSyntaxTree.ParseText(source);
            var cstree = (CSharpSyntaxTree)tree;
            return cstree.GetRoot();
        }
    }

    public class ASTContractResolver : DefaultContractResolver
    {
        protected override IList<JsonProperty> CreateProperties(Type type, MemberSerialization memberSerialization)
        {
            IList<JsonProperty> properties = base.CreateProperties(type, memberSerialization);

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
                ValueProvider = new StringValueProvider(type.Name)
            });

            return properties;
        }
    }

    public class StringValueProvider : IValueProvider {
        string value;
        public StringValueProvider(string v)
        {
            value = v;
        }
        public Object GetValue(Object target)
        {
            return value;
        }
        public void SetValue(Object target, Object value)
        {
            // do nothing
        }
    }
}
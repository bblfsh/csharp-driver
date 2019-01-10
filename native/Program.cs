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
                ReferenceLoopHandling = ReferenceLoopHandling.Ignore,
                DefaultValueHandling = DefaultValueHandling.Ignore,
                ContractResolver = new ASTContractResolver(),
            };

            string line;
            while ((line = Console.ReadLine()) != null)
            {
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

            properties = properties.Where((p) => {
                switch (p.PropertyName)
                {
                case "Language":
                case "SyntaxTree":
                    return false;
                default:
                    return true;
                }
            }).ToList();

            return properties;
        }
    }
}
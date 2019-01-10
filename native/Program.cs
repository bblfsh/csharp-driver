using System;
using System.Collections.Generic;

using Newtonsoft.Json;

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
                string json = JsonConvert.SerializeObject(resp);
                Console.Write(json);
            }
        }

        static Object Parse(string source)
        {
            return source; // TODO
        }
    }
}
using System;

namespace native
{
    class Program
    {
        static void Main(string[] args)
        {
            string line;
            while ((line = Console.ReadLine()) != null)
            {
                Console.Write("{\"status\":\"ok\", \"ast\": "+line+"}\n");
            }
        }
    }
}
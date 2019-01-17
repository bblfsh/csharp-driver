delegate void Printer(string s);

class TestClass
{
    static void Main()
    {
        Printer p = delegate(string j)
        {
            System.Console.WriteLine(j);
        };

        p("The delegate using the anonymous method is called.");
        p = DoWork;
        p("The delegate using the named method is called.");
    }

    static void DoWork(string k)
    {
        System.Console.WriteLine(k);
    }
}

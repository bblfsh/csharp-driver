class Program {
    static void Main(string[] args) {
        for (int i = 0; i < 5; i++) {
            Console.WriteLine(i);
        }

        int i;
        int j = 10;
        for (i = 0, Console.WriteLine("foo"); i < j; i++, j--, Console.WriteLine("bar")) {}

        for (;;) {
        }
    }
}

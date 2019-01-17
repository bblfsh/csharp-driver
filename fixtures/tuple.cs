class Program {
    static void Main(string[] args) {
        var unnamed = ("one", "two");
        string s = unnamed.Item1;
        var named = (first: "one", second: "two");
        s = named.first;
        (int, (int, int)) nestedTuple = (1, (2, 3));
        (int first, int second) = named;
        unnamed = named;
    }

    public static (int Count, double Sum) Returntuple() {
        return (1, 3.0);
    }

    private static (int, double) ReturnUnnamedTuple() {
        return (2, 4.0);
    }
}

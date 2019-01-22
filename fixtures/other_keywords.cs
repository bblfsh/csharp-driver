class Program {
    public void TakeVariableLengthArgs(__arglist)
    {
        var args = new ArgIterator(__arglist);
    }

    static void Main(string[] args) {
        int x = 1;
        var xref = __makeref(x);
        __reftype(xref);
        __refvalue(xref);
        TakeVariableLengthArgs(__arglist(new StringBuilder(), 12));
    }
}

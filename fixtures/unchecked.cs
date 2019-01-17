class Program {
    static void Main(string[] args) {
        int ten = 10;

        unchecked 
        {
            int overflow = 2147483647 + ten;
        }

        int overflow2 = unchecked(2147483647 + ten);
    }
}

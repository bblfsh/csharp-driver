class Program {
    static void Main(string[] args) {
        int ten = 10;

        checked 
        {
            int overflow = 2147483647 + ten;
        }

        int overflow2 = checked(2147483647 + ten);
    }
}

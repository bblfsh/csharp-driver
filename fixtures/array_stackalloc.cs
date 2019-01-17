class Program {
    static void unsafe Main(string[] args) {
        int* foo = stackalloc int[] {1, 2, 3, 4};
        int* bar = stackalloc int[10];
        bar[0];
    }
}

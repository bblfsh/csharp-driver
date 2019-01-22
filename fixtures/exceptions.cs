class Program {
    static void Main(string[] args) {
        int a;
        try {
            a = 1;
        } catch (System.Exception ex) {
            a = 2;
        } finally {
            a = 3;
        }
    }
}

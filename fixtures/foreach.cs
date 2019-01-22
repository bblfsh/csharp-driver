
class Program {
    static void Main(string[] args) {
        var fibNumbers = new List<int> { 0, 1, 1, 2, 3, 5, 8, 13 };
        foreach (int element in fibNumbers) {
            Console.WriteLine(element);
        }

        foreach(ref int i in fibNumbers) {
            ++i;
        }

        foreach(ref readonly j in fibNumbers) {
            Console.WriteLine(j);
        }
    }
}

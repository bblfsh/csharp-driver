class Program {
    static void Main(string[] args) {
        throw new IndexOutOfRangeException();
    }

    string Name;
    void ThrowExp(string name) => Name = name ?? throw new ArgumentNullException(name);
}

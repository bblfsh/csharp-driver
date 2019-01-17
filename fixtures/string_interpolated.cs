class Program {
    static void Main(string[] args) {
        string name = "Blas";
        string s1 = $"What is your {name}?";
        string s2 = $"What is your {name}? {{ }}";
        string s3 = $"|{"Left", -7}|{"Right", 7}";
        string s4 = $"{Math.PI,FieldWidthRightAligned:F3};
    }
}

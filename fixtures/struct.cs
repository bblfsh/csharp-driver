public struct Coords
{
    public int x, y;

    public Coords(int p1, int p2)
    {
        x = p1;
        y = p2;
    }
}

class Program {
    static void Main(string[] args) {
        Coords c = new Coords();
        Coords c2 = new Coords(1, 2);
    }
}

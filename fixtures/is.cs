class C1 {}
class C2: C1 {}

const TEST = 100;

class Program {
    static void Main(string[] args) {
        C2 c2 = new C2();
        c2 is C1;

        int t = 100;
        if (t is TEST) {}

        string s;
        if (s is null) {}

        if (s is var x) {}
    }

    void Compare(Object o) {
        if (o is C1 c) {}
    }


}

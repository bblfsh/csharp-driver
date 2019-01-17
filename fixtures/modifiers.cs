abstract class Foo
{
    abstract public int area();
}

sealed class Bar : Foo
{
    sealed override public int area() { return 1; }
}

class ReadOnly
{
    readonly int year;
    public virtual string Name {get; set;}
    ReadOnly(int year) { this.year = year; }
}

[DllImport("avifil32.dll")]
static extern void AVIFileInit();
extern alias Foo2;

private fixed char name[30];

class Program {
    volatile int vol;

    static void Main(string[] args) {
        const int a = 0;

    }

    unsafe private Danger(byte* ps) {
        unsafe
        {
            fixed (int *p) {}
        }
    }
}

struct MyStruct 
{
    public void a() {}
}

class Program {
    static void unsafe Main(string[] args) {
        int x = 10;
        int* ptr1 = &x; 
        int y = *ptr1;
        MyStruct st = new MyStruct();
        MyStruct* ptrst = &st;
        ptrstr->a();
    }
}

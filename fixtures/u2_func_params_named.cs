class Testcls1 {
    void testfnc1(int a, string b) {}
    void testfnc2() {
        this.testfnc1(a: 1, b: "foo");
        this.testfnc1(1, b: "foo");
        this.testfnc1(a: 1, "foo");
    }
}

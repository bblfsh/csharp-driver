class Program {
    static void Main(string[] args) {
        int[] numbers = { 5, 4, 1, 3, 9, 8, 6, 7, 2, 0 };
        var lowNums = from num in numbers
            let doub = num * 2
            join num in numbers on num < 5 into testGroup
            where num < 5
            group num by num into foo
            orderby num ascending
            orderby num descending
            select num
    }
}

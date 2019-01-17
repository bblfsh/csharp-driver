class Base { }
class Test<T, U, V>
    where U : struct
    where T : Base, new()
    where V : class
{ }

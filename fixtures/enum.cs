enum Day {Sat, Sun, Mon, Tue, Wed, Thu, Fri};
enum Day2 {Sat=1, Sun, Mon, Tue, Wed, Thu, Fri};
enum Day3 : byte {Sat, Sun, Mon, Tue, Wed, Thu, Fri};

[Flags]
public enum CarOptions
{
    SunRoof = 0x01,
    Spoiler = 0x02,
    FogLights = 0x04,
    TintedWindows = 0x08,
}

class Program {
    static void Main(string[] args) {
        int x = (int)Day.Wed;
        CarOptions options = CarOptions.SunRoof | CarOptions.ForLights;
    }
}

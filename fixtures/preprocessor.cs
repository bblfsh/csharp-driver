#if DEBUG
#define BAR

#elif FOO
#undef BAR

#else
#line
#warning something
#endif

#error something

#region someRegion
const int a = 1;
#endregion

#pragma warning warn1
#pragma checksum "file.cs" "{406EA660-64CF-4C82-B6F0-42D48172A799}" "ab007f1d23d9"

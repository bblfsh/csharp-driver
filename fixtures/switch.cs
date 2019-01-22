using System;

public class Example
{
   public static void Main()
   {
      int caseSwitch = 1;
      
      switch (caseSwitch)
      {
          case 1:
              Console.WriteLine("Case 1");
              break;
          case 2:
              Console.WriteLine("Case 2");
              break;
          case 3:
              goto 1;
              break;
          case 4:
              goto default;
              break;
          case 5 when caseSwitch < 100:
              break;
          default:
              Console.WriteLine("Default case");
              break;
      }
   }
}

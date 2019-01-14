using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
 
namespace NQueens
{
    class Program
    {
        const int N = 8;
 
        static bool Allowed(bool[,] board, int x, int y)
        {
            for (int i=0; i<=x; i++)
            {
                if (board[i,y] || (i <= y && board[x-i,y-i]) || (y+i < N && board[x-i,y+i]))
                {
                    return false;
                }
            }
            return true;
        }
 
        static bool FindSolution(bool[,] board, int x)
        {
            for (int y = 0; y < N; y++)
            {
                if (Allowed(board, x, y))
                {
                    board[x, y] = true;
                    if (x == N-1 || FindSolution(board, x + 1))
                    {
                        return true;
                    }
                    board[x, y] = false;
                }
            }
            return false;
        }
 
        static void Main(string[] args)
        {
            bool[,] board = new bool[N, N];
 
            if (FindSolution(board, 0))
            {
                for (int i = 0; i < N; i++)
                {
                    for (int j = 0; j < N; j++)
                    {
                        Console.Write(board[i, j] ? "|Q" : "| ");
                    }
                    Console.WriteLine("|");
                }
            }
            else
            {
                Console.WriteLine("No solution found for n = " + N + ".");
            }
 
            Console.ReadKey(true);
        }
    }
}

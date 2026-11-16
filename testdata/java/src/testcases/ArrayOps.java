package testcases;

public class ArrayOps {
    static int[] ints(int size) {
        return new int[size];
    }

    static String[] strings(int size) {
        return new String[size];
    }

    static int[][] matrix(int rows, int columns) {
        return new int[rows][columns];
    }

    static String[][] jagged(int rows) {
        return new String[rows][];
    }

    static int length(int[][] matrix) {
        return matrix.length;
    }
}

package testcases;

import java.io.IOException;

public class Exceptions {
    void declares() throws IOException, InterruptedException {
    }

    int catches(String value) {
        try {
            return Integer.parseInt(value);
        } catch (NumberFormatException e) {
            return -1;
        } finally {
            System.gc();
        }
    }
}

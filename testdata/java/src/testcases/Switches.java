package testcases;

public class Switches {
    static int tableSwitch(int value) {
        switch (value) {
            case 0:
                return 1;
            case 1:
                return 2;
            case 2:
                return 3;
            default:
                return -1;
        }
    }

    static int lookupSwitch(int value) {
        switch (value) {
            case 10:
                return 10;
            case 50:
                return 50;
            case 100:
                return 100;
            default:
                return -1;
        }
    }

    static int stringSwitch(String value) {
        switch (value) {
            case "alpha":
                return 1;
            case "omega":
                return 2;
            default:
                return -1;
        }
    }
}

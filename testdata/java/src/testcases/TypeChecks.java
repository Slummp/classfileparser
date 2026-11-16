package testcases;

public class TypeChecks {
    static boolean isString(Object value) {
        return value instanceof String;
    }

    static String cast(Object value) {
        return (String) value;
    }

    static void locked(Object value) {
        synchronized (value) {
            value.notify();
        }
    }

    static void fail() {
        throw new IllegalStateException();
    }
}

package testcases;

public class NestedClasses {
    public class Inner {
        Object outer() {
            return NestedClasses.this;
        }
    }

    public static class Nested {
    }

    Runnable local() {
        class Local implements Runnable {
            public void run() {
            }
        }
        return new Local();
    }

    Runnable anonymous() {
        return new Runnable() {
            public void run() {
            }
        };
    }
}

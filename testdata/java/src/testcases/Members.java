package testcases;

public class Members {
    public static final int ANSWER = 42;
    private String name;
    protected volatile long counter;
    transient boolean flag;

    public Members(String name) {
        this.name = name;
    }

    public String greet(String prefix) {
        return prefix + name;
    }

    public static int count(String... values) {
        return values.length;
    }

    private synchronized void lock() {
    }

    protected native void nativeCall();

    public final double compute(int a, long b, float c, double d, char e, byte f, short g, boolean h, Object i) {
        return a + b + c + d + e + f + g + (h ? 1 : 0) + i.hashCode();
    }
}

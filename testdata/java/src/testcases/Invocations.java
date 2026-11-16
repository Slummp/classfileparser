package testcases;

import java.util.Comparator;
import java.util.List;
import java.util.function.Supplier;

interface Greeter {
    default String greet() {
        return "hi";
    }
}

public class Invocations implements Greeter {
    private static int total;
    private int count;

    public int fields() {
        count = total;
        total = count;
        return count;
    }

    public String virtualCall() {
        return toString();
    }

    public String superCall() {
        return super.toString();
    }

    public int staticCall() {
        return Math.abs(-1);
    }

    public int interfaceCall(List<String> list) {
        return list.size();
    }

    public Comparator<String> staticInterfaceCall() {
        return Comparator.naturalOrder();
    }

    public String defaultSuperCall() {
        return Greeter.super.greet();
    }

    public Supplier<String> lambda() {
        return () -> "x";
    }
}

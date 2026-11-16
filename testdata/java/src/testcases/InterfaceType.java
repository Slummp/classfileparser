package testcases;

public interface InterfaceType {
    double UNIT = 1.0;

    double area();

    default String describe() {
        return "area";
    }

    static InterfaceType unit() {
        return null;
    }

    private double twice() {
        return area() * 2;
    }
}

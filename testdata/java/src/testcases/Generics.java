package testcases;

import java.util.ArrayList;
import java.util.List;

public class Generics<T extends Number> {
    private T value;

    public <U extends Comparable<U>> List<U> process(List<U> input, T limit) {
        List<U> copy = new ArrayList<>(input);
        return copy;
    }
}

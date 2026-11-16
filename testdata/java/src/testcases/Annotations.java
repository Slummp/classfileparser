package testcases;

import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;

@Retention(RetentionPolicy.RUNTIME)
@interface Visible {
}

@Retention(RetentionPolicy.CLASS)
@interface Invisible {
}

@Visible
@Invisible
public class Annotations {
    @Visible
    @Invisible
    private int field;

    @Visible
    @Invisible
    public void method(@Visible @Invisible int parameter) {
    }
}

package testcases;

public sealed class SealedType permits SealedType.First, SealedType.Second {
    public static final class First extends SealedType {
    }

    public static non-sealed class Second extends SealedType {
    }
}

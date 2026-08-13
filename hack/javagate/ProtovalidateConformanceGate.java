package hack.javagate;

import build.buf.protovalidate.Validator;
import build.buf.protovalidate.ValidatorFactory;
import build.buf.protovalidate.exceptions.ValidationException;
import com.google.common.reflect.ClassPath;
import com.google.protobuf.Descriptors;
import com.google.protobuf.DynamicMessage;

import java.lang.reflect.Method;
import java.util.ArrayList;
import java.util.List;

/**
 * The producer-side rule-engine conformance gate: every CEL validation rule in
 * this repository's generated catalog must COMPILE on protovalidate-java with
 * default configuration -- the exact engine the platform's control plane
 * validates documents with (and the strictest of the three consumer engines:
 * Go and the browser's ES evaluator accept everything Java's defaults do).
 *
 * <p>Why per-descriptor: protovalidate compiles rules lazily, per validated
 * message TYPE. Validating one sample document compiles only the types that
 * document instantiates -- rules on nested types nobody's fixture sets (a
 * backup-schedule sub-message, a node-class requirement) are never compiled,
 * which is exactly how a Java-incompilable rule once shipped in a release and
 * emptied every fresh local instance's chart catalog. This gate walks EVERY
 * message descriptor reachable from the generated catalog classes and
 * validates a default instance of each type, forcing every rule through the
 * Java compiler. Violations on default instances are expected and ignored;
 * only compile failures fail.
 *
 * <p>Runs as part of {@code make protos} (and therefore in the proto PR lane
 * and ahead of every release): a rule the Java engine cannot evaluate can
 * never reach a tag. Authored rules must stay within protovalidate's portable
 * standard library -- no engine-specific CEL extensions.
 */
public final class ProtovalidateConformanceGate {

    /** The generated catalog's package root (protoc java_package for this repo's protos). */
    private static final String STUB_PACKAGE_ROOT = "com.dev.planton";

    /**
     * A floor on the number of message types walked, so a classpath or
     * package-root regression can never turn the gate into a vacuous green.
     * The catalog carries thousands of message types; this trips long before
     * a real shrink could be legitimate.
     */
    private static final int MINIMUM_TYPES_WALKED = 1000;

    public static void main(String[] args) throws Exception {
        Validator validator = ValidatorFactory.newBuilder().build();

        var failures = new ArrayList<String>();
        int typesWalked = 0;
        int filesSeen = 0;

        var classPath = ClassPath.from(ProtovalidateConformanceGate.class.getClassLoader());
        for (var info : classPath.getTopLevelClassesRecursive(STUB_PACKAGE_ROOT)) {
            Descriptors.FileDescriptor file = fileDescriptorOf(info);
            if (file == null) {
                continue;
            }
            filesSeen++;
            for (var message : file.getMessageTypes()) {
                typesWalked += validateTree(message, validator, failures);
            }
        }

        if (typesWalked < MINIMUM_TYPES_WALKED) {
            System.err.printf(
                    "GATE INVALID: walked only %d message types across %d files (floor: %d) -- the classpath scan "
                            + "or package root regressed; a gate that sees nothing proves nothing.%n",
                    typesWalked, filesSeen, MINIMUM_TYPES_WALKED);
            System.exit(2);
        }

        if (!failures.isEmpty()) {
            System.err.printf("%d message type(s) carry CEL rules protovalidate-java cannot compile:%n", failures.size());
            for (var failure : failures) {
                System.err.println("  - " + failure);
            }
            System.err.println("Rewrite the rule(s) using protovalidate's portable standard library "
                    + "(plain matches()/contains()/size() forms; no engine-specific string extensions).");
            System.exit(1);
        }

        System.out.printf("protovalidate-java conformance: %d message types across %d files, every rule compiles.%n",
                typesWalked, filesSeen);
    }

    /** Answers the FileDescriptor of a generated per-file outer class ("...Proto"), else null. */
    private static Descriptors.FileDescriptor fileDescriptorOf(ClassPath.ClassInfo info) {
        if (!info.getSimpleName().endsWith("Proto")) {
            return null;
        }
        try {
            Class<?> clazz = info.load();
            Method descriptor = clazz.getMethod("getDescriptor");
            Object result = descriptor.invoke(null);
            return result instanceof Descriptors.FileDescriptor ? (Descriptors.FileDescriptor) result : null;
        } catch (ReflectiveOperationException | LinkageError e) {
            return null;
        }
    }

    /**
     * Validates a default instance of this type and every nested type,
     * recording compile failures. Returns the number of types walked.
     */
    private static int validateTree(Descriptors.Descriptor descriptor, Validator validator, List<String> failures) {
        if (descriptor.getOptions().getMapEntry()) {
            return 0; // synthetic map-entry types carry no authored rules
        }
        int walked = 1;
        try {
            validator.validate(DynamicMessage.getDefaultInstance(descriptor));
        } catch (ValidationException e) {
            failures.add(descriptor.getFullName() + ": " + e.getMessage());
        }
        for (var nested : descriptor.getNestedTypes()) {
            walked += validateTree(nested, validator, failures);
        }
        return walked;
    }

    private ProtovalidateConformanceGate() {
    }
}

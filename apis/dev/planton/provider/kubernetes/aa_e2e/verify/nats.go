package verify

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	natsclient "github.com/nats-io/nats.go"
	natsjetstream "github.com/nats-io/nats.go/jetstream"
	"github.com/pkg/errors"
)

// NatsVerifier checks a NATS system to the point a customer could build
// on it: the server StatefulSet rolled out, and THE MESSAGING PROOF on
// every lane — a real NATS client (the official nats.go SDK) connects
// through the client Service and completes a pub/sub round-trip; with
// JetStream on (the kind's default) a verifier-owned stream stores a
// published marker and a consumer reads it back (messaging that cannot
// persist a message is not persistent messaging).
//
// With auth declared, the proof authenticates AS THE FIRST DECLARED USER
// (password read from the module-generated `<name>-auth` Secret — the
// paired-credential contract) and the AUTH GATE is asserted: an
// unauthenticated connection must be rejected outright (or, with a
// no_auth_user declared, must land as that guest identity).
//
// The behavioral-durability scenario (recognized by name) additionally
// DELETES the server pod after the marker is stored, waits for a
// REPLACEMENT (a new UID), reconnects and reads the SAME marker back —
// stream data survives on the JetStream PVC.
//
// Destroy is clean by design: NATS installs no CRDs — everything leaves
// with the release (stream data on the PVC leaves with the StatefulSet's
// volume claims).
type NatsVerifier struct {
	Namespace string
	Name      string
	// JetStreamEnabled gates the stream proof (the spec default is ON).
	JetStreamEnabled bool
	// FirstUsername is the user the proof connects as ("" without auth).
	FirstUsername string
	// NoAuthUser is the declared guest identity ("" = unauthenticated
	// connections must be rejected when auth is on).
	NoAuthUser string
	// NatsBoxEnabled gates the nats-box deployment check.
	NatsBoxEnabled bool
	// DurabilityProof switches on the pod-replacement arm.
	DurabilityProof bool
}

// natsFirstUsername reads the first declared username — flat users
// first, then the first account's first user ("" without auth).
func natsFirstUsername(spec map[string]interface{}) string {
	auth, _ := spec["auth"].(map[string]interface{})
	if auth == nil {
		return ""
	}
	if users, ok := auth["users"].([]interface{}); ok && len(users) > 0 {
		if entry, ok := users[0].(map[string]interface{}); ok {
			if name, ok := entry["username"].(string); ok {
				return name
			}
		}
	}
	if accounts, ok := auth["accounts"].([]interface{}); ok && len(accounts) > 0 {
		if account, ok := accounts[0].(map[string]interface{}); ok {
			if users, ok := account["users"].([]interface{}); ok && len(users) > 0 {
				if entry, ok := users[0].(map[string]interface{}); ok {
					if name, ok := entry["username"].(string); ok {
						return name
					}
				}
			}
		}
	}
	return ""
}

// natsNoAuthUser reads spec.auth.no_auth_user ("" when not declared).
func natsNoAuthUser(spec map[string]interface{}) string {
	auth, _ := spec["auth"].(map[string]interface{})
	if auth == nil {
		return ""
	}
	for _, key := range []string{"no_auth_user", "noAuthUser"} {
		if raw, ok := auth[key].(string); ok && raw != "" {
			return raw
		}
	}
	return ""
}

// natsJetStreamEnabled reads spec.jet_stream.enabled (default true — the
// kind's persistent posture).
func natsJetStreamEnabled(spec map[string]interface{}) bool {
	for _, key := range []string{"jet_stream", "jetStream"} {
		if js, ok := spec[key].(map[string]interface{}); ok {
			if raw, ok := js["enabled"]; ok {
				if enabled, ok := raw.(bool); ok {
					return enabled
				}
			}
		}
	}
	return true
}

// natsBoxEnabled reads spec.nats_box_enabled (default true).
func natsBoxEnabled(spec map[string]interface{}) bool {
	for _, key := range []string{"nats_box_enabled", "natsBoxEnabled"} {
		if raw, ok := spec[key]; ok {
			if enabled, ok := raw.(bool); ok {
				return enabled
			}
		}
	}
	return true
}

func (v *NatsVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] nats %q in namespace %q\n", v.Name, v.Namespace)

	if err := kubectlRolloutStatus(ctx, kubeconfig, "statefulset/"+v.Name, v.Namespace, 10*time.Minute); err != nil {
		return errors.Wrap(err, "the nats StatefulSet never rolled out")
	}
	if err := KubectlResourceExists(ctx, kubeconfig, "service", v.Name, v.Namespace); err != nil {
		return errors.Wrap(err, "the client Service not found")
	}
	if v.NatsBoxEnabled {
		if err := kubectlRolloutStatus(ctx, kubeconfig, "deployment/"+v.Name+"-box", v.Namespace, 5*time.Minute); err != nil {
			return errors.Wrap(err, "the nats-box deployment never rolled out")
		}
	}

	const clientPort = "14222"
	cancel, err := startPortForward(ctx, kubeconfig, "svc/"+v.Name, v.Namespace, clientPort+":4222")
	if err != nil {
		return errors.Wrap(err, "starting port-forward to the client Service")
	}
	defer cancel()
	url := "nats://127.0.0.1:" + clientPort

	options, err := v.clientOptions(ctx, kubeconfig)
	if err != nil {
		return err
	}

	conn, err := natsConnectRetry(url, options, 3*time.Minute)
	if err != nil {
		return errors.Wrap(err, "connecting to the nats server")
	}
	defer conn.Close()
	if v.FirstUsername != "" {
		fmt.Printf("  [verify] AUTH: connected as declared user %q (credential from the %s-auth Secret)\n", v.FirstUsername, v.Name)
	}

	// THE MESSAGING PROOF — core pub/sub round-trip on every lane.
	if err := v.provePubSub(conn); err != nil {
		return err
	}

	// The auth gate: without credentials the server must reject (or,
	// with a no_auth_user, admit as the guest).
	if v.FirstUsername != "" {
		if err := v.proveAuthGate(url); err != nil {
			return err
		}
	}

	if !v.JetStreamEnabled {
		return nil
	}

	// THE PERSISTENCE PROOF — a marker stored in a file-backed stream
	// and read back (plus the pod-replacement arm on the behavioral
	// lane).
	return v.proveJetStream(ctx, kubeconfig, conn, url, options)
}

func (v *NatsVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	if err := KubectlResourceAbsent(ctx, kubeconfig, "statefulset", v.Name, v.Namespace); err != nil {
		return err
	}
	if err := KubectlResourceAbsent(ctx, kubeconfig, "service", v.Name, v.Namespace); err != nil {
		return err
	}
	fmt.Printf("  [verify] DESTROY: the nats StatefulSet and client Service gone (NATS installs no CRDs — destroy is clean by design)\n")
	return nil
}

// clientOptions builds the connection options: with auth declared, the
// first user's password is read from the module-generated auth Secret
// (one key per username).
func (v *NatsVerifier) clientOptions(ctx context.Context, kubeconfig string) ([]natsclient.Option, error) {
	options := []natsclient.Option{natsclient.Timeout(20 * time.Second)}
	if v.FirstUsername == "" {
		return options, nil
	}
	encoded, err := kubectlGetJSONPath(ctx, kubeconfig, "secret", v.Name+"-auth", v.Namespace,
		fmt.Sprintf("{.data.%s}", strings.ReplaceAll(v.FirstUsername, ".", "\\.")))
	if err != nil {
		return nil, errors.Wrapf(err, "reading user %q from the %s-auth Secret", v.FirstUsername, v.Name)
	}
	password, err := base64.StdEncoding.DecodeString(strings.TrimSpace(encoded))
	if err != nil {
		return nil, errors.Wrap(err, "decoding the auth Secret value")
	}
	return append(options, natsclient.UserInfo(v.FirstUsername, string(password))), nil
}

// natsConnectRetry dials until the budget runs out (the port-forward and
// a freshly-elected server both need a beat).
func natsConnectRetry(url string, options []natsclient.Option, budget time.Duration) (*natsclient.Conn, error) {
	deadline := time.Now().Add(budget)
	var lastErr error
	for time.Now().Before(deadline) {
		conn, err := natsclient.Connect(url, options...)
		if err == nil {
			return conn, nil
		}
		lastErr = err
		time.Sleep(5 * time.Second)
	}
	return nil, lastErr
}

// provePubSub completes a subscribe → publish → receive round-trip.
func (v *NatsVerifier) provePubSub(conn *natsclient.Conn) error {
	subject := "e2e.proof.roundtrip"
	marker := fmt.Sprintf("e2e-proof-%d", time.Now().UnixNano())

	sub, err := conn.SubscribeSync(subject)
	if err != nil {
		return errors.Wrap(err, "subscribing to the proof subject")
	}
	defer func() { _ = sub.Unsubscribe() }()

	if err := conn.Publish(subject, []byte(marker)); err != nil {
		return errors.Wrap(err, "publishing the proof message")
	}
	msg, err := sub.NextMsg(30 * time.Second)
	if err != nil {
		return errors.Wrap(err, "the proof message never arrived")
	}
	if string(msg.Data) != marker {
		return errors.Errorf("the proof message arrived corrupted: %q", string(msg.Data))
	}
	fmt.Printf("  [verify] PUBSUB: publish → subscribe round-trip OK on %q\n", subject)
	return nil
}

// proveAuthGate asserts the unauthenticated posture: with a no_auth_user
// declared the connection lands as the guest; without one the server
// must reject it.
func (v *NatsVerifier) proveAuthGate(url string) error {
	conn, err := natsclient.Connect(url, natsclient.Timeout(20*time.Second))
	if v.NoAuthUser != "" {
		if err != nil {
			return errors.Wrapf(err, "no_auth_user %q is declared — an unauthenticated connection should land as the guest", v.NoAuthUser)
		}
		conn.Close()
		fmt.Printf("  [verify] AUTH GATE: unauthenticated connection admitted as the declared no_auth_user %q\n", v.NoAuthUser)
		return nil
	}
	if err == nil {
		conn.Close()
		return errors.New("an UNAUTHENTICATED connection was accepted although users are declared and no no_auth_user exists — the auth gate is open")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "authorization") {
		return errors.Wrap(err, "the unauthenticated connection failed, but not with an authorization rejection")
	}
	fmt.Printf("  [verify] AUTH GATE: unauthenticated connection REJECTED (authorization violation) — credentials are enforced\n")
	return nil
}

// proveJetStream stores a marker in a verifier-owned file-backed stream
// and reads it back; the durability arm replaces the server pod between
// store and read.
func (v *NatsVerifier) proveJetStream(ctx context.Context, kubeconfig string, conn *natsclient.Conn, url string, options []natsclient.Option) error {
	js, err := natsjetstream.New(conn)
	if err != nil {
		return errors.Wrap(err, "building the JetStream context")
	}

	jsCtx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()

	const streamName = "E2EPROOF"
	subject := "e2e.js.marker"
	marker := fmt.Sprintf("e2e-js-proof-%d", time.Now().UnixNano())

	stream, err := js.CreateOrUpdateStream(jsCtx, natsjetstream.StreamConfig{
		Name:     streamName,
		Subjects: []string{"e2e.js.>"},
		Storage:  natsjetstream.FileStorage,
	})
	if err != nil {
		return errors.Wrap(err, "creating the proof stream (JetStream is expected ON — the kind's default posture)")
	}
	// Zero-residue duty: the proof stream leaves with the verifier (a
	// hard lane death leaves it on the PVC, which the destroy phase
	// removes with the volume claims).
	defer func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cleanupCancel()
		_ = js.DeleteStream(cleanupCtx, streamName)
	}()

	if _, err := js.Publish(jsCtx, subject, []byte(marker)); err != nil {
		return errors.Wrap(err, "publishing the marker into the stream")
	}
	fmt.Printf("  [verify] JETSTREAM: marker stored in file-backed stream %q\n", streamName)

	readMarker := func(readCtx context.Context, readStream natsjetstream.Stream) error {
		consumer, err := readStream.CreateOrUpdateConsumer(readCtx, natsjetstream.ConsumerConfig{
			DeliverPolicy: natsjetstream.DeliverAllPolicy,
		})
		if err != nil {
			return errors.Wrap(err, "creating the proof consumer")
		}
		msg, err := consumer.Next(natsjetstream.FetchMaxWait(60 * time.Second))
		if err != nil {
			return errors.Wrap(err, "the stored marker never arrived")
		}
		if string(msg.Data()) != marker {
			return errors.Errorf("the stored marker arrived corrupted: %q", string(msg.Data()))
		}
		return msg.Ack()
	}

	if !v.DurabilityProof {
		if err := readMarker(jsCtx, stream); err != nil {
			return err
		}
		fmt.Printf("  [verify] JETSTREAM: the stored marker read back through a consumer — messages persist\n")
		return nil
	}

	// ---- the durability proof: server pod replacement --------------------
	// The marker is on the PVC; the pod is not the storage.
	if err := deletePodAwaitReplacement(ctx, kubeconfig, v.Namespace,
		"app.kubernetes.io/instance="+v.Name+",app.kubernetes.io/component=nats", 10*time.Minute); err != nil {
		return errors.Wrap(err, "the nats server pod did not recover after deletion")
	}

	recovered, err := natsConnectRetry(url, options, 3*time.Minute)
	if err != nil {
		return errors.Wrap(err, "reconnecting after the pod replacement")
	}
	defer recovered.Close()

	recoveredJs, err := natsjetstream.New(recovered)
	if err != nil {
		return errors.Wrap(err, "rebuilding the JetStream context after the pod replacement")
	}
	readCtx, readCancel := context.WithTimeout(ctx, 3*time.Minute)
	defer readCancel()
	recoveredStream, err := recoveredJs.Stream(readCtx, streamName)
	if err != nil {
		return errors.Wrap(err, "the proof stream is gone after the pod replacement — stream state did not survive the PVC")
	}
	if err := readMarker(readCtx, recoveredStream); err != nil {
		return errors.Wrap(err, "reading the marker AFTER the pod replacement")
	}
	fmt.Printf("  [verify] DURABILITY: the marker survived a server pod replacement through the JetStream PVC\n")
	return nil
}

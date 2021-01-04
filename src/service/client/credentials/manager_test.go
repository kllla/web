package credentials

import (
	"fmt"
	"github.com/kllla/web/src/test"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"testing"
)

const FirestoreEmulatorHost = "FIRESTORE_EMULATOR_HOST"

func TestMain(m *testing.M) {
	// command to start firestore emulator
	cmd := exec.Command("gcloud", "beta", "emulators", "firestore", "start", "--host-port=localhost:9133")

	// this makes it killable
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// we need to capture it's output to know when it's started
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	defer stderr.Close()

	// start her up!
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	// ensure the process is killed when we're finished, even if an error occurs
	// (thanks to Brian Moran for suggestion)
	var result int
	defer func() {
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		os.Exit(result)
	}()

	// we're going to wait until it's running to start
	var wg sync.WaitGroup
	wg.Add(1)

	// by starting a separate go routine
	go func() {
		// reading it's output
		buf := make([]byte, 256, 256)
		for {
			n, err := stderr.Read(buf[:])
			if err != nil {
				// until it ends
				if err == io.EOF {
					break
				}
				log.Fatalf("reading stderr %v", err)
			}

			if n > 0 {
				d := string(buf[:n])

				// only required if we want to see the emulator output
				log.Printf("%s", d)

				// checking for the message that it's started
				if strings.Contains(d, "Dev App Server is now running") {
					wg.Done()
				}

				// and capturing the FIRESTORE_EMULATOR_HOST value to set
				pos := strings.Index(d, FirestoreEmulatorHost+"=")
				if pos > 0 {
					host := d[pos+len(FirestoreEmulatorHost)+1 : len(d)-1]
					os.Setenv(FirestoreEmulatorHost, host)
				}
			}
		}
	}()

	// wait until the running message has been received
	wg.Wait()

	// now it's running, we can run our unit tests
	result = m.Run()
}

func TestNewTestManager(t *testing.T) {
	tm := newTestManager()
	defer tm.Close()

	test.AssertNotNil(t, tm)
}

func TestNewCredentials(t *testing.T) {
	tm := TestManager
	defer tm.Close()

	creds := NewCredentials("test", "password", true)

	test.AssertNil(t, tm.CreateCredentials(creds))
	test.AssertTrue(t, tm.IsCredentialsValid(creds))
	test.AssertNil(t, tm.DeleteCredentials(creds))
}

func TestNewCredentialsFailureDuplicate(t *testing.T) {
	tm := newTestManager()
	defer tm.Close()

	creds := NewCredentials("test", "password", true)

	test.AssertNil(t, tm.CreateCredentials(creds))
	test.AssertTrue(t, tm.IsCredentialsValid(creds))
	test.AssertError(t, fmt.Sprintf("username %s is unavailable", creds.Username), tm.CreateCredentials(creds))
	test.AssertNil(t, tm.DeleteCredentials(creds))
}

func TestDeleteCredentials(t *testing.T) {
	tm := TestManager
	defer tm.Close()

	creds := NewCredentials("test", "password", true)

	test.AssertNil(t, tm.CreateCredentials(creds))
	test.AssertNil(t, tm.DeleteCredentials(creds))
}

func TestDeleteCredentialsPasswordFailure(t *testing.T) {
	tm := TestManager
	defer tm.Close()

	creds := NewCredentials("test", "password", true)
	invalidCreds := NewCredentials("test", "password1", true)
	test.AssertNil(t, tm.CreateCredentials(creds))
	test.AssertError(t, "credentials invalid failed to delete", tm.DeleteCredentials(invalidCreds))
	test.AssertNil(t, tm.DeleteCredentials(creds))
}

func TestDeleteCredentialsUsernameNotFoundFailure(t *testing.T) {
	tm := TestManager
	defer tm.Close()

	creds := NewCredentials("test", "password", true)
	invalidCreds := NewCredentials("test1", "password1", true)
	test.AssertNil(t, tm.CreateCredentials(creds))
	test.AssertError(t, "credentials invalid failed to delete", tm.DeleteCredentials(invalidCreds))
	test.AssertNil(t, tm.DeleteCredentials(creds))
}

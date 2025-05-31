package testhelper

import (
	"log/slog"
	"strings"
	"testing"
)

// ExpSlogMsg holds details used to check an expected slog message. It must
// be created by the MkExpSlogMsg func.
type ExpSlogMsg struct {
	expSlogLevel slog.Level
	expContent   []string
}

// Level returns the slog Level for the checker
func (esm ExpSlogMsg) Level() slog.Level {
	return esm.expSlogLevel
}

// splitKey returns a string used to split the slog message into the part that
// comes before the Level Attr and the part coming after.
func (esm ExpSlogMsg) splitKey() string {
	return " " + slog.LevelKey + "=" + esm.Level().String() + " "
}

// MsgShldContain returns the list of strings that the mesage should contain
// for the checker
func (esm ExpSlogMsg) MsgShldContain() []string {
	return esm.expContent
}

// MkExpSlogMsg creates a single expected slog message.
func MkExpSlogMsg(level slog.Level, s ...string) ExpSlogMsg {
	return ExpSlogMsg{
		expSlogLevel: level,
		expContent:   s,
	}
}

// ExpSlogMsgList is an ordered list of expected slog message checkers.
type ExpSlogMsgList []ExpSlogMsg

// ExpSlogMessages returns the list of expected slog message checkers.
func (esml ExpSlogMsgList) ExpSlogMessages() []ExpSlogMsg {
	return []ExpSlogMsg(esml)
}

// MkExpSlogMsgList constructs an ordered list of expected slog message
// checkers, each of which must be created with the MkExpSlogMsg func.
func MkExpSlogMsgList(em ...ExpSlogMsg) ExpSlogMsgList {
	return ExpSlogMsgList(em)
}

// TestSlogMsgList is an interface matching the slog message expectation list
// methods.
type TestSlogMsgList interface {
	ExpSlogMessages() []ExpSlogMsg
}

// TestCaseWithSlogMsgList combines the TestCase and TestSlogMsgList
// interfaces.
type TestCaseWithSlogMsgList interface {
	TestCase
	TestSlogMsgList
}

// CheckExpSlogMessages calls CheckSlogMessages using the details from the
// TestCaseWithSlogMsgList to supply the parameters (the test ID and the list
// of expected slog messages).
//
// It will return false if the slog messages do not pass the checks, true
// otherwise.
func CheckExpSlogMessages(
	t *testing.T, messages string, tcsml TestCaseWithSlogMsgList,
) bool {
	t.Helper()

	return CheckSlogMessages(t,
		tcsml.IDStr(),
		messages,
		tcsml.ExpSlogMessages()...)
}

// CheckExpSlogMessagesWithID calls CheckSlogMessages using the details from
// the TestSlogMsgList to supply the parameters (the list of expected slog
// messages). The testID is passed as a parameter.
//
// It will return false if the slog messages do not pass the checks, true
// otherwise.
func CheckExpSlogMessagesWithID(
	t *testing.T, testID string,
	messages string, tsml TestSlogMsgList,
) bool {
	t.Helper()

	return CheckSlogMessages(t,
		testID,
		messages,
		tsml.ExpSlogMessages()...)
}

// CheckSlogMessages first splits the messages into separate lines. Then for
// each line it checks that the corresponging ExpSlogMsg matches it.
//
// It will return false if the slog messages do not pass the checks, true
// otherwise.
func CheckSlogMessages(
	t *testing.T,
	testID string,
	messages string,
	expSMsgs ...ExpSlogMsg,
) bool {
	t.Helper()

	msgList := strings.Split(messages, "\n")

	if len(msgList) > 0 {
		if msgList[len(msgList)-1] != "" {
			t.Log(testID)
			t.Log("\t: the last log message did not end in a newline")
			t.Log("\t: ", messages)
			t.Errorf("\t: bad slog messages")

			return false
		}

		// remove the last (empty) log message
		msgList = msgList[:len(msgList)-1]
	}

	if DiffInt(t, testID, "number of slog messages",
		len(msgList), len(expSMsgs)) {
		return false
	}

	var errorCount int

	for i, msg := range msgList {
		esm := expSMsgs[i]

		_, afterLevel, ok := strings.Cut(msg, esm.splitKey())
		if !ok {
			t.Log(testID)
			t.Logf("\t: unexpected slog Level at message %d", i)
			t.Logf("\t: the expected slog Level was %s",
				esm.Level().String())
			t.Log("\t: message:", msg)
			t.Errorf("\t: bad slog message")

			errorCount++
		}

		if ShouldContain(t,
			testID, "slog message",
			afterLevel, esm.expContent) {
			errorCount++
		}
	}

	if errorCount > 0 {
		t.Errorf("%s: errors found: %d", testID, errorCount)
	}

	return errorCount == 0
}

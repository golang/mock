package fail

import context "fmt"

func fail(msg string) string {
	return context.Sprintf("demo %s", msg)
}

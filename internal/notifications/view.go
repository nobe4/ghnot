package notifications

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/cli/go-gh/v2/pkg/tableprinter"
	"github.com/cli/go-gh/v2/pkg/term"
	"github.com/nobe4/gh-not/internal/colors"
)

func (n Notification) ToString() string {
	return fmt.Sprintf("%s %s %s by %s: '%s' ", n.prettyType(), n.prettyState(), n.Repository.FullName, n.Author.Login, n.Subject.Title)
}

var prettyTypes = map[string]string{
	"Issue":       colors.Blue("IS"),
	"PullRequest": colors.Cyan("PR"),
}

var prettyState = map[string]string{
	"open":   colors.Green("OP"),
	"closed": colors.Red("CL"),
	"merged": colors.Magenta("MG"),
}

func (n Notification) prettyType() string {
	if p, ok := prettyTypes[n.Subject.Type]; ok {
		return p
	}

	return colors.Yellow("T?")
}

func (n Notification) prettyState() string {
	if p, ok := prettyState[n.Subject.State]; ok {
		return p
	}

	return colors.Yellow("S?")
}

func (n Notifications) ToString() string {
	out := ""
	for _, n := range n {
		out += n.ToString() + "\n"
	}
	return out
}

func (n Notifications) ToTable() (string, error) {
	out := bytes.Buffer{}

	t := term.FromEnv()
	w, _, err := t.Size()
	if err != nil {
		return "", err
	}

	printer := tableprinter.New(&out, t.IsTerminalOutput(), w)

	for _, n := range n {
		printer.AddField(n.prettyType())
		printer.AddField(n.prettyState())
		printer.AddField(n.Repository.FullName)
		printer.AddField(n.Author.Login)
		printer.AddField(n.Subject.Title)
		printer.EndRow()
	}

	if err := printer.Render(); err != nil {
		return "", err
	}

	return strings.TrimRight(out.String(), "\n"), nil
}
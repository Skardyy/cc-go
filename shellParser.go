package main

import "strings"

type Alias struct {
	Name   string
	Target string
}

type ShellConfigParser struct {
	shellConfigPath    string
	partitionedContent []string
	reproduceContent   []string
	ShellParser
}

type ShellParser interface {
	Add(content []string, alias Alias) []string
	GetPartitionDel() string
}

func (scp *ShellConfigParser) With(shellConfigPath string, shellParser ShellParser) {
	scp.shellConfigPath = shellConfigPath
	scp.ShellParser = shellParser
}
func (scp *ShellConfigParser) RemoveAll() {
	scp.partitionedContent = nil
}
func (scp *ShellConfigParser) Add(aliases ...Alias) {
	scp.ReproduceAlias(aliases...)
	for _, a := range aliases {
		scp.partitionedContent = scp.ShellParser.Add(scp.partitionedContent, a)
	}
}
func (scp *ShellConfigParser) ReproduceAlias(aliases ...Alias) {
	for _, a := range aliases {
		var target string
		if strings.Contains(a.Target, " ") {
			target = "\"" + a.Target + "\""
		} else {
			target = a.Target
		}
		content := "sam add -alias " + a.Name + " " + target
		scp.reproduceContent = append(scp.reproduceContent, content)
	}
}
func (scp *ShellConfigParser) ReproducePath(aliases ...Alias) {
	for _, a := range aliases {
		target := "\"" + a.Target + "\""
		content := "sam add -path " + a.Name + " " + target
		scp.reproduceContent = append(scp.reproduceContent, content)
	}
}
func (scp *ShellConfigParser) confirm() error {
	file, err := getFile(scp.shellConfigPath)
	if err != nil {
		return err
	}
	err = replaceFilePartition(scp.ShellParser.GetPartitionDel(), file, false, scp.partitionedContent...)
	if err != nil {
		return err
	}
	return nil
}

func getDynShellParser(cf configFile) ShellConfigParser {
	//---------- here give other code different parsers ----------
	// can switch between goos.os for different ones
	parser := ShellConfigParser{}
	parser.With(cf.ShellConfigPath, &PwshConfigParsser{})
	return parser
}

// -------------------- create new config parsers to support more shells --------------------

// ---------- PowerShell file parser ----------
type PwshConfigParsser struct {
}

func (psp *PwshConfigParsser) Add(content []string, alias Alias) []string {
	pipes := strings.Split(alias.Target, "|")
	pipes[0] += " $Arguments "
	target := strings.Join(pipes, "|")

	newValue := "\nfunction " + alias.Name + " { param($Arguments) " + target + " } "
	content = append(content, newValue)
	return content
}
func (psp *PwshConfigParsser) GetPartitionDel() string {
	return "#SAM"
}

// ---------- PowerShell file parser ----------

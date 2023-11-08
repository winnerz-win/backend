package zzp

import "log"

type command struct {
	Text     string    `json:"text"`
	TagParis TagPair   `json:"tag"`
	ArgTypes []ArgType `json:"types"`
}

func (my command) isValidBody(args []string) error {
	body_cnt := len(args)
	for i, tag := range my.ArgTypes {
		switch tag {
		case argEmpty:
		case ArgJson, ArgJsonArray, ArgText:
			if i >= body_cnt {
				return Error("json param count mismatch :", args)
			}
		case ArgTextComma:
		} //
	}
	return nil
}

func Cmd(cmd string, tag Tag, args ...ArgType) *command {
	opt := &command{
		Text:     cmd,
		TagParis: tagMap[tag],
	}
	if len(args) > 0 {
		opt.ArgTypes = append(opt.ArgTypes, args...)
	} else {
		opt.ArgTypes = append(opt.ArgTypes, argEmpty)
	}

	return opt
}

func Brunch(params ...*command) []*command {
	list := []*command{}
	list = append(list, params...)
	return list
}

type nodeList []*node
type node struct {
	Cmd *command `json:"cmd,omitempty"`

	Brunch []*command `json:"brunch,omitempty"`
	Fork   nodeList   `json:"fork,omitempty"`

	IsSingleTextCheck bool `json:"is_single_text_check"`
}

func (my node) String() string { return toJsonString(my) }

// func (my *nodeData) AddFork(ns ...*nodeData) *nodeData {
// 	my.Fork = append(my.Fork, ns...)
// 	return my
// }

func newNodeData(cmd *command, brunch []*command, nexts nodeList, isSingleText bool) *node {
	my := &node{
		Cmd:               cmd,
		Brunch:            brunch,
		Fork:              nexts,
		IsSingleTextCheck: isSingleText,
	}
	emptyCmdCount := 0
	for _, v := range my.Fork {
		if v.Cmd == nil {
			emptyCmdCount++
		}

	} //for
	if emptyCmdCount > 1 {
		log.Println("[empty node] must be followed by only one. (empty-next count :", emptyCmdCount, ")")
	}
	return my
}

func Node(cmd_brunch_next ...interface{}) *node {
	var opt *command
	var brunch []*command
	var nexts nodeList
	isSingleTextCheck := false
	for _, v := range cmd_brunch_next {
		switch instance := v.(type) {
		case *command:
			opt = instance

		case []*command:
			brunch = instance

		case *node:
			nexts = append(nexts, instance)

		case bool:
			isSingleTextCheck = instance
		}
	} //for
	return newNodeData(opt, brunch, nexts, isSingleTextCheck)
}

type NodeFunc func(cmd_brunch_next ...interface{}) *node
type CmdFunc func(cmd string, tag Tag, args ...ArgType) *command
type BrunchFunc func(params ...*command) []*command

func NodeMaker() (NodeFunc, CmdFunc, BrunchFunc) {
	return Node, Cmd, Brunch
}

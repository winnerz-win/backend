package zzp

import (
	"strings"
)

func get_front_dot(text string) (string, int) {
	idx := strings.Index(text, ".")
	if idx == -1 {
		return "", idx
	}
	return strings.TrimSpace(text[:idx]), idx
}
func cut_front_dot(text string, idx int) string {
	return strings.TrimSpace(text[idx+1:])
}

func parseOptBody(text string, opt *command) ([]string, string, error) { //body, tail

	list := []rune(text)

	tag := opt.TagParis
	stack := fStack{}
	stack.Push(tag.Front())

	list = list[1:]

	bodys := []string{}
	tail := ""

	ignore := false
	idx := 0
	for _, bt := range opt.ArgTypes {
		if idx >= len(list) || stack.IsEmpty() {
			break
		}

		bodyText := ""
		for {
			elem := string(list[idx])
			switch elem {
			case "{", "[":
				switch bt {
				case ArgJson, ArgJsonArray:
					stack.Push(elem)
				} //switch
				bodyText += elem

			case ")", "}", "]", ">":
				if ignore {
					bodyText += elem
				} else {
					isLast := false
					tagpair := tagPair(elem)
					if stack.Cmp(tagpair) {
						if stack.IsFirstOne(tagpair) {
							isLast = true
						}
						stack.Pop()
					}

					if stack.Count() == 0 || isLast {
						if bt != argEmpty {
							bodys = append(bodys, strings.TrimSpace(bodyText))
						}
						bodyText = ""
						break
					}
					bodyText += elem
				}
			case ".":
				if ignore {
					bodyText += elem
				} else {
					if bt == ArgTextDot {
						bodys = append(bodys, strings.TrimSpace(bodyText))
						bodyText = ""
					} else {
						bodyText += elem
					}
				}

			case ",":
				if ignore {
					bodyText += elem
				} else {
					switch bt {
					case ArgText, ArgJsonArray:
						bodyText += elem

					case ArgJson:
						if stack.ContainCount("{", "[") == 0 {
							bodys = append(bodys, strings.TrimSpace(bodyText))
							bodyText = ""
						} else {
							bodyText += elem
						}

					default:
						if bt != argEmpty {
							bodys = append(bodys, strings.TrimSpace(bodyText))
						}
						bodyText = ""
						break
					} //switch
				}

			case "\"":
				ignore = !ignore
				bodyText += elem
			default:
				bodyText += elem
			} //switch

			idx++
			if idx >= len(list) || stack.IsEmpty() {
				break
			}

		} //for

	} //for

	if !stack.IsEmpty() {
		return nil, "", Error("cmd body cut fail")
	}

	tail = strings.TrimSpace(text[idx+1:])
	//cc.CyanItalic("tail [", tail, "]")

	return bodys, tail, nil
}

func compair_cmd(text string, opt *command) (int, error) {
	if opt == nil {
		return 0, Error("opt is nil")
	}
	if text == "" {
		return 0, Error("cmd is empty")
	}
	if len(text) < len(opt.Text)+2 {
		return 0, Error("cmd is mismatch :", opt.Text)
	}

	ti := 0
	for _, v := range opt.Text {
		a, b := string(v), text[ti:ti+1]
		if a != b {
			return 0, Error("cmd is not same")
		}
		ti++
	}
	if ti != len(opt.Text) {
		return 0, Error("cmd is not same")
	}
	cut := text[ti:]
	for _, v := range cut {
		first_word := string(v)
		switch first_word {
		case " ":
		default:
			if opt.TagParis.Front() == first_word {
				return ti, nil
			} else {
				return 0, Error("cmd is not found first_tag :", opt.TagParis)
			}
		}
	} //for
	return 0, Error("cmd is not found first_tag :", opt.TagParis)
}
func cut_tag(text string, opt *command) ([]string, string, error) { //optBody , tail , err
	ti, err := compair_cmd(text, opt)
	if err != nil {
		return nil, "", err
	}

	text = strings.TrimSpace(text[ti:])
	tag := opt.TagParis
	if tag.Front() != string(text[0]) {
		return nil, "", Error("cmd is mismatch :", opt.Text)
	}

	//cc.YellowItalic(text)

	return parseOptBody(text, opt)
}

func cut_brunch(r *normal, text string, brunch []*command) string {
	cut_func := func() bool {
		for _, opt := range brunch {
			opt_bodys, tail, err := cut_tag(text, opt)
			if err != nil {
				continue
			}

			text = cut_tail_dot(tail)
			r.brunch[opt.Text] = opt_bodys
			return true
		} //for
		return false
	}
	for {
		if !cut_func() {
			break
		}
	}
	return text
}
func cut_tail_dot(text string) string {
	if len(text) > 0 {
		if text[0] == '.' {
			text = text[1:]
		}
	}
	return text
}

func (my *node) Parse(text string) Normalizer {

	result := newNormalizer()
	r := result.normal
	cur := my
	depth := 0
	for {
		if cur == nil {
			break
		}

		if cur.Cmd == nil {
			msg, idx := get_front_dot(text)
			if idx > 0 {
				r.cmd = msg
				text = cut_front_dot(text, idx)
			} else if my.IsSingleTextCheck {
				if depth == 0 {
					r.cmd = strings.TrimSpace(text)
					text = ""
					return result
				}

			}
		} else {
			opt_bodys, tail, err := cut_tag(text, cur.Cmd)
			if err != nil {
				result.err = err
				return result
			}
			if err := cur.Cmd.isValidBody(opt_bodys); err != nil {
				result.err = err
				return result
			}
			r.cmd = cur.Cmd.Text
			r.params = opt_bodys
			text = cut_tail_dot(tail)
		}

		if len(cur.Brunch) > 0 {
			text = cut_brunch(r, text, cur.Brunch)
		}

		if len(cur.Fork) == 0 {
			break
		}

		nr := newNormal()
		r.next = nr
		r = nr
		if len(cur.Fork) == 1 {
			cur = cur.Fork[0]
		} else {
			isNext := false
			var nextNil *node
			for i, next := range cur.Fork {
				if next == nil {
					continue
				}
				if next.Cmd == nil {
					nextNil = next
				}
				_, err := compair_cmd(text, next.Cmd)
				if err == nil {
					cur = cur.Fork[i]
					isNext = true
					break
				}
			} //for
			if !isNext {
				if nextNil != nil {
					cur = nextNil
				} else {
					result.err = Error("next is mismatch :", text)
					return result
				}
			}
		}
		depth++
	} //for
	return result
}

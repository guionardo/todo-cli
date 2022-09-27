package internal

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

/* MarkDown
- [ ] 001 #tag (2020-01-01) Title {2020-01-01 00:00:00}
*/

func (item *ToDoItem) ToMarkDown() string {
	completed := " "
	if item.Completed {
		completed = "x"
	}
	tag := ""
	if len(item.Tags) > 0 {
		tags := make([]string, len(item.Tags))
		for i, tag := range item.Tags {
			tags[i] = fmt.Sprintf("#%s", tag)

		}
		tag = strings.Join(tags, " ") + " "
	}
	dueTo := ""
	if !item.DueTo.IsZero() {
		dueTo = fmt.Sprintf("(%s) ", item.DueTo.Format(TimeFormat))
	}
	lastAction := ""
	if !item.LastAction.IsZero() {
		lastAction = fmt.Sprintf(" {%s}", item.LastAction.Format("2006-01-02 15:04:05"))
	}
	return fmt.Sprintf("- [%s] %03d %s%s%s%s", completed, item.Index, tag, dueTo, item.Title, lastAction)
}

func extractCompleted(line string) (bool, string) {
	completed := !strings.HasPrefix(line, "- [ ]")
	return completed, strings.Trim(line[5:], " ")
}

func extractIndex(line string) (int, string) {
	r, _ := regexp.Compile("[0-9]{3}")
	indexStr := r.FindString(line)
	index, err := strconv.Atoi(indexStr[2:])
	if err != nil {
		return 0, line
	}
	line = strings.Replace(line, indexStr, "", 1)
	return index, line
}

func extractTags(line string) ([]string, string) {
	words := strings.Split(line, " ")
	tags := make([]string, 0)
	for i, word := range words {
		if strings.HasPrefix(word, "#") {
			tags = append(tags, strings.Trim(word[1:], " "))
			words[i] = ""
		}
	}
	newWords := make([]string, 0)
	for _, word := range words {
		if len(word) > 0 {
			newWords = append(newWords, word)
		}
	}

	return tags, strings.Join(newWords, " ")
}

func extractDueTo(line string) (time.Time, string) {
	r, _ := regexp.Compile("\\([0-9]{4}-[0-9]{2}-[0-9]{2}\\)")
	dueToStr := r.FindString(line)
	if len(dueToStr) == 0 {
		return time.Time{}, line
	}
	dueTo, err := time.Parse(TimeFormat, dueToStr[1:len(dueToStr)-1])
	if err != nil {
		return time.Time{}, line
	}
	line = strings.Replace(line, dueToStr, "", 1)
	return dueTo, strings.Trim(line, " ")
}

func extractLastAction(line string) (time.Time, string) {
	r, _ := regexp.Compile("\\{[0-9]{4}-[0-9]{2}-[0-9]{2}\\}")
	dueToStr := r.FindString(line)
	if len(dueToStr) == 0 {
		return time.Time{}, line
	}
	dueTo, err := time.Parse(TimeFormat, dueToStr[1:len(dueToStr)-1])
	if err != nil {
		return time.Time{}, line
	}
	line = strings.Replace(line, dueToStr, "", 1)
	return dueTo, strings.Trim(line, " ")
}

func (item *ToDoItem) FromMarkDown(line string) error {
	if !strings.HasPrefix(line, "- [") {
		return fmt.Errorf("Invalid line: %s", line)
	}

	completed, line := extractCompleted(line)
	index, line := extractIndex(line)
	tags, line := extractTags(line)
	dueTo, line := extractDueTo(line)
	lastAction, line := extractLastAction(line)
	item.Index = index
	item.Completed = completed
	item.Tags = tags
	item.DueTo = dueTo
	item.Title = line
	item.LastAction = lastAction
	return nil

}

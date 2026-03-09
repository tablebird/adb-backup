package shell

import (
	"adb-backup/internal/utils"
	"errors"
	"iter"
	"regexp"
	"strconv"
	"strings"
)

const (
	CONTENT_QUERY_URI_SMS = "content://sms"
)

type ContentQuery struct {
	Uri string

	Where      string
	Projection string
	Sort       string
	Limit      int
	Offset     int
	Distinct   bool
}

func (c *ContentQuery) Query(s Shell) (string, error) {
	if c.Uri == "" {
		return "", errors.New("uri is empty")
	}
	commandStr := "content query --uri " + c.Uri
	if c.Where != "" {
		commandStr += " --where '" + c.Where + "'"
	}
	if c.Projection != "" {
		commandStr += " --projection '" + c.Projection + "'"
	}
	if c.Offset != 0 {
		commandStr += " --offset " + strconv.Itoa(c.Offset)
	}
	if c.Sort != "" {
		commandStr += " --sort '" + c.Sort + "'"
	}
	if c.Limit != 0 {
		commandStr += " --limit " + strconv.Itoa(c.Limit)
	}
	if c.Distinct {
		commandStr += " --distinct"
	}

	command, err := s.RunCommand(commandStr)
	if err != nil {
		return command, err
	}
	if strings.HasPrefix(command, "No result found.") {
		return "", nil
	} else if strings.HasPrefix(command, "Row: 0") {
		return command, nil
	} else {
		return command, errors.New(command)
	}
}

func (c *ContentQuery) QueryRow(s Shell) (iter.Seq2[int, string], error) {
	result, err := c.Query(s)
	if err != nil {
		return nil, err
	}
	return parseQueryResult(result), nil
}

func (c *ContentQuery) QueryRowMap(s Shell) (iter.Seq2[int, map[string]string], error) {
	result, err := c.QueryRow(s)
	if err != nil {
		return nil, err
	}
	return func(yield func(int, map[string]string) bool) {
		for i, item := range result {
			fields := ContentQueryParseItem(item)
			if !yield(i, fields) {
				return
			}
		}
	}, nil
}

func parseQueryResult(queryResult string) iter.Seq2[int, string] {
	return func(yield func(int, string) bool) {
		rowRegex := regexp.MustCompile(`^Row:\s*(\d+)\s*(.*)`)
		lines := strings.Split(queryResult, "\n")
		currentLineNum := -1
		currentContent := ""
		for _, line := range lines {
			// Row: 0 _id=6000, thread_id=1871, address=106873402503251103, person=0, date=1752721713333, date_sent=1752721709000, protocol=0, read=1, status=-1, type=1, reply_path_present=0, subject=NULL, body=【小象超市】, service_center=+8613010113500, locked=0, sub_id=0, network_type=13, error_code=0, creator=com.android.mms, seen=1, si_id=NULL, mid=NULL, created=NULL, mtype=0, hw_is_satellite=0, privacy_mode=0, group_id=6391, addr_body=0,, time_body=0,, risk_url_body=0,, is_secret=0, resent_im=0
			// 匹配 Row: 数字 开头的行
			matches := rowRegex.FindStringSubmatch(line)
			if len(matches) == 3 {
				// 提取行号并转换为整数
				numStr := matches[1]
				contentPart := matches[2]
				num, err := strconv.Atoi(numStr)
				if err != nil {
					// 行号不是数字，归属到当前数据块
					currentContent += "\n" + line
					continue
				}

				// 验证行号是否严格递增（当前有效行号+1）
				if num == currentLineNum+1 {
					// 是新的有效行号：保存上一个数据块，初始化新块
					if currentLineNum != -1 {
						if !yield(currentLineNum, currentContent) {
							return
						}
					}
					currentLineNum = num
					currentContent = contentPart // 初始化新块内容
				} else {
					// 行号不递增，归属到当前数据块
					currentContent += "\n" + line
				}
			} else {
				// 非 Row: 数字 开头的行，归属到当前数据块
				if currentLineNum != -1 { // 仅当已有有效行号时追加
					currentContent += "\n" + line
				}
			}
		}
		if currentLineNum != -1 && currentContent != "" {
			if !yield(currentLineNum, currentContent) {
				return
			}
		}
	}
}

func ContentQueryParseItem(text string) map[string]string {
	fields := make(map[string]string)
	// 匹配单个字段（如 _id=1, address=10086）
	fieldRegex := regexp.MustCompile(`^ ?\w+=[\s\S]+`)
	currentKey := ""
	currentValue := ""
	text = utils.CleanString(text)
	split := strings.Split(text, ",")
	for i, item := range split {
		match := fieldRegex.MatchString(item)
		if match {
			data := strings.SplitN(item, "=", 2)
			if i > 0 {
				fields[currentKey] = currentValue
			}
			currentKey = strings.TrimSpace(data[0])
			currentValue = data[1]
		} else {
			currentValue += "," + item
		}
	}
	if currentKey != "" {
		fields[currentKey] = currentValue
	}
	return fields
}

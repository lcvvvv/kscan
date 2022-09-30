package appfinger

import (
	"errors"
	"fmt"
	"github.com/lcvvvv/appfinger/httpfinger"
	"io"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var database []*fingerPrint

func add(productName string, expression string) error {
	httpFinger, err := parseFingerPrint(productName, expression)
	if err != nil {
		return err
	}
	database = append(database, httpFinger)
	return nil
}

func InitDatabase(path string) (n int, lastErr error) {
	fs, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	return InitDatabaseFS(fs)
}

func InitDatabaseFS(fs io.Reader) (n int, lastErr error) {
	sourceBuf, err := io.ReadAll(fs)
	if err != nil {
		return 0, err
	}
	source := strings.Split(string(sourceBuf), "\n")
	for _, line := range source {
		line = strings.TrimSpace(line)
		r := strings.SplitAfterN(line, "\t", 2)
		if len(r) != 2 {
			lastErr = errors.New(line + "invalid")
			continue
		}
		err := add(r[0], r[1])
		if err != nil {
			lastErr = err
		}
	}
	return len(source), lastErr
}

func Clear() {
	database = []*fingerPrint{}
}

func search(banner *httpfinger.Banner) []string {
	var products []string
	for _, fingerPrint := range database {
		if productName := fingerPrint.Match(banner); productName != "" {
			products = append(products, productName)
		}
	}

	return removeDuplicate(products)
}

type Expression struct {
	//表达式数组
	paramSlice []*Param
	//表达式原文
	value string
	//表达式逻辑字符串
	//value = (body="test" || header="tt") && response="aaaa"
	//expr  = (${1} || ${2}) && ${3}
	expr string
}

func parseExpression(expr string) (*Expression, error) {
	e := &Expression{}
	e.value = expr
	//去除表达式尾部空格
	expr = strings.TrimSpace(expr)
	//修饰expr
	expr = strings.ReplaceAll(expr, `\"`, `[quota]`)
	//字符合法性校验
	if err := exprCharVerification(expr); err != nil {
		return nil, err
	}
	//提取param数组
	var paramSlice []*Param
	paramRawSlice := paramRegx.FindAllStringSubmatch(expr, -1)
	//对param进行解析
	for index, value := range paramRawSlice {
		expr = strings.Replace(expr, value[0], "${"+strconv.Itoa(index+1)+"}", 1)
		param, err := parseParam(value[0])
		if err != nil {
			return nil, err
		}
		paramSlice = append(paramSlice, param)
	}
	//语义合法性校验
	if err := exprSyntaxVerification(expr); err != nil {
		return nil, err
	}
	e.expr = expr
	e.paramSlice = paramSlice
	return e, nil
}

func (e *Expression) match(banner *httpfinger.Banner) bool {
	expr := e.makeBoolExpression(banner)
	b, _ := parseBoolFromString(expr)
	return b
}

func (e *Expression) makeBoolExpression(banner *httpfinger.Banner) string {
	var expr = e.expr
	for index, param := range e.paramSlice {
		b := param.match(banner)
		expr = strings.Replace(expr, "${"+strconv.Itoa(index+1)+"}", strconv.FormatBool(b), 1)
	}
	return expr
}

func (e *Expression) Split() []string {
	r := recursiveSplitExpression(e.expr)
	for i, v := range r {
		r[i] = e.Reduction(v)
	}
	return r
}

func (e *Expression) Reduction(s string) string {
	for i, v := range e.paramSlice {
		param := fmt.Sprintf("${%d}", i+1)
		s = strings.ReplaceAll(s, param, v.String())
	}
	return s
}

func exprCharVerification(expr string) error {
	//把所有param替换为空
	str := paramRegx.ReplaceAllString(expr, "")
	//把所有逻辑字符替换为空
	str = regexp.MustCompile(`[&| ()]`).ReplaceAllString(str, "")
	//检测是否存在其他字符
	if str != "" {
		str = strings.ReplaceAll(str, `[quota]`, `\"`)
		return errors.New(strconv.Quote(str) + " is unknown")
	}
	//检测语法合法性
	return nil
}

var regxSyntaxVerification = regexp.MustCompile(`\${\d+}`)

func exprSyntaxVerification(expr string) error {
	expr = regxSyntaxVerification.ReplaceAllString(expr, "true")
	_, err := parseBoolFromString(expr)
	if err != nil {
		return errors.New(expr + ":" + err.Error())
	}
	return nil
}

type Param struct {
	keyword  string
	value    string
	operator Operator
}
type Operator string

const (
	unequal    Operator = "!=" // !=
	equal               = "="  // =
	regxEqual           = "~=" // ~=
	superEqual          = "==" // ==
)

var keywordSlice = []string{
	"Title",
	"Header",
	"Body",
	"Response",
	"Protocol",
	"Cert",
	"Port",
	"Hash",
	"Icon",
}

var paramRegx = regexp.MustCompile(`([a-zA-Z0-9]+) *(!=|=|~=|==) *"([^"\n]+)"`)
var keywordRegx = regexp.MustCompile("^" + strings.Join(keywordSlice, "|") + "$")

func parseParam(expr string) (*Param, error) {
	p := paramRegx.FindStringSubmatch(expr)

	keyword := p[1]
	valueRaw := p[3]

	keyword = strings.ToUpper(keyword[:1]) + keyword[1:]

	if keywordRegx.MatchString(keyword) == false {
		return nil, errors.New(keyword + " keyword is unknown")
	}

	operator := convOperator(p[2])

	valueRaw = strings.ReplaceAll(valueRaw, `[quota]`, `\"`)
	value, err := strconv.Unquote("\"" + valueRaw + "\"")
	if err != nil {
		return nil, err
	}
	if operator == regxEqual {
		_, err = regexp.Compile(value)
		if err != nil {
			return nil, err
		}
	}

	return &Param{
		keyword:  keyword,
		value:    value,
		operator: operator,
	}, nil
}

func (p *Param) match(banner *httpfinger.Banner) bool {
	subStr := p.value
	keyword := p.keyword

	v := reflect.ValueOf(*banner)
	str := v.FieldByName(keyword).String()

	switch p.operator {
	case unequal:
		return !strings.Contains(str, subStr)
	case equal:
		return strings.Contains(str, subStr)
	case regxEqual:
		return regexp.MustCompile(subStr).MatchString(str)
	case superEqual:
		return str == subStr
	default:
		return false
	}
}

func (p *Param) String() string {
	return fmt.Sprintf("%s%s%s", p.keyword, p.operator, strconv.Quote(p.value))
}

const (
	and = iota // &&
	or         // ||
)

func convOperator(expr string) Operator {
	switch expr {
	case "!=":
		return unequal
	case "=":
		return equal
	case "~=":
		return regxEqual
	case "==":
		return superEqual
	default:
		panic(expr)
	}
}

func parseBoolFromString(expr string) (bool, error) {
	//去除空格
	expr = strings.ReplaceAll(expr, " ", "")
	//如果存在其他异常字符，则报错
	s := regexp.MustCompile(`true|false|&|\||\(|\)`).ReplaceAllString(expr, "")
	if s != "" {
		return false, errors.New(s + "is known")
	}
	return stringParse(expr)
}

func stringParse(expr string) (bool, error) {
	first := true
	operator := and
	if expr == "true" {
		return true, nil
	}
	if expr == "false" {
		return false, nil
	}

	for i := 0; i < len(expr); i++ {
		char := expr[i : i+1]
		if char == "t" {
			first = parseCoupleBool(first, true, operator)
			i += 3
		}
		if char == "f" {
			first = parseCoupleBool(first, false, operator)
			i += 4
		}
		if char == "&" {
			operator = and
			i += 1

		}
		if char == "|" {
			operator = or
			i += 1
		}
		if char == "(" {
			length, err := findCoupleBracketIndex(expr[i:])
			if err != nil {
				return false, err
			}
			next, err := stringParse(expr[i+1 : i+length])
			if err != nil {
				return false, err
			}
			first = parseCoupleBool(first, next, operator)

			i += length
		}

	}
	return first, nil
}

func parseCoupleBool(first bool, next bool, operator int) bool {
	if operator == or {
		return first || next
	}
	if operator == and {
		return first && next
	}
	return false
}

func findCoupleBracketIndex(expr string) (int, error) {
	var leftIndex []int
	var rightIndex []int

	for index, value := range expr {
		if value == '(' {
			leftIndex = append(leftIndex, index)
		}
		if value == ')' {
			rightIndex = append(rightIndex, index)
		}
	}

	if len(leftIndex) != len(rightIndex) {
		return 0, errors.New("bracket is not couple")
	}
	for i, index := range rightIndex {
		countLeft := strings.Count(expr[:index], "(")
		if countLeft == i+1 {
			return index, nil
		}

	}
	return 0, errors.New("bracket is not couple")
}

func trimParentheses(s string) string {
	if s[0:1] != "(" {
		return s
	}
	length := len(s)
	if length == 0 {
		return s
	}
	if s[length-1:length] != ")" {
		return s
	}
	index, err := findCoupleBracketIndex(s)
	if err != nil {
		return s
	}
	if index+1 == length {
		return s[1 : length-1]
	}
	return s
}

func splitExpression(s string) []string {
	var data []string
	for i := 0; i < len(s); i++ {
		char := s[i : i+1]
		if char == "(" {
			length, err := findCoupleBracketIndex(s[i:])
			if err != nil {
				panic(err)
			}
			data = append(data, s[i+1:i+length])
			i += length
		}
		if char == "&" {
			data = append(data, "&&")
			i += 1
		}
		if char == "|" {
			data = append(data, "||")
			i += 1
		}
		if char == "$" {
			var j = 0
			for j = 0; j < len(s)-i; j++ {
				index := i + j
				if s[index:index+1] == "}" {
					break
				}
			}
			data = append(data, s[i:i+j+1])
			i += j
		}
	}
	var r []string
	for i := 0; i < len(data); i++ {
		var s = data[i]
		if s == "||" {
			r = append(r, data[i+1])
			i += 1
			continue
		}
		if s == "&&" {
			if len(r) == 0 {
				panic("expression invalid")
			}
			for j := 1; j < len(r); j++ {
				r[0] = r[0] + " || " + r[j]
			}
			r[0] = r[0] + " && " + data[i+1]
			r = r[:1]
			i += 1
			continue
		}
		r = append(r, data[i])
	}
	return r
}

func recursiveSplitExpression(s string) []string {
	var result []string
	s = trimParentheses(s)
	for _, str := range splitExpression(s) {
		sr := splitExpression(s)
		if len(sr) == 1 {
			result = append(result, sr[0])
			continue
		}
		result = append(result, recursiveSplitExpression(str)...)
	}
	return result
}

type fingerPrint struct {
	//指纹适用范围
	//Protocol string
	//Path     string

	//指纹适配产品
	ProductName string

	//指纹识别规则
	Keyword *Expression
}

func parseFingerPrint(productName string, expressionString string) (*fingerPrint, error) {

	//校验合法，序列化表达式
	expression, err := parseExpression(expressionString)
	if err != nil {
		return nil, err
	}

	return &fingerPrint{
		//指纹适用范围
		//Protocol: "",
		//Path:     "",
		//指纹适配产品
		ProductName: productName,
		//指纹识别规则
		//Icon:    "",
		//Hash:    "",
		Keyword: expression,
	}, nil
}

func (f *fingerPrint) Match(banner *httpfinger.Banner) string {
	if f.Keyword != nil {
		if f.Keyword.match(banner) {
			return f.ProductName
		}
	}
	return ""
}

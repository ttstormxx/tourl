package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/dlclark/regexp2"
	"io"
	"io/ioutil"
	"log"
	// "net"
	"errors"
	"os"
	"sort"
	"strconv"
	"strings"
)

func ReadLine(filename string) ([]string, error) {

	var result []string
	// pip begins

	stat, _ := os.Stdin.Stat()

	if (stat.Mode() & os.ModeCharDevice) == 0 {

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			result = append(result, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		// pipe ends
	} else {
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer func() {
			f.Close()
			// fmt.Println("文件关闭成功")
		}()
		reader := bufio.NewReader(f)

		for {
			// 这里文本最后一行读不到，需要处理,已处理
			// 或者最后一行置空
			// line,err :=buf.ReadString('\n')
			line, _, err := reader.ReadLine()

			if err != nil {
				if err == io.EOF {
					// 读取结束，报EOF

					// fmt.Println("读取结束")
					break
				}
				return nil, err
			}
			linestr := string(line)
			result = append(result, linestr)
		}
	}
	var temp []string
	for _, value := range result {

		// 处理两头空白字符
		value = strings.TrimSpace(value)
		// 抛弃空行
		expr := `^$`
		reg, _ := regexp2.Compile(expr, 0)
		if isMatch, _ := reg.MatchString(value); !isMatch {
			temp = append(temp, value)

		}
	}
	result = temp
	return result, nil
}

func UrlToIpsWithPort(urls []string) []string {
	// 处理URL为IP:PORT列表
	// 正则表达式稍有问题
	// URL后面如果没有/则无法匹配
	// http://123.123.123.123 无法匹配
	expr := `(?<=://).+?(?=/)`
	reg, _ := regexp2.Compile(expr, 0)

	// 对于URL后方无/的，主动添加/
	// http://123.123.123.123
	// 变为
	// http://123.123.123.123/
	// 同时针对单行为IP/域名的情况，主动添加http://xxxx/
	//
	expr_http_finder := `^(http://|https://)`
	reg2, _ := regexp2.Compile(expr_http_finder, 0)

	var temp []string
	for _, value := range urls {
		// 任何value，后方均➕/
		value = value + "/"
		// 查找开头是否为http://或https://，没有则加上

		if isMatch, _ := reg2.MatchString(value); !isMatch {
			value = "http://" + value
		}

		match, _ := reg.FindStringMatch(value)

		ipPort := match.String()
		// ipPort = strings.Split(ipPort, ":")[0]
		// 处理 IP:PORT列表为IP列表
		temp = append(temp, ipPort)

	}
	return temp
}

func UrlToIpsNoPort(urls []string) []string {
	// 处理URL为IP:PORT列表
	// 正则表达式稍有问题
	// URL后面如果没有/则无法匹配
	// http://123.123.123.123 无法匹配
	expr := `(?<=://).+?(?=/)`
	reg, _ := regexp2.Compile(expr, 0)

	// 对于URL后方无/的，主动添加/
	// http://123.123.123.123
	// 变为
	// http://123.123.123.123/
	// 同时针对单行为IP/域名的情况，主动添加http://xxxx/
	//
	expr_http_finder := `^(http://|https://)`
	reg2, _ := regexp2.Compile(expr_http_finder, 0)

	var temp []string
	for _, value := range urls {
		// 任何value，后方均➕/
		value = value + "/"
		// 查找开头是否为http://或https://，没有则加上

		if isMatch, _ := reg2.MatchString(value); !isMatch {
			value = "http://" + value
		}

		match, _ := reg.FindStringMatch(value)

		ipPort := match.String()
		ipPort = strings.Split(ipPort, ":")[0]
		// 处理 IP:PORT列表为IP列表
		temp = append(temp, ipPort)

	}
	return temp
}

func GetAsignPorts(ports string) (intTempPorts []int, err error) {
	var tempPorts []string
	var dashPorts []string
	var stringPorts []string
	// 切割"," "-"
	tempPorts = strings.Split(ports, ",")
	for _, value := range tempPorts {
		if value == "0" {
			err = errors.New("0不与其他端口同时指定，请检查")
			return
		}
	}
	for _, value := range tempPorts {
		if strings.Contains(value, "-") {
			dashPorts = append(dashPorts, value)
		} else {
			isInt := IsInt(value)
			if !isInt {
				err = errors.New("输入的值：" + value + " 不是数字，请检查")
				return
			} else {
				if big, _ := strconv.Atoi(value); big > 65535 {
					err = errors.New("端口值" + value + "： 数字大于65535，请检查")
					return
				}
				stringPorts = append(stringPorts, value)
			}
		}
	}

	for _, value := range stringPorts {
		int, _ := strconv.Atoi(value)
		intTempPorts = append(intTempPorts, int)
	}

	// 扩展所有使用横杠标识的端口范围
	if len(dashPorts) != 0 {
		// fmt.Println("范围端口处理开始。。。")
		for _, value := range dashPorts {
			if dashCount := strings.Count(value, "-"); dashCount != 1 {
				err = errors.New("端口范围不规范:" + value + " 请检查")
				return
			}
			pointPorts := strings.Split(value, "-")

			expr := `^\d+$`
			reg, _ := regexp2.Compile(expr, 0)
			for _, value1 := range pointPorts {
				if isMatch, _ := reg.MatchString(value1); !isMatch {
					err = errors.New("输入的端口范围非数字")
					fmt.Println("存在问题的端口范围为：", value)
					return
				}
			}

			// fmt.Println("未发现非数字，范围端口处理继续。。。")
			i, _ := strconv.Atoi(pointPorts[0])
			j, _ := strconv.Atoi(pointPorts[1])

			if i > j {
				err = errors.New("端口范围" + value + ": 数字大小有误，请检查")
				return
			}
			if j > 65535 {
				err = errors.New("端口范围" + value + "： 数字大于65535，请检查")
				return
			}
			for n := i; n <= j; n++ {
				if n == 0 {
					continue
				} else {
					intTempPorts = append(intTempPorts, n)
				}

			}
		}
	}

	// 去重
	var finalTemp []int
	finalTemp = uniqueArr(intTempPorts)
	// 排序
	finalTemp = ascArr(finalTemp)
	intTempPorts = finalTemp

	return intTempPorts, nil
}

// 去重
func uniqueArr(m []int) []int {
	d := make([]int, 0)
	tempMap := make(map[int]bool, len(m))
	for _, v := range m { // 以值作为键名
		if tempMap[v] == false {
			tempMap[v] = true
			d = append(d, v)
		}
	}
	return d
}

// 升序
func ascArr(e []int) []int {
	sort.Ints(e[:])
	return e
}

func IsInt(str string) bool {
	var temp bool = false
	_, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("输入的不是数字，请检查")
		temp = false
	} else {
		temp = true
	}
	return temp
}
func IpWriteToFile(ips []string, save_to_file string) bool {
	var tempString string
	// var tempStringCsv string

	// expr := `.txt`
	if !strings.HasSuffix(save_to_file, ".txt") {
		save_to_file += ".txt"
	}
	for _, value := range ips {

		tempString += value + "\n"
		// tempStringCsv +=
		// 写入txt
		err := ioutil.WriteFile(save_to_file, []byte(tempString), 0666)
		if err != nil {
			fmt.Println(err)
		}
		// // 写入CSV文件
		// err1 := ioutil.WriteFile("results.csv", []byte(tempStringCsv), 0666)
		// if err1 != nil {
		// 	fmt.Println(err1)
		// }
	}
	return true
}

func main() {
	// 转化IP/域名为指定URL形式
	// 指定PORT
	// 指定多PORT
	// 指定http|https
	// 指定URI
	// 单URL/IP/域名 对 多PORT 单URI
	// 设置WEB常见端口列表
	// 指定URI  path

	// 处理URL成为IP:PORT
	// 或IP形式
	url_file_path := flag.String("l", "ips.txt", "url文件路径")

	keep_port := flag.String("p", "", "指定port")
	is_common_web_port := flag.Bool("P", false, "在常见WEB端口组的基础上指定port，可单独使用")
	is_major_big_web_port_list := flag.Bool("PP", false, "大容量WEB端口")
	to_http := flag.Bool("th", false, "转换为HTTP")
	to_https := flag.Bool("ts", false, "转换为HTTPS")
	path := flag.String("path", "/", "URI路径")
	quiet_mod := flag.Bool("q", false, "安静模式，减少输出")
	save_to_file := flag.String("o", "result.txt", "输出到指定文件")

	flag.Parse()
	args := os.Args[1:]

	var common_port_list = []int{81, 3443, 5601, 7001, 8000, 8001, 8080, 8081, 8088, 8090, 8161, 8443, 8888, 9000, 9200, 9999}
	var major_big_web_port_list = []int{81, 82, 88, 1443, 3443, 5601, 7001, 8000, 8001, 8080, 8081, 8088, 8090, 8161, 8181, 8443, 8888, 8899, 9000, 9090, 9200, 9999}

	// 启用big超级WEB端口组
	if *is_major_big_web_port_list {
		common_port_list = major_big_web_port_list
		*is_common_web_port = true
	}
	fmt.Println("P是启动的吗：", *is_common_web_port)
	fmt.Println("PP是启动的吗：", *is_major_big_web_port_list)
	fmt.Println("目前附加的端口组是： 总计：", len(common_port_list))
	fmt.Println(common_port_list)
	urls, err := ReadLine(*url_file_path)
	if err != nil {
		fmt.Println(err)
	}

	var hasp bool
	var hasP bool
	var Pnop bool

	exp_p := `^-p$`
	reg_p, _ := regexp2.Compile(exp_p, 0)
	// exp_P := `^-P$`
	exp_P := `^(-P|-PP)$`
	reg_P, _ := regexp2.Compile(exp_P, 0)
	for _, value := range args {
		if ismatch_p, _ := reg_p.MatchString(value); ismatch_p {
			hasp = true
		}
		if ismatch_P, _ := reg_P.MatchString(value); ismatch_P {
			hasP = true
		}
	}
	if !hasp && hasP {
		Pnop = true
	}

	if !*quiet_mod {
		if Pnop {
			fmt.Println("无指定端口，仅启用默认端口组")
		}
	}

	var misc []string
	if *keep_port == "" && !Pnop {
		misc = UrlToIpsWithPort(urls)
	} else if Pnop {
		misc = UrlToIpsWithPort(urls)
		ports := common_port_list
		if *to_http || *to_https {
			var temp []int
			for _, value := range ports {
				if value != 80 && value != 443 {
					temp = append(temp, value)
				}
			}
			ports = temp
		}
		if !*quiet_mod {
			fmt.Println("扩展PORT列表为：")
			fmt.Println(ports)
		}

		var temp []string
		for _, port := range ports {
			strPort := strconv.Itoa(port)
			for _, url := range misc {
				tempUrl := url + ":" + strPort
				temp = append(temp, tempUrl)
			}
		}
		misc = temp
	} else {

		// p 0  去除端口
		if *keep_port == "0" {
			misc = UrlToIpsNoPort(urls)
		} else {
			misc = UrlToIpsNoPort(urls)
			// 以空格分割port组
			ports, err := GetAsignPorts(*keep_port)
			if err != nil {
				fmt.Println("出现错误：", err)
			}

			if *is_common_web_port {
				if !*quiet_mod {
					fmt.Println("指定端口，同时启用常见WEB端口组。。。")
				}

				ports = append(ports, common_port_list...)
				// 去重
				ports = uniqueArr(ports)
				// 升序
				ports = ascArr(ports)

			}
			if *to_http || *to_https {
				var temp []int
				for _, value := range ports {
					if value != 80 && value != 443 {
						temp = append(temp, value)
					}
				}
				ports = temp
			}

			if !*quiet_mod {
				fmt.Println("扩展PORT列表为：")
				fmt.Println(ports)
			}

			var temp []string
			for _, port := range ports {
				strPort := strconv.Itoa(port)
				for _, url := range misc {
					tempUrl := url + ":" + strPort
					temp = append(temp, tempUrl)
				}
			}
			misc = temp

		}
	}

	if *path != "/" {
		if !strings.HasPrefix(*path, "/") {
			*path = "/" + *path
		}
		var tempSlice []string
		for _, value := range misc {
			temp := value + *path
			tempSlice = append(tempSlice, temp)
		}
		misc = tempSlice
	}

	// to http

	if *to_http {
		var tempSlice []string
		for _, value := range misc {
			temp := "http://" + value
			tempSlice = append(tempSlice, temp)
		}
		misc = tempSlice
	} else if *to_https {
		var tempSlice []string
		for _, value := range misc {
			temp := "https://" + value
			tempSlice = append(tempSlice, temp)
		}
		misc = tempSlice
	}

	var for_http_list []string
	origin_ips := UrlToIpsNoPort(urls)
	if *to_http && *keep_port != "0" && *is_common_web_port {
		var temp []string
		for _, value := range origin_ips {
			var tempStr string
			if *path == "/" {
				tempStr = "http://" + value
			} else {
				tempStr = "http://" + value + *path
			}

			temp = append(temp, tempStr)
		}
		for_http_list = append(for_http_list, temp...)
		for_http_list = append(for_http_list, misc...)
		misc = for_http_list
	} else if *to_https && *keep_port != "0" && *is_common_web_port {
		var temp []string
		for _, value := range origin_ips {
			var tempStr string
			if *path == "/" {
				tempStr = "https://" + value
			} else {
				tempStr = "https://" + value + *path
			}
			temp = append(temp, tempStr)
		}
		for_http_list = append(for_http_list, temp...)
		for_http_list = append(for_http_list, misc...)
		misc = for_http_list
		// fmt.Println("初始附加完毕")
	}

	// 不生成http或https时，指定P的情况下，附上去除port的原始结果
	for_http_list = []string{}
	if *path == "/" {
		for_http_list = append(for_http_list, origin_ips...)
	} else {
		for _, value := range origin_ips {
			for_http_list = append(for_http_list, value+*path)
		}
	}

	if !*to_http && !*to_https && *is_common_web_port {

		for_http_list = append(for_http_list, misc...)
		misc = for_http_list
	}

	if !*quiet_mod {
		fmt.Println("\n")
		fmt.Println("IP总计", len(urls), "个")
		fmt.Println("扩展结果总计：", len(misc), "个")
		fmt.Println("\n")
	}

	for _, value := range misc {
		fmt.Println(value)
	}

	for _, value := range args {
		if strings.Contains(value, "-o") {
			IpWriteToFile(misc, *save_to_file)

			if !*quiet_mod {
				fmt.Println("\n")
				fmt.Println("文件写入完毕。。。")
			}

		}
	}

}

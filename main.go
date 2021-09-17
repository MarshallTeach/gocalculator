package main

import (
	"bufio"
	"fmt"
	"github.com/MarshallTeach/gods"
	"math"
	"os"
	"strings"
)
/*
双栈解决 [通用表达式] 问题的通用解法

对于「任何表达式」而言，我们都使用两个栈 nums 和 ops：

nums ： 存放所有的数字
ops ：存放所有的数字以外的操作
然后从前往后做，对遍历到的字符做分情况讨论：

空格 : 跳过
( : 直接加入 ops 中，等待与之匹配的 )
) : 使用现有的 nums 和 ops 进行计算，直到遇到左边最近的一个左括号为止，计算结果放到 nums
数字 : 从当前位置开始继续往后取，将整一个连续数字整体取出，加入 nums
+ - * / ^ % : 需要将操作放入 ops 中。在放入之前先把栈内可以算的都算掉（只有「栈内运算符」比「当前运算符」优先级高/同等，才进行运算），使用现有的 nums 和 ops 进行计算，直到没有操作或者遇到左括号，计算结果放到 nums
      我们可以通过 🌰 来理解 只有「栈内运算符」比「当前运算符」优先级高/同等，才进行运算 是什么意思：

因为我们是从前往后做的，假设我们当前已经扫描到 2 + 1 了（此时栈内的操作为 + ）。

如果后面出现的 + 2 或者 - 1 的话，满足「栈内运算符」比「当前运算符」优先级高/同等，可以将 2 + 1 算掉，把结果放到 nums 中；
如果后面出现的是 * 2 或者 / 1 的话，不满足「栈内运算符」比「当前运算符」优先级高/同等，这时候不能计算 2 + 1。
一些细节：

由于第一个数可能是负数，为了减少边界判断。一个小技巧是先往 nums 添加一个 0
为防止 () 内出现的首个字符为运算符，将所有的空格去掉，并将 (- 替换为 (0-，(+ 替换为 (0+（当然也可以不进行这样的预处理，将这个处理逻辑放到循环里去做）
从理论上分析，nums 最好存放的是 long，而不是 int。因为可能存在 大数 + 大数 + 大数 + … - 大数 - 大数 的表达式导致中间结果溢出，最终答案不溢出的情况

时间复杂度：O(n)
空间复杂度：O(n)
*/

type Calculator interface {
	calculate(string)
}

type cal struct {
	calStr  string
	calChan chan string
	calPriority map[byte]int
}

func newCal() *cal {
	return &cal{
		calChan: make(chan string, 1),
		calPriority: map[byte]int{
			'-': 1,
			'+': 1,
			'*': 2,
			'/': 2,
			'%': 2,
			'^': 3,
		},
	}
}
//(1 + 2) * 3
func main() {
	scanStr:
		c := newCal()
		fmt.Println("请输入待计算的公式")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			c.calStr = scanner.Text()
			break
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("输入为空")
			goto scanStr
		}
	c.calChan<-c.calStr
	select {
		case str := <-c.calChan:
			var x Calculator
			x = c
			x.calculate(str)
			goto scanStr
	}
}

func (cl *cal) calculate(str string) {
	// (1 + 2) * 5 / 2^2
	str = strings.ReplaceAll(str, " ", "")
	// fmt.Println(strings.TrimSpace(str))
	n := len(str)
	// 存放所有数字
	nums := gods.NewDeque()
	nums.Append(0)
	// 存放所有非数字操作
	ops := gods.NewDeque()
	for i := 0; i < n; i++ {
		c := str[i]
		if c == '(' {
			ops.Append(c)
		} else if c == ')' {
			// 计算到最近一个左括号为止
			for !ops.Empty() {
				if ops.Last().(byte) != '(' {
					calc(nums, ops)
				} else {
					ops.Pop()
					break
				}
			}
		} else {
			if isNumber(string(c)) {
				u := 0
				j := i
				// 将从i位置开始后面的连续数字整体取出 加入 nums
				for j < n && isNumber(string(str[j])) {
					u = u * 10 + int(str[j] - '0')
					j++
					nums.Append(u)
					i = j - 1
				}
			} else {
				if i > 0 && (str[i - 1] == '(' || str[i - 1] == '+' || str[i - 1] == '-') {
					nums.Append(0)
				}
				// 有一个新操作要入栈时，先把栈内可以算的都算了 
				// 只有满足「栈内运算符」比「当前运算符」优先级高/同等，才进行运算
				for !ops.Empty() && ops.Last().(byte) != '(' {
					prev := ops.Last().(byte)
					fmt.Println(cl.calPriority[prev], cl.calPriority[c])
					if cl.calPriority[prev] >= cl.calPriority[c] {
						calc(nums, ops)
					} else {
						break
					}
				}
				ops.Append(c)
			}
		}
	}
	for !ops.Empty() {
		calc(nums, ops)
	}

	fmt.Println(nums.Last())
}

func calc(nums, ops *gods.Deque)  {
	if nums.Empty() || nums.Size() < 2 {
		return
	}

	if ops.Empty() {
		return
	}

	b := nums.Pop().(int)
	a := nums.Pop().(int)
	op := ops.Pop().(byte)
	ans := 0
	if op == '+' {
		ans = a + b
	} else if op == '-' {
		ans = a - b
	} else if op == '*' {
		ans = a * b
	} else if op == '/' {
		ans = a / b
	} else if op == '^' {
		ans = int(math.Pow(float64(a), float64(b)))
	} else if op == '%' {
		ans = a % b
	}
	nums.Append(ans)
}

func isNumber(str string) bool {
	if len(str) == 0 {
		return false
	}
	for _, s := range str {
		if s < '0' || s > '9' {
			return false
		}
	}
	return true
}

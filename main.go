package main

import "fmt"
import "strconv"
import "strings"
import "os"
import "bufio"
import "flag"

type big_num struct {
    chars []byte
    sign byte
}

type represent interface {
    restruct() string
    rm_head0() []byte
}

func (bn big_num) restruct() string {
    var s string
    if bn.sign == 45 {
        s += "-"
    }
    for _, ch := range bn.chars {
        s += fmt.Sprintf("%c", ch + 48)
    }
    return s
}

func (bn big_num) rm_head0() []byte {
    for bn.chars[0] == 0 && len(bn.chars) > 1 {
        bn.chars = bn.chars[1:]
    }
    return bn.chars
}

func convert(s string) big_num {
    var bn big_num

    bn.sign = 43
    if s[0] == 45 {
        bn.sign += 2
        s = s[1:]
    } else if s[0] == 43 {
        s = s[1:]
    }
    for _, c := range s {
        bn.chars = append(bn.chars, byte(c) - 48)
    }
    bn.chars = bn.rm_head0()
    return bn
}

func compare(bn1, bn2 big_num) int8 {
    sti := strconv.ParseInt
    s1, s2 := bn1.restruct(), bn2.restruct()
    _n1, _ := sti(s1, 10, 64)
    _n2, _ := sti(s2, 10, 64)
    n1, n2 := abs(int(_n1)), abs(int(_n2))
    if n1 < n2 {
        return -1
    } else if n1 == n2 {
        return 0
    }
    return 1
}

func fix_len(re_len int, x []byte) []byte {
    l := len(x)
    if l == re_len {
        return x
    }
    new_x := make([]byte, re_len)
    fix_idx := re_len - l
    for i := 0; i < l; i++ {
        new_x[i + fix_idx] = x[i]
    }
    return new_x
}

func abs(n int) int {
    if n < 0 {
        return -n
    }
    return n
}

func max(a, b int) int {
    return (a + b + abs(a-b))/2
}

func add(a, b []byte) []byte {
    la, lb := len(a), len(b)
    max_len := max(la, lb)
    a, b = fix_len(max_len, a), fix_len(max_len, b)

    var c, temp byte
    res := make([]byte, max_len)

    for i := max_len - 1; i >= 0; i-- {
        temp = a[i] + b[i] + c
        if temp >= 10 {
            res[i] = temp - 10
            c = 1
            continue
        }
        res[i] = temp
        c = 0
    }
    if c == 1 {
        return append([]byte{1}, res...)
    }
    return res
}

func sub(a, b []byte) []byte {
    la, lb := len(a), len(b)
    max_len := max(la, lb)
    a, b = fix_len(max_len, a), fix_len(max_len, b)

    var c, temp byte
    res := make([]byte, max_len)

    for i := max_len - 1; i >= 0; i-- {
        temp = b[i] + c
        if a[i] < temp {
            res[i] = 10 + a[i] - temp
            c = 1
            continue
        }
        res[i] = a[i] - temp
        c = 0
    }

    return res
}

func xyz(x, y, z byte) byte {
    return 4*(x-43) + 2*(y-43) + (z-43)
}

func process(bn1, bn2 big_num, ope byte) big_num {
    var res big_num
    k := xyz(bn1.sign, bn2.sign, ope)
    if k == 0 || k == 6 {
        res.chars = add(bn1.chars, bn2.chars)
        res.sign = 43
    } else if k == 2 || k == 4 {
        if compare(bn1, bn2) == -1 {
            res.chars = sub(bn2.chars, bn1.chars)
            res.sign = 45
        } else {
            res.chars = sub(bn1.chars, bn2.chars)
            res.sign = 43
        }
    } else if k == 8 || k == 14 {
        bn2.sign = 43
        res = process(bn2, bn1, 43)
    } else {
        bn1.sign , bn2.sign = 43, 43
        res = process(bn1, bn2, 43)
        res.sign = 45
    }
    res.chars = res.rm_head0()
    return res
}

func main() {
    path := flag.String("p", "", "path of input file")
    flag.Parse()
    if *path == "" {
        return
    }

    f, err := os.Open(*path)
    if err != nil {
        fmt.Println("input find not found")
        return
    }
    defer f.Close()

    var calated [][]string
    var temp []string
    var a, b big_num
    var ope byte
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        temp_s := scanner.Text()
        temp = append(temp, temp_s)
        data := strings.Split(temp_s, " ")
        a, b, ope = convert(data[0]), convert(data[1]), data[2][0]
        temp = append(temp, process(a, b, ope).restruct())
        calated = append(calated, temp)
        temp = []string{}
    }

    for _, item := range calated {
        fmt.Println(item)
    }
}

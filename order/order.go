package order

import (
	"fmt"
	"strconv"
	//"strings"
	"strings"
)

//订单均分
func test() {
	ids := make([]int, 3)
	ids[0] = 1
	ids[1] = 2
	ids[2] = 3
	num := 5
	ms := GetSalesman()
	fmt.Println("##########")
	Distribution1(len(ids), num, ms, ids)

}

func GetSalesman() []*SalesmanAllotment {
	ms := make([]*SalesmanAllotment, 0)
	s1 := new(SalesmanAllotment)
	s1.Id = 101
	s1.PmId = 10
	ms = append(ms, s1)

	s2 := new(SalesmanAllotment)
	s2.Id = 102
	s2.PmId = 10
	ms = append(ms, s2)

	s3 := new(SalesmanAllotment)
	s3.Id = 103
	s3.PmId = 10
	ms = append(ms, s3)

	s4 := new(SalesmanAllotment)
	s4.Id = 104
	s4.PmId = 10
	ms = append(ms, s4)

	s5 := new(SalesmanAllotment)
	s5.Id = 105
	s5.PmId = 10
	ms = append(ms, s5)

	s6 := new(SalesmanAllotment)
	s6.Id = 106
	s6.PmId = 10
	ms = append(ms, s6)
	return ms
}

type SalesmanAllotment struct {
	Id         int
	SalesmanId int //业务员ID
	PmId       int //经理ID
}

//ids 催收员人数
//num 订单数
func Distribution1(ids, num int, salemens []*SalesmanAllotment, userIds []int) {
	keepIds := make([]int, ids)         //记录每个人分配的订单下表
	ordersKeeper := make([]string, ids) //每个人得到的订单
	for i := 0; i < num; i++ {
		present := findMinMoneyIndex(keepIds)
		ordersKeeper[present] += strconv.Itoa(salemens[i].Id) + ","
		keepIds[present] += salemens[i].Id
	}
	for j := 0; j < ids; j++ {
		sp := strings.Split(ordersKeeper[j], ",")
		for n := 0; n < len(sp); n++ {
			if sp[n] != "" {
				fmt.Println("n", n, sp[n], "j", j)
				//写update，修改数据即可
				fmt.Println("用户:", userIds[j])
			}
		}
	}
	/*for j := 0; j < num; j++ {
		rstr := substring(ordersKeeper[j], 0, len(ordersKeeper[j])-1)
		ids := strings.Split(rstr, ",")
		for _, v := range ids {
			fmt.Println(v)
		}
	}*/
}

//找到最小下标
func findMinMoneyIndex(keepAmount []int) int {
	result := len(keepAmount) - 1
	if result <= 0 {
		return 0
	}
	min := keepAmount[len(keepAmount)-1]
	index := len(keepAmount) - 2
	//避免先给后面的人，反向操作
	for ; index > -1; index-- {
		if keepAmount[index] <= min {
			result = index
		}
	}
	return result
}

func substring(source string, start int, end int) string {
	var r = []rune(source)
	length := len(r)

	if start < 0 || end > length || start > end {
		return ""
	}

	if start == 0 && end == length {
		return source
	}
	return string(r[start:end])
}

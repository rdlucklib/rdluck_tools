package paging

const (
	PageSize20=20
)

type PagingReq struct {
	CurrentIndex  int  `description:"当前页页码"`
	PageSize      int  `description:"每页数据条数,如果不传，默认每页20条"`
}

type PagingItem struct {
	IsStart       bool `description:"是否首页"`
	IsEnd         bool `description:"是否最后一页"`
	PreviousIndex int  `description:"上一页页码"`
	NextIndex     int  `description:"下一页页码"`
	CurrentIndex  int  `description:"当前页页码"`
	Pages         int  `description:"总页数"`
	Totals        int  `description:"总数据量"`
	PageSize      int  `description:"每页数据条数"`
}

func GetPaging(currentIndex,pageSize,total int)(item *PagingItem)  {
	if pageSize<=0 {
		pageSize=PageSize20
	}
	if currentIndex<=0 {
		currentIndex=1
	}
	item=new(PagingItem)
	item.PageSize=pageSize
	item.Totals=total
	item.CurrentIndex=currentIndex

	if total<=0 {
		item.IsStart=true
		item.IsEnd=true
		return
	}
	pages:=PageCount(total,pageSize)
	item.Pages=pages
	if pages<=1 {
		item.IsStart=true
		item.IsEnd=true
		item.PreviousIndex=1
		item.NextIndex=1
		return
	}
	if pages == currentIndex {
		item.IsStart=false
		item.IsEnd=true
		item.PreviousIndex=currentIndex-1
		item.NextIndex=currentIndex
		return
	}
	if currentIndex==1 {
		item.IsStart=true
		item.IsEnd=false
		item.PreviousIndex=1
		item.NextIndex=currentIndex+1
		return
	}
	item.IsStart=false
	item.IsEnd=false
	item.PreviousIndex=currentIndex-1
	item.NextIndex=currentIndex+1
	return
}


func StartIndex(page, pagesize int) int {
	if page > 1 {
		return (page - 1) * pagesize
	}
	return 0
}
func PageCount(count, pagesize int) int {
	if count%pagesize > 0 {
		return count/pagesize + 1
	} else {
		return count / pagesize
	}
}
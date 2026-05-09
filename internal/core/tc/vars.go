package tc

import "github.com/kakeetopius/qosm/internal/core/filter"

//1: Root (The HTB classfull qdisc to be attached at egress)
//└── 1:1 ParentClass (parent of the all classes.)
//    ├── 1:10 HighClass (high priority class. Only Packets matched by HighPrioClassFilter reach here.)
//    ├── 1:19 LowClass  (low priority class. Only Packets matched by LowPrioClassFilter reach here.)
//    └── 1:15 DefaultClass (class where packets go if no filter matches on them)

var Root = HTBQdisc{
	Handle:       HTBQDISCHANDLE,
	Parent:       ROOTHANDLE,
	DefaultClass: uint32(HTBDEFAULTCLASS),
}

var ParentClass = HTBClass{
	Handle:       HTBPARENTCLASSHANDLE,
	ParentHandle: HTBQDISCHANDLE,
	Rate:         1_000_000,
}

var HighClass = HTBClass{
	Handle:       HTBHIGHPRIOCLASSHANDLE,
	ParentHandle: HTBPARENTCLASSHANDLE,
	Priority:     Priority(HTBHIGHCLASSPRIO),
	Rate:         600_000,
}

var LowClass = HTBClass{
	Handle:       HTBLOWPRIOCLASSHANDLE,
	ParentHandle: HTBPARENTCLASSHANDLE,
	Priority:     Priority(HTBLOWCLASSPRIO),
	Rate:         100_000,
}

var DefaultClass = HTBClass{
	Handle:       HTBDEFAULTCLASSHANDLE,
	ParentHandle: HTBPARENTCLASSHANDLE,
	Priority:     Priority(HTBDEFAULTCLASSPRIO),
	Rate:         400_000,
}

var HighPrioClassFilter = FWFilter{
	Handle:       filter.HIGHPRIOMARK,
	ParentHandle: HTBQDISCHANDLE,
	ClassID:      HTBHIGHPRIOCLASSHANDLE,
}

var LowPrioClassFilter = FWFilter{
	Handle: filter.LOWPRIOMARK,

	ParentHandle: HTBQDISCHANDLE,
	ClassID:      HTBLOWPRIOCLASSHANDLE,
}

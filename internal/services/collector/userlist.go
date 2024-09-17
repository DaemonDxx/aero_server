package collector

import "github.com/daemondxx/lks_back/entity"

type element struct {
	next *element
	prev *element
	v    *entity.User
}

type userList struct {
	head *element
	last *element
	len  int
}

func newUserList(users []entity.User) *userList {
	l := userList{}
	for i := range users {
		l.Insert(&users[i])
	}

	return &l
}

func (l *userList) First() *element {
	return l.head
}

func (l *userList) Insert(u *entity.User) {
	el := &element{
		next: nil,
		prev: nil,
		v:    u,
	}

	if l.head == nil {
		l.head = el
		l.last = el
	} else {
		l.last.next = el
		el.prev = l.last
		l.last = el
	}

	l.len++
}

func (l *userList) Remove(el *element) {
	if l.head == nil {
		return
	}

	if el == l.head {
		if l.len == 1 {
			l.head = nil
			l.last = nil
		} else {
			l.head = l.head.next
			l.head.prev = nil
		}
	} else if el == l.last {
		l.last = l.last.prev
		l.last.next = nil
	} else {
		el.prev.next = el.next
		el.next.prev = el.prev
	}

	el.v = nil
	el.prev = nil
	el.next = nil

	l.len--
}

func (l *userList) Len() int {
	return l.len
}

func (l *userList) Array() []*entity.User {
	arr := make([]*entity.User, 0, l.len)
	el := l.First()
	if el == nil {
		return arr
	}

	for {
		arr = append(arr, el.v)
		if el.next != nil {
			el = el.next
		} else {
			return arr
		}
	}
}

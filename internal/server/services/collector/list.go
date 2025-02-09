package collector

import "github.com/daemondxx/lks_back/entity"

type element struct {
	next *element
	prev *element
	v    *entity.Credential
}

type credentialList struct {
	head *element
	last *element
	len  int
}

func newCredentialList(creds []entity.Credential) *credentialList {
	l := credentialList{}
	for i := range creds {
		l.Insert(&creds[i])
	}

	return &l
}

func (l *credentialList) First() *element {
	return l.head
}

func (l *credentialList) Insert(u *entity.Credential) {
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

func (l *credentialList) Remove(el *element) {
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

func (l *credentialList) Len() int {
	return l.len
}

func (l *credentialList) Array() []*entity.Credential {
	arr := make([]*entity.Credential, 0, l.len)
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

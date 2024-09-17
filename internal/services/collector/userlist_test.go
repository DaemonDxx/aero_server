package collector

import (
	"github.com/daemondxx/lks_back/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var casesInsert = []struct {
	name         string
	insert       []*entity.User
	init         []entity.User
	resultLength int
}{
	{
		name:         "init list",
		insert:       []*entity.User{},
		init:         []entity.User{entity.User{ID: 0}, entity.User{ID: 1}},
		resultLength: 2,
	},
	{
		name:         "insert one object",
		insert:       []*entity.User{&entity.User{ID: 0}},
		init:         []entity.User{},
		resultLength: 1,
	},
	{
		name:         "insert some object with init array",
		insert:       []*entity.User{&entity.User{ID: 3}, &entity.User{ID: 4}},
		init:         []entity.User{entity.User{ID: 0}, entity.User{ID: 1}},
		resultLength: 4,
	},
}

func TestUserList_Insert(t *testing.T) {
	for _, tt := range casesInsert {
		t.Run(tt.name, func(t *testing.T) {
			l := newUserList(tt.init)
			for _, u := range tt.insert {
				l.Insert(u)
			}
			assert.Equal(t, tt.resultLength, l.len)
		})
	}
}

var casesRemove = []struct {
	name        string
	init        []entity.User
	removeIndex []int
	resultIDs   []uint
}{
	{
		name:        "remove head el",
		init:        []entity.User{entity.User{ID: 0}, entity.User{ID: 1}, entity.User{ID: 2}},
		removeIndex: []int{0},
		resultIDs:   []uint{1, 2},
	},
	{
		name:        "remove last el",
		init:        []entity.User{entity.User{ID: 0}},
		removeIndex: []int{0},
		resultIDs:   []uint{},
	},
	{
		name:        "remove last",
		init:        []entity.User{entity.User{ID: 0}, entity.User{ID: 1}, entity.User{ID: 2}},
		removeIndex: []int{2},
		resultIDs:   []uint{0, 1},
	},
	{
		name:        "remove middle",
		init:        []entity.User{entity.User{ID: 0}, entity.User{ID: 1}, entity.User{ID: 2}, entity.User{ID: 3}},
		removeIndex: []int{1, 2},
		resultIDs:   []uint{0, 3},
	},
}

func TestUserList_Remove(t *testing.T) {
	for _, tt := range casesRemove {
		t.Run(tt.name, func(t *testing.T) {
			l := newUserList(tt.init)

			el := l.First()
			var next *element
			rmi := tt.removeIndex[0]
			rmiNext := 1

			for i := 0; i < len(l.Array()); i++ {
				next = el.next

				if i == rmi {
					l.Remove(el)
					if rmiNext >= len(tt.removeIndex) {
						break
					} else {
						rmi = tt.removeIndex[rmiNext]
						rmiNext++
					}
				}

				if next == nil {
					break
				}

				el = next
			}
			res := l.Array()
			require.Equal(t, len(tt.resultIDs), len(res))

			for i := 0; i < len(tt.resultIDs); i++ {
				assert.Equal(t, tt.resultIDs[i], res[i].ID)
			}
		})
	}
}

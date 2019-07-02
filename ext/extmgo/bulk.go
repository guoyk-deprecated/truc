package extmgo

import "github.com/globalsign/mgo"

type Bulk struct {
	coll *mgo.Collection
	size int
	docs []interface{}
}

func NewBulk(coll *mgo.Collection, size int) (b *Bulk) {
	b = &Bulk{coll: coll, size: size}
	if b.size < 1 {
		b.size = 1
	}
	b.docs = make([]interface{}, 0, b.size)
	return
}

func (b *Bulk) Append(doc interface{}) (err error) {
	b.docs = append(b.docs, doc)
	if len(b.docs) >= b.size {
		if err = b.coll.Insert(b.docs...); err != nil {
			return
		}
		b.docs = b.docs[0:0]
	}
	return
}

func (b *Bulk) Finish() (err error) {
	if len(b.docs) > 0 {
		if err = b.coll.Insert(b.docs...); err != nil {
			return
		}
		b.docs = b.docs[0:0]
	}
	return
}

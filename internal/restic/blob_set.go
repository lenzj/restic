package restic

import "sort"

// A BlobSet is a set of BlobHandles.
// The zero BlobSet is the empty set.
type BlobSet struct {
	byType [NumBlobTypes]map[ID]struct{}
}

// Has returns true iff id is contained in the set.
func (s *BlobSet) Has(h BlobHandle) bool {
	m := s.byType[h.Type]
	if m == nil {
		return false
	}
	_, ok := m[h.ID]
	return ok
}

// Insert adds id to the set.
func (s *BlobSet) Insert(h BlobHandle) {
	m := s.byType[h.Type]
	if m == nil {
		m = make(map[ID]struct{})
		s.byType[h.Type] = m
	}
	m[h.ID] = struct{}{}
}

// Delete removes id from the set.
func (s BlobSet) Delete(h BlobHandle) {
	if m := s.byType[h.Type]; m != nil {
		delete(m, h.ID)
	}
}

// Equals returns true iff s equals other.
func (s *BlobSet) Equals(other *BlobSet) bool {
	for t := range s.byType {
		if len(s.byType[t]) != len(other.byType[t]) {
			return false
		}
	}

	for t, st := range s.byType {
		ot := other.byType[t]

		for id := range st {
			if _, ok := ot[id]; !ok {
				return false
			}
		}
	}

	return true
}

func (s *BlobSet) ForEach(fn func(BlobHandle) error) error {
	for t, m := range s.byType {
		for id := range m {
			err := fn(BlobHandle{Type: BlobType(t), ID: id})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Len returns the number of BlobHandles in the set s.
func (s *BlobSet) Len() (n int) {
	for _, m := range s.byType {
		n += len(m)
	}
	return n
}

// Merge adds the blobs in other to the current set.
func (s *BlobSet) Merge(other *BlobSet) {
	for t, st := range s.byType {
		ot := other.byType[t]

		for id := range ot {
			st[id] = struct{}{}
		}
	}
}

// Intersect returns a new set containing the handles that are present in both sets.
func (s BlobSet) Intersect(other BlobSet) (result BlobSet) {
	for t, m1 := range s.byType {
		m2 := other.byType[t]

		// Iterate over the smaller map.
		if len(m2) < len(m1) {
			m1, m2 = m2, m1
		}

		for id := range m1 {
			if _, ok := m2[id]; ok {
				result.Insert(BlobHandle{Type: BlobType(t), ID: id})
			}
		}
	}

	return result
}

// Sub returns a new set containing all handles that are present in s but not in
// other.
func (s BlobSet) Sub(other BlobSet) (result BlobSet) {
	for t, st := range s.byType {
		ot := other.byType[t]

		for id := range st {
			if _, ok := ot[id]; !ok {
				result.Insert(BlobHandle{Type: BlobType(t), ID: id})
			}
		}
	}

	return result
}

// List returns a sorted slice of all BlobHandle in the set.
func (s *BlobSet) List() BlobHandles {
	list := make(BlobHandles, 0, s.Len())
	for t, m := range s.byType {
		for id := range m {
			list = append(list, BlobHandle{Type: BlobType(t), ID: id})
		}
	}

	sort.Sort(list)

	return list
}

func (s BlobSet) String() string {
	str := s.List().String()
	if len(str) < 2 {
		return "{}"
	}

	return "{" + str[1:len(str)-1] + "}"
}

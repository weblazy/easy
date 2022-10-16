package set_test

// func TestStringSet(t *testing.T) {
// 	s := set.NewStringSet()
// 	assert.Equal(t, 0, s.Size())

// 	s.Add("1")
// 	assert.Equal(t, 1, s.Size())
// 	assert.True(t, s.Has("1"))
// 	s.Add("1")
// 	assert.Equal(t, 1, s.Size())
// 	assert.True(t, s.Has("1"))
// 	s.Delete("1")
// 	assert.Equal(t, 0, s.Size())
// 	assert.False(t, s.Has("1"))

// 	s.BatchAdd("2", "3")
// 	assert.True(t, s.Has("2"))
// 	assert.True(t, s.Has("3"))
// 	assert.Equal(t, 2, s.Size())
// 	assert.Equal(t, "2,3", s.String())
// 	assert.Equal(t, []string{"2", "3"}, s.ToArray())
// }

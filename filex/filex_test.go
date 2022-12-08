package filex

import "testing"

func TestWrite(t *testing.T) {
	path := "test.csv"

	err := Write(path, []byte("666"), false)
	if err != nil {
		t.Fatalf("err:%#v,", err)
	}
}

func TestRead(t *testing.T) {
	path := "test.csv"
	excepted := "666"
	ret, _ := Read(path)
	result := string(ret)
	if result != excepted {
		t.Fatalf("result:%#v,excepted:%#v", result, excepted)
	}
}

package ac

type AcNode struct {
	fail   *AcNode
	next   map[byte]*AcNode
	length []int
}

func newAcNode() *AcNode {
	return &AcNode{
		fail: nil,
		next: map[byte]*AcNode{},
	}
}

type AcAutoMachine struct {
	root *AcNode
	size int64
}

func NewAcAutoMachine() *AcAutoMachine {
	return &AcAutoMachine{
		root: newAcNode(),
	}
}

// 构造前缀树
func (ac *AcAutoMachine) AddPattern(pattern string) {
	chars := []byte(pattern)
	iter := ac.root
	var length int
	for _, c := range chars {
		if _, ok := iter.next[c]; !ok {
			iter.next[c] = newAcNode()
		}
		iter = iter.next[c]
		length++
	}
	iter.length = append(iter.length, length)
	ac.size++
}

// 构建fail指针
func (ac *AcAutoMachine) Build() {
	queue := []*AcNode{}
	queue = append(queue, ac.root)
	// 规则1:对节点层序遍历
	for len(queue) != 0 {
		parent := queue[0]
		queue = queue[1:]
		// 遍历第一个元素的子节点
		for char, child := range parent.next {
			if parent == ac.root {
				// fail指针规则2:第二层节点(根节点的孩子)fail指针指向根节点
				child.fail = ac.root
			} else {
				// 规则3:查找父节点的fail指针是否有与自己相同的子节点
				failAcNode := parent.fail
				for failAcNode != nil {
					if _, ok := failAcNode.next[char]; ok {
						child.fail = failAcNode.next[char]
						child.length = append(child.length, failAcNode.next[char].length...)
						break
					}
					failAcNode = failAcNode.fail
				}
				// failAcNode == ac.root找不到匹配的fail指针,指向根节点
				if failAcNode == nil {
					child.fail = ac.root
				}
			}
			queue = append(queue, child)
		}
	}
}

// 匹配敏感词
func (ac *AcAutoMachine) Query(content string) (results []string) {
	chars := []byte(content)
	iter := ac.root

	respMap := make(map[string]bool)
	data := []string{}
	for i, c := range chars {
		_, ok := iter.next[c]
		for !ok && iter != ac.root {
			// 匹配失败从fail指针开始尝试子串
			iter = iter.fail
			_, ok = iter.next[c]
		}

		iter = iter.next[c]
		if iter == nil {
			iter = ac.root
		}
		parent := iter
		if parent != ac.root && len(parent.length) > 0 {
			// 匹配成功
			for _, length := range parent.length {
				respMap[string([]byte(content)[i+1-length:i+1])] = true
				data = append(data, string([]byte(content)[i+1-length:i+1]))
			}
		}
	}
	for word := range respMap {
		results = append(results, word)
	}
	return
}

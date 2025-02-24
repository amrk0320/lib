package rtree

import (
	"fmt"
	"math"
)

type (
	RTree struct {
		Root *Node
		cnf  *Config
	}
	Config struct {
		MaxEntrySize int
	}
	Node struct {
		Tree      *RTree
		Parent    *Node
		Rectangle Rectangle
		DataID    *uint64 // リーフエントリーのみ存在
		Children  Nodes
		// depth     uint8
	}

	Inteval struct {
		First  float64
		Second float64
	}

	Nodes     []*Node
	Rectangle []*Inteval
)

const (
	dimCount = 2
)

func NewRTree(cnf *Config) (result *RTree) {
	result = new(RTree)
	result.cnf = cnf
	result.Root = result.NewNode(nil)
	result.Root.Rectangle = maxRectangle()

	return
}

// ゼロ空間
func zeroRectangle() Rectangle {
	return Rectangle{&Inteval{First: 0.0, Second: 0.0}, &Inteval{First: 0.0, Second: 0.0}}
}

// 最大空間
func maxRectangle() Rectangle {
	return Rectangle{&Inteval{First: 0.0, Second: math.MaxFloat64}, &Inteval{First: 0.0, Second: math.MaxFloat64}}
}

func min(a, b float64) float64 {
	if a < b {
		return a
	} else {
		return b
	}
}

func max(a, b float64) float64 {
	if a > b {
		return a
	} else {
		return b
	}
}

func (tree *RTree) NewNode(parent *Node) (node *Node) {
	node = new(Node)
	node.Tree = tree
	node.Parent = parent
	node.Children = make([]*Node, 0, tree.cnf.MaxEntrySize)
	node.Rectangle = make(Rectangle, dimCount)
	node.Rectangle = zeroRectangle()

	return
}

func (node *Node) Print(depth uint8) {
	line := ""
	loop := int(depth + 1)

	for i := 0; i < loop; i++ {
		line += "-"
	}

	fmt.Println(line)
	fmt.Println("node", node)

	if node.DataID != nil {
		fmt.Println("dataID", *node.DataID)
	}

	fmt.Println("rectangle ", node.Rectangle[0], node.Rectangle[1])
	fmt.Println("depth", depth)
	fmt.Println("children size", len(node.Children))

	for _, v := range node.Children {
		v.Print(depth + 1)
	}
}

func (node *Node) isLeaf() bool {
	return node.hasDataNode() || (node.isRoot() && len(node.Children) == 0)
}

func (node *Node) hasDataNode() bool {
	for i := range node.Children {
		if node.Children[i] != nil && node.Children[i].DataID != nil {
			return true
		}
	}

	return false
}

func (node *Node) isRoot() bool {
	return node.Parent == nil
}

func (node *Node) AddEntry(src *Node) {
	node.Children = append(node.Children, src)
	src.Parent = node
}

func (node *Node) deleteAllEntry() {
	node.Children = nil
}

// エントリー数が上限を超えているか判定する
func (node *Node) isOverFlow() bool {
	return node.Tree.cnf.MaxEntrySize < len(node.Children)
}

func (node *Node) isFullEntry() bool { //nolint
	return node.Tree.cnf.MaxEntrySize == len(node.Children)
}

// 次元毎の区間長
func (node *Node) listIntervalDistance() (distance []float64) {
	distance = make([]float64, len(node.Rectangle))

	for dim := range node.Rectangle {
		distance[dim] = node.Rectangle[dim].Second - node.Rectangle[dim].First
	}

	return
}

func (nodes *Nodes) delete(src *Node) {
	for i := range *nodes {
		if (*nodes)[i] == src {
			*nodes = append((*nodes)[:i], (*nodes)[i+1:]...)
			return
		}
	}
}

// 区間が重なるか判定
func (internal *Inteval) overlap(other Inteval) bool {
	return max(internal.First, other.First) <= min(internal.Second, other.Second)
}

// 重なる区間
func (internal *Inteval) overlapArea(other Inteval) (area float64) {
	return min(internal.Second, other.Second) - max(internal.First, other.First)
}

// 短形が重なるか判定
func (rectangle *Rectangle) overlap(other Rectangle) bool {
	for i := range *rectangle {
		// 重ならない
		if !(*rectangle)[i].overlap(*other[i]) {
			return false
		}
	}

	// 全て重なる
	return true
}

// 重なる面積
func (rectangle Rectangle) overlapArea(other Rectangle) (area float64) {
	for i := range rectangle {
		if !rectangle[i].overlap(*other[i]) {
			return 0
		}

		if area == 0.0 {
			area = 1
		}

		area *= rectangle[i].overlapArea(*other[i])
	}

	return
}

// 短形間の中心同士の距離
func (rectangle Rectangle) distance(other Rectangle) (distance float64) {
	for i := range rectangle {
		mid1 := (rectangle[i].First + rectangle[i].Second) / 2
		mid2 := (other[i].First + other[i].Second) / 2

		d := (mid1 - mid2)

		distance += d * d
	}

	return math.Sqrt(distance)
}

// 区間を完全に包含する
func (internal Inteval) cover(other Inteval) bool {
	return max(internal.First, other.First) <= min(other.Second, internal.Second)
}

// 短形を包含するか判定
func (rectangle *Rectangle) cover(other Rectangle) bool {
	for i := range *rectangle {
		// 重ならない
		if !(*rectangle)[i].cover(*other[i]) {
			return false
		}
	}

	// 全て重なる
	return true
}

// 探索短形が重なっている区間のIDを返却する
func (tree *RTree) FindAreas(root *Node, rectangle Rectangle) (results []*uint64, err error) {
	if root != nil {
		switch {
		case root.isLeaf():
			for i := range root.Children {
				if root.Children[i].Rectangle.cover(rectangle) {
					return []*uint64{root.Children[i].DataID}, nil
				}
			}
		default:
			for i := range root.Children {
				if root.Children[i].Rectangle.overlap(rectangle) {
					areas, err := tree.FindAreas(root.Children[i], rectangle)
					if err != nil {
						return nil, err
					}

					results = append(results, areas...)
				}
			}
		}
	}

	return
}

// ノードを挿入する
func (tree *RTree) AddNode(src *Node) (err error) {
	leaf, err := tree.findLeaf(tree.Root, src)
	if err != nil {
		return err
	}

	leaf.AddEntry(src)

	var newLeaf *Node

	if !leaf.isRoot() && leaf.isOverFlow() {
		newLeaf = leaf.SplitNode()
	}

	leaf.Adjust(newLeaf)

	tree.Root.AdjustRoot()

	return
}

func (node *Node) SplitNode() (newNode *Node) {
	// Linear-Cost Algorithm O(M)
	// https://tanishiking24.hatenablog.com/entry/introduction_rtree_index
	newNode = node.Tree.NewNode(node.Parent)

	if node.Parent != nil {
		node.Parent.AddEntry(newNode)
	}

	existsChildren := make(Nodes, len(node.Children))
	copy(existsChildren, node.Children)

	node.deleteAllEntry()

	baseDistances := node.listIntervalDistance()

	for 0 < len(existsChildren) {
		one, another := existsChildren.GetFarthestChildren(baseDistances)
		node.AddEntry(one)
		existsChildren.delete(one)

		if another != nil {
			newNode.AddEntry(another)
			existsChildren.delete(another)
		}
	}

	return
}

// 全ての次元内で最も離れたエントリーを取得
func (nodes *Nodes) GetFarthestChildren(distancesInDim []float64) (one *Node, another *Node) {
	maxDistance := -1.0

	for dim, baseDistance := range distancesInDim {
		distance, tmpOne, tmpAnother := nodes.getFarthestChildrenInDim(dim, baseDistance)

		if maxDistance < distance {
			one, another = tmpOne, tmpAnother
			maxDistance = distance
		}
	}

	return
}

// ある次元内で最も離れたエントリーを取得
func (nodes *Nodes) getFarthestChildrenInDim(dim int, baseDistance float64) (distance float64, one *Node, another *Node) {
	if len(*nodes) == 1 {
		return 0, (*nodes)[0], nil
	}

	farthestPairs := make([]int, 2)

	minSecond := math.MaxFloat64
	maxFirst := -1.0

	for j := range *nodes {
		if (*nodes)[j].Rectangle[dim].Second < minSecond {
			minSecond = (*nodes)[j].Rectangle[dim].Second
			farthestPairs[0] = j
		}

		if maxFirst < (*nodes)[j].Rectangle[dim].First {
			maxFirst = (*nodes)[j].Rectangle[dim].First
			farthestPairs[1] = j
		}
	}

	distance = (maxFirst - minSecond) / baseDistance
	one = (*nodes)[farthestPairs[0]]
	another = (*nodes)[farthestPairs[1]]

	return
}

func (node *Node) Adjust(newNode *Node) {
	// ルートに到達したら終了
	if !node.isRoot() {
		node.AdjustCoverRectangles()

		if newNode != nil {
			newNode.AdjustCoverRectangles()
		}

		var newParentNode *Node

		// 親がオーバーフローしたら分割する
		if !node.Parent.isRoot() && node.Parent.isOverFlow() {
			newParentNode = node.Parent.SplitNode()
		}

		node.Parent.Adjust(newParentNode)
	}
}

// ルートの調整、ルートがオーバーフローしたら分割する
func (node *Node) AdjustRoot() {
	// ルートの調整
	node.AdjustCoverRectangles()

	if node.isOverFlow() {
		oldRoot := node
		// 現在のルートを分割する
		anotherOldRoot := oldRoot.SplitNode()
		oldRoot.AdjustCoverRectangles()
		anotherOldRoot.AdjustCoverRectangles()

		// 新しいルートを作成する
		newRoot := oldRoot.Tree.NewNode(nil)

		// 子供となった旧ルートの調整
		oldRoot.Parent = newRoot
		anotherOldRoot.Parent = newRoot

		newRoot.AddEntry(oldRoot)
		newRoot.AddEntry(anotherOldRoot)

		// 新しいルートの木への設定と調整
		newRoot.AdjustCoverRectangles()
		node.Tree.Root = newRoot
	}
}

// // 親のエントリーに自ノードを追加する
// func (node *Node) joinParentEntry() {
// 	node.Parent.addEntry(node)
// }

// エントリー要素の区間を包含するように親区間を修正する
func (node *Node) AdjustCoverRectangles() {
	for dim := range node.Rectangle {
		newFirst := math.MaxFloat64
		newSecond := -1.0

		for i := range node.Children {
			newFirst = min(newFirst, node.Children[i].Rectangle[dim].First)
			newSecond = max(newSecond, node.Children[i].Rectangle[dim].Second)
		}

		node.Rectangle[dim].First = newFirst
		node.Rectangle[dim].Second = newSecond
	}
}

// ノードを挿入する葉ノードを探索する
func (tree *RTree) findLeaf(root *Node, src *Node) (node *Node, err error) {
	switch {
	case root.isLeaf():
		node = root // for lint
	//	return root, nil
	default:
		var (
			next        *Node
			minDistance float64
			nearestNode *Node
			maxArea     float64
		)

		minDistance = math.MaxFloat64

		for i := range root.Children {
			tmpArea := root.Children[i].Rectangle.overlapArea(src.Rectangle)
			// 重なりあり
			if maxArea < tmpArea {
				next = root.Children[i]
				maxArea = tmpArea
			} else {
				// 重なりなし
				distance := root.Children[i].Rectangle.distance(src.Rectangle)
				if distance < minDistance {
					minDistance = distance
					nearestNode = root.Children[i]
				}
			}
		}

		// 該当範囲がないなら最も近いノードを探索する
		if next == nil {
			next = nearestNode
		}

		node, err = tree.findLeaf(next, src)
	}

	return //nolint
}

func (tree *RTree) TakePlace(id uint64, lat, lon float64) (node *Node) {
	node = tree.NewNode(nil)

	rectangle := Rectangle{
		&Inteval{First: lat, Second: lat},
		&Inteval{First: lon, Second: lon},
	}

	node.Rectangle = rectangle
	node.DataID = &id

	return
}

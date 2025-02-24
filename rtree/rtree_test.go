package rtree_test

import (
	"math"
	"rtree"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindAreas(t *testing.T) {

	t.Run("add one, search one", func(t *testing.T) {
		tree := rtree.NewRTree(&rtree.Config{MaxEntrySize: 1})

		place := tree.TakePlace(1, 1, 1)

		_ = tree.AddNode(place)

		ids, err := tree.FindAreas(tree.Root, place.Rectangle)
		assert.NoError(t, err)
		assert.Contains(t, ids, place.DataID)
	})

	t.Run("add and overflow, search one, area", func(t *testing.T) {
		tree := rtree.NewRTree(&rtree.Config{MaxEntrySize: 2})

		place := []*rtree.Node{
			tree.TakePlace(1, 1, 1),
			tree.TakePlace(2, 1, 2),
			tree.TakePlace(3, 2, 3),
		}

		for _, p := range place {
			_ = tree.AddNode(p)
		}

		tree.Root.Print(0)

		// root
		// node node
		// 1 2  3

		// search one
		ids, err := tree.FindAreas(tree.Root, place[1].Rectangle)
		assert.NoError(t, err)
		assert.Contains(t, ids, place[1].DataID)

		// search area
		ids, err = tree.FindAreas(tree.Root, rtree.Rectangle{&rtree.Inteval{First: 1, Second: 2}, &rtree.Inteval{First: 2, Second: 3}})
		assert.NoError(t, err)
		assert.ElementsMatch(t, ids, []*uint64{place[1].DataID, place[2].DataID})
	})
}

func TestAddNode(t *testing.T) {
	t.Run("add children node", func(t *testing.T) {
		tree := rtree.NewRTree(&rtree.Config{MaxEntrySize: 3})

		place := rtree.Nodes{
			tree.TakePlace(1, 1, 1),
			tree.TakePlace(2, 50, 34),
			tree.TakePlace(3, 3, 2000),
		}

		for _, p := range place {
			_ = tree.AddNode(p)
		}

		tree.Root.Print(0)

		// root
		// 1 2 3

		assert.Equal(t, 3, len(tree.Root.Children))
		assert.ElementsMatch(t, place, tree.Root.Children)
		assert.EqualValues(t, 1, tree.Root.Rectangle[0].First)
		assert.EqualValues(t, 50, tree.Root.Rectangle[0].Second)
		assert.EqualValues(t, 1, tree.Root.Rectangle[1].First)
		assert.EqualValues(t, 2000, tree.Root.Rectangle[1].Second)
	})

	t.Run("add grand grand child node", func(t *testing.T) {
		tree := rtree.NewRTree(&rtree.Config{MaxEntrySize: 2})

		place := []*rtree.Node{
			tree.TakePlace(1, 1, 1),    // 子１
			tree.TakePlace(2, 2, 2),    // 子１
			tree.TakePlace(3, 3, 2000), //
			tree.TakePlace(4, 1, 2),    //
		}

		for _, p := range place {
			_ = tree.AddNode(p)
		}

		tree.Root.Print(0)

		// root
		// node      node
		// node node node
		// 1 4  2     3

		assert.NotNil(t, tree.Root)                      // depth 0
		assert.Equal(t, 2, len(tree.Root.Children))      // depth 1
		assert.NotNil(t, tree.Root.Children[0])          // depth 1
		assert.NotNil(t, tree.Root.Children[0].Children) // depth 2
		assert.NotNil(t, tree.Root.Children[1].Children) // depth 2

		assert.NotNil(t, tree.Root.Children[0].Children[0].Children[0].DataID) // depth 3
		assert.NotNil(t, tree.Root.Children[0].Children[1].Children[0].DataID) // depth 3
	})
}

// t.Run("add cover grand child node", func(t *testing.T) {
// 	rtree := NewRTree(&Config{maxEntrySize: 2})

// 	place := []*Node{
// 		rtree.takePlace(1, 1, 1),    // 子１
// 		rtree.takePlace(2, 50, 34),  // 子１
// 		rtree.takePlace(3, 3, 2000), // 子１-孫１
// 		rtree.takePlace(4, 3, 2001), // 子１-孫２
// 	}

// 	for _, p := range place {
// 		rtree.AddNode(p)
// 	}

// 	assert.Equal(t, 2, len(rtree.Root.Children))
// 	assert.NotNil(t, rtree.Root.Children[0])
// 	assert.NotNil(t, rtree.Root.Children[1])
// 	assert.Equal(t, 0, len(rtree.Root.Children[0].Children))
// 	assert.Equal(t, 2, len(rtree.Root.Children[1].Children))
// })

func TestGetFarthestEntries(t *testing.T) {
	t.Run("take fatherest pair", func(t *testing.T) {
		tree := rtree.NewRTree(&rtree.Config{MaxEntrySize: 3})

		entries := rtree.Nodes{
			tree.TakePlace(0, 0, 1), // ok 緯度で最も最も遠い
			tree.TakePlace(0, 1, 0), // 経度で最も最も遠い
			tree.TakePlace(0, 1, 2),
			tree.TakePlace(0, 1, 50),  // 経度で最も最も遠い
			tree.TakePlace(0, 100, 0), // ok 緯度で最も最も遠い
		}

		one, another := entries.GetFarthestChildren([]float64{math.MaxFloat64, math.MaxFloat64})

		assert.Equal(t, one, entries[0])
		assert.Equal(t, another, entries[4])
	})
}

func _TestAdjustCoverRectangles(t *testing.T) {
	t.Run("adjsut rectangle", func(t *testing.T) {
		tree := rtree.NewRTree(&rtree.Config{MaxEntrySize: 3})

		entries := rtree.Nodes{
			tree.TakePlace(0, 32.789789, 23.7878),
			tree.TakePlace(0, 2.4343, 88.9090),
			tree.TakePlace(0, 9.00000, 23.8989),
			tree.TakePlace(0, 9.000001, 21.877),
			tree.TakePlace(0, 2.4344, 135.999),
		}

		n := tree.NewNode(nil)

		for _, e := range entries {
			n.AddEntry(e)
		}

		n.AdjustCoverRectangles()
		assert.Equal(t, 2.4343, n.Rectangle[0].First)
		assert.Equal(t, 32.789789, n.Rectangle[0].Second)
		assert.Equal(t, 21.877, n.Rectangle[1].First)
		assert.Equal(t, 135.999, n.Rectangle[1].Second)
	})
}

func _TestSplitNode(t *testing.T) {
	t.Run("SplitNode 5", func(t *testing.T) {
		tree := rtree.NewRTree(&rtree.Config{MaxEntrySize: 5})

		nodes := rtree.Nodes{
			tree.TakePlace(0, 0, 1),   // グループ１ 緯度で最も最も遠い
			tree.TakePlace(1, 1, 0),   // グループ１ 経度で最も最も遠い
			tree.TakePlace(2, 1, 2),   // グループ１ 残り、どっちでもいい
			tree.TakePlace(3, 1, 50),  // グループ２ 経度で最も最も遠い
			tree.TakePlace(4, 100, 0), // グループ２ 緯度で最も最も遠い
		}

		n := tree.NewNode(tree.Root)

		for _, e := range nodes {
			n.AddEntry(e)
		}

		newNode := n.SplitNode()

		assert.Equal(t, 3, len(n.Children))
		assert.Equal(t, 2, len(newNode.Children))
		assert.ElementsMatch(t, rtree.Nodes{
			nodes[0],
			nodes[1],
			nodes[2],
		}, n.Children)

		assert.ElementsMatch(t, rtree.Nodes{
			nodes[3],
			nodes[4],
		}, newNode.Children)
	})

	t.Run("SplitNode 4", func(t *testing.T) {
		tree := rtree.NewRTree(&rtree.Config{MaxEntrySize: 4})

		nodes := rtree.Nodes{
			tree.TakePlace(0, 0, 1),
			tree.TakePlace(1, 3000, 3000),
			tree.TakePlace(3, 4000, 4000),
			tree.TakePlace(4, 10, 100),
		}

		n := tree.NewNode(tree.Root)

		for _, e := range nodes {
			n.AddEntry(e)
		}

		newNode := n.SplitNode()

		assert.Equal(t, 2, len(n.Children))
		assert.Equal(t, 2, len(newNode.Children))
		assert.ElementsMatch(t, rtree.Nodes{
			nodes[0],
			nodes[3],
		}, n.Children)

		assert.ElementsMatch(t, rtree.Nodes{
			nodes[1],
			nodes[2],
		}, newNode.Children)
	})
}

func _TestAdjust(t *testing.T) {
	t.Run("adjust node", func(t *testing.T) {
		tree := rtree.NewRTree(&rtree.Config{MaxEntrySize: 2})

		nodes := []*rtree.Node{
			tree.TakePlace(0, 100, 100),
			tree.TakePlace(1, 200, 200),
			tree.TakePlace(2, 10000, 20000),
		}

		// 適当な中間ノード
		node := tree.NewNode(tree.Root)

		for i := range nodes {
			node.AddEntry(nodes[i])
		}

		// オーバーフローしているので分割
		newNode := node.SplitNode()

		node.Adjust(newNode)

		assert.EqualValues(t, 100, node.Rectangle[0].First)
		assert.EqualValues(t, 200, node.Rectangle[0].Second)
		assert.EqualValues(t, 100, node.Rectangle[1].First)
		assert.EqualValues(t, 200, node.Rectangle[1].Second)
		assert.EqualValues(t, 2, len(node.Children))
		assert.Equal(t, tree.Root, node.Parent)

		assert.EqualValues(t, 10000, newNode.Rectangle[0].First)
		assert.EqualValues(t, 10000, newNode.Rectangle[0].Second)
		assert.EqualValues(t, 20000, newNode.Rectangle[1].First)
		assert.EqualValues(t, 20000, newNode.Rectangle[1].Second)
		assert.EqualValues(t, 1, len(newNode.Children))
		assert.Equal(t, tree.Root, newNode.Parent)
	})

	t.Run("adjust parent(split root)", func(t *testing.T) {
		tree := rtree.NewRTree(&rtree.Config{MaxEntrySize: 2})

		nodes := []*rtree.Node{
			tree.TakePlace(0, 100, 100),
			tree.TakePlace(1, 200, 200),
			tree.TakePlace(2, 10000, 20000),
		}

		// 適当な中間ノード
		oldRoot := tree.Root

		for i := range nodes {
			oldRoot.AddEntry(nodes[i])
		}

		// オーバーフローしているのでルートを分割
		oldRoot.AdjustRoot()

		// ルートが変わっているか確認
		assert.NotEqual(t, tree.Root, oldRoot)
		assert.Nil(t, tree.Root.Parent)
		assert.Equal(t, tree.Root, tree.Root.Children[0].Parent)
		assert.Equal(t, tree.Root, tree.Root.Children[1].Parent)
		assert.EqualValues(t, 2, len(tree.Root.Children))

		// 新ルートの確認
		assert.EqualValues(t, 100, tree.Root.Rectangle[0].First)
		assert.EqualValues(t, 10000, tree.Root.Rectangle[0].Second)
		assert.EqualValues(t, 100, tree.Root.Rectangle[1].First)
		assert.EqualValues(t, 20000, tree.Root.Rectangle[1].Second)
	})
}

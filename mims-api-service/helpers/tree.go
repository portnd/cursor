package helpers

// // Tree
// type Tree struct {
// 	Id       int    `json:"id"`
// 	ParentId int    `json:"parent_id"`
// 	Title    string `json:"title"`
// 	Name     string `json:"name"`
// 	Route    string `json:"route"`
// 	Icon     string `json:"icon"`
// 	// Data       interface{} `json:"data"`
// 	IsChildren bool   `json:"is_children"`
// 	Children   []Tree `json:"children"`
// }

// // ConvertToINodeArray
// type INode interface {
// 	GetId() int
// 	GetParentId() int
// 	GetTitle() string
// 	GetName() string
// 	GetRoute() string
// 	GetIcon() string
// 	// GetData
// 	GetData() interface{}
// 	IsRoot() bool
// }
// type INodes []INode

// func (nodes INodes) Len() int {
// 	return len(nodes)
// }
// func (nodes INodes) Swap(i, j int) {
// 	nodes[i], nodes[j] = nodes[j], nodes[i]
// }
// func (nodes INodes) Less(i, j int) bool {
// 	return nodes[i].GetId() < nodes[j].GetId()
// }

// func GenerateTree(nodes, selectedNodes []INode) (trees []Tree) {
// 	trees = []Tree{}
// 	var roots, childs []INode
// 	for _, v := range nodes {
// 		if v.IsRoot() {
// 			roots = append(roots, v)
// 		}
// 		childs = append(childs, v)
// 	}

// 	for _, v := range roots {
// 		childTree := &Tree{
// 			Id:       v.GetId(),
// 			ParentId: v.GetParentId(),
// 			Title:    v.GetTitle(),
// 			Name:     v.GetName(),
// 			Route:    v.GetRoute(),
// 			Icon:     v.GetIcon(),

// 			// Data:  v.GetData(),
// 		}
// 		recursiveTree(childTree, childs, selectedNodes)
// 		childTree.IsChildren = len(childTree.Children) == 0
// 		trees = append(trees, *childTree)
// 	}
// 	return
// }

// func recursiveTree(tree *Tree, nodes, selectedNodes []INode) {
// 	// data := tree.Data.(INode)
// 	id := tree.Id

// 	for _, v := range nodes {
// 		if v.IsRoot() {
// 			continue
// 		}
// 		if id == v.GetParentId() {
// 			childTree := &Tree{
// 				Id:       v.GetId(),
// 				ParentId: v.GetParentId(),
// 				Title:    v.GetTitle(),
// 				Name:     v.GetName(),
// 				Route:    v.GetRoute(),
// 				Icon:     v.GetIcon(),

// 				// Data:  v.GetData(),
// 			}
// 			recursiveTree(childTree, nodes, selectedNodes)
// 			childTree.IsChildren = len(childTree.Children) == 0
// 			tree.Children = append(tree.Children, *childTree)
// 		}
// 	}
// }

// // func FindRelationNode(nodes, allNodes []INode) (respNodes []INode) {
// // 	nodeMap := make(map[int]INode)
// // 	for _, v := range nodes {
// // 		recursiveFindRelationNode(nodeMap, allNodes, v, 0)
// // 	}

// // 	for _, v := range nodeMap {
// // 		respNodes = append(respNodes, v)
// // 	}
// // 	sort.Sort(INodes(respNodes))
// // 	return
// // }

// // func recursiveFindRelationNode(nodeMap map[int]INode, allNodes []INode, node INode, t int) {
// // 	nodeMap[node.GetId()] = node
// // 	for _, v := range allNodes {
// // 		if _, ok := nodeMap[v.GetId()]; ok {
// // 			continue
// // 		}

// // 		if t == 0 || t == 1 {
// // 			if node.GetParentId() == v.GetId() {
// // 				nodeMap[v.GetId()] = v
// // 				if v.IsRoot() {
// // 					continue
// // 				}
// // 				recursiveFindRelationNode(nodeMap, allNodes, v, 1)
// // 			}
// // 		}
// // 		if t == 0 || t == 2 {
// // 			if node.GetId() == v.GetParentId() {
// // 				nodeMap[v.GetId()] = v
// // 				recursiveFindRelationNode(nodeMap, allNodes, v, 2)
// // 			}
// // 		}
// // 	}
// // }

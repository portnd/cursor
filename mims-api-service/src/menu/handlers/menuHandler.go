package handlers

import (
	"fmt"

	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/menu/domains"

	"github.com/gin-gonic/gin"
)

// init handler
type MenuHandler struct {
	menuUseCase domains.MenuUseCase
}

// init handler
func NewMenuHandler(usecase domains.MenuUseCase) *MenuHandler {
	return &MenuHandler{
		menuUseCase: usecase,
	}
}

// ================================== start function  ==================================

// request form
type LoginCredentials struct {
	Email    string `form:"email" validate:"min=1"`
	Password string `form:"password" validate:"min=1"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type SystemMenu struct {
	Id       int    `json:"id"`
	ParentId int    `json:"parent_id"`
	Name     string `json:"name"`
	Route    string `json:"route"`
	Icon     string `json:"icon"`
}

func (s SystemMenu) GetTitle() string {
	return s.Name
}
func (s SystemMenu) GetId() int {
	return s.Id
}
func (s SystemMenu) GetParentId() int {
	return s.ParentId
}

func (s SystemMenu) GetName() string {
	return s.Name
}

func (s SystemMenu) GetRoute() string {
	return s.Route
}

func (s SystemMenu) GetIcon() string {
	return s.Icon
}

func (s SystemMenu) GetData() interface{} {
	return s
}
func (s SystemMenu) IsRoot() bool {
	return s.ParentId == 0 || s.ParentId == s.Id
}

type SystemMenus []SystemMenu

func (s SystemMenus) ConvertToINodeArray() (nodes []INode) {
	for _, v := range s {
		nodes = append(nodes, v)
	}
	return
}

// @summary
// @description
// @tags menu
// @id get_menu
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} Tree "OK"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/menu [get]
func (t *MenuHandler) GetMenu(c *gin.Context) {
	accCtrl, _ := c.Get("accessControl")
	aInterface := accCtrl.([]interface{})
	accCtrls := make(map[string]string)

	for _, v := range aInterface {
		accCtrls[v.(string)] = v.(string)
	}
	// c.JSON(200, accCtrls)
	// return
	menu, err := t.menuUseCase.GetMenu(accCtrls)
	if err != nil {
		fmt.Println("error", err)
	}
	// c.JSON(200, menu)
	// return
	var systemMenu []SystemMenu
	for _, item := range menu {
		var sys SystemMenu
		sys.Id = item.Id
		sys.ParentId = item.ParentId
		sys.Name = item.Name
		sys.Route = item.Route
		sys.Icon = item.Icon
		systemMenu = append(systemMenu, sys)
	}
	// c.JSON(200, SystemMenus.ConvertToINodeArray(systemMenu))
	// return
	resp := t.GenerateTree(SystemMenus.ConvertToINodeArray(systemMenu), nil)
	// c.JSON(200, resp)
	c.JSON(200, responses.SuccessResponse(resp, 200))
	// return
	// bytes, _ := json.MarshalIndent(resp, "", "\t")
	// fmt.Println(string(pretty.Color(pretty.PrettyOptions(bytes, pretty.DefaultOptions), nil)))
	//
}

// --------------------------------------------- Tree ---------------------------------------------
type Tree struct {
	Id       int    `json:"id"`
	ParentId int    `json:"parent_id"`
	Title    string `json:"title"`
	Name     string `json:"name"`
	Route    string `json:"route"`
	Icon     string `json:"icon"`
	// Data       interface{} `json:"data"`
	IsChildren bool   `json:"is_children"`
	Children   []Tree `json:"children"`
}

// ConvertToINodeArray
type INode interface {
	GetId() int
	GetParentId() int
	// GetTitle() string
	GetName() string
	GetRoute() string
	GetIcon() string
	// GetData
	GetData() interface{}
	IsRoot() bool
}

type INodes []INode

func (nodes INodes) Len() int {
	return len(nodes)
}
func (nodes INodes) Swap(i, j int) {
	nodes[i], nodes[j] = nodes[j], nodes[i]
}
func (nodes INodes) Less(i, j int) bool {
	return nodes[i].GetId() < nodes[j].GetId()
}

func (t *MenuHandler) GenerateTree(nodes, selectedNodes []INode) (trees []Tree) {
	trees = []Tree{}
	var roots, childs []INode
	for _, v := range nodes {
		if v.IsRoot() {
			roots = append(roots, v)
		}
		childs = append(childs, v)
	}

	for _, v := range roots {
		childTree := &Tree{
			Id:       v.GetId(),
			ParentId: v.GetParentId(),
			// Title:    v.GetTitle(),
			Name:  v.GetName(),
			Route: v.GetRoute(),
			Icon:  v.GetIcon(),

			// Data:  v.GetData(),
		}
		t.recursiveTree(childTree, childs, selectedNodes)
		childTree.IsChildren = len(childTree.Children) == 0
		// fmt.Println("v.GetRoute()", v.GetRoute())
		if childTree.IsChildren && childTree.Route == "" {
			continue
		}
		trees = append(trees, *childTree)
	}
	return
}

func (t *MenuHandler) recursiveTree(tree *Tree, nodes, selectedNodes []INode) {
	// data := tree.Data.(INode)
	id := tree.Id

	for _, v := range nodes {
		if v.IsRoot() {
			continue
		}
		if id == v.GetParentId() {
			childTree := &Tree{
				Id:       v.GetId(),
				ParentId: v.GetParentId(),
				// Title:    v.GetTitle(),
				Name:  v.GetName(),
				Route: v.GetRoute(),
				Icon:  v.GetIcon(),

				// Data:  v.GetData(),
			}

			t.recursiveTree(childTree, nodes, selectedNodes)
			childTree.IsChildren = len(childTree.Children) == 0

			if childTree.IsChildren && childTree.Route == "" {
				fmt.Println("v.GetRoute()v.GetRoute()", v.GetRoute())
				continue
			}
			tree.Children = append(tree.Children, *childTree)
		}
	}
}

package hn

import "clx/item"

type Service interface {
	FetchItems(int, int) []*item.Item
	FetchItem(id int) *item.Item
	FetchComments(int) *item.Item
}

// FlatSortable 组件的 Props
export interface FlatSortableProps {
  dir?: boolean // 是否启用排序方向
  modal?: boolean // 是否启用模态行为
}

export type FlatSortableEmits = {
  'update:open': [value: boolean]
}

// FlatSortableContent 的 Props
export interface FlatSortableContentProps {
  direction?: 'row' | 'column'; // 布局方向，可选
  gap?: number; // 间距（像素）
}

// FlatSortableContent 的 Emits，与 FlatSortable 共享
export type FlatSortableContentEmits = FlatSortableEmits;

// FlatSortableItem 的 Props
export interface FlatSortableItemProps {

}

// FlatSortableItem 的 Emits
export type FlatSortableItemEmits = {
  'select': [event: Event] // 选择事件
};

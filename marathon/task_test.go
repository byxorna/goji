package marathon
import (
  "testing"
  "sort"
)

func TestSortTaskListById(t *testing.T){
  l := TaskList{
    Task{Id: "z"},Task{Id: "b"},Task{Id: "a"},Task{Id:"x"},
  }
  sort.Sort(l)
  s := ""
  for _, t := range l {
    s += t.Id
  }
  if s != "abxz" {
    t.Error("TaskList doesnt sort properly")
  }
}

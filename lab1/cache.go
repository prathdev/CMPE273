package main

// DO NOT CHANGE THIS CACHE SIZE VALUE
const CACHE_SIZE int = 3

type Node struct {
   	key int
    value int
    pre *Node
    next *Node
}
	
var items map[int]Node
var head *Node
var end *Node

func remove(n *Node){
        if n.pre!=nil {
            n.pre.next = n.next
        }else{
            head = n.next
        }
 
        if n.next!=nil {
            n.next.pre = n.pre
        }else{
            end = n.pre
        }
 
}
 
func setHead(n *Node){
        n.next = head
        n.pre = nil
	
 
		if head!=nil {
            head.pre = n
		}
 
        head = n
 
		if end == nil {
            end = head
		}
	
}


func Set(key int, value int) {
	// TODO: add your code here!
	if len(items)==0 {
			
			items = make(map[int]Node)
	}
		
		if val, ok := items[key]; ok {
			old := val
            old.value = value
            remove(&old)
            setHead(&old)
		
		} else {
		//remove(end)
		created := Node{key: key, value: value}
            if len(items)>=CACHE_SIZE {
				
				delete(items,end.key)
                
                remove(end)
                setHead(&created)
 
            }else{
				
                setHead(&created)
            }    
 
			
			items[key]=created
			            
        }
}

func Get(key int) int {
	// TODO: add your code here!
	if val1, ok := items[key]; ok {
	i := val1
	return i.value
		
	}
	return -1
}
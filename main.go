package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func append_queue(buffer []string,ql *int,elem string){
	buffer[*ql] = elem
	*ql++
}

func create_gopher_entry(file string, path_from_root string) string{
	buffer := make([]byte,512)
	fp,err := os.Open(file)
	if err != nil{
		panic(err)
	}
	defer fp.Close()
	n,err := fp.Read(buffer)
	if err != nil && err != io.EOF{
		panic(err)
	}
	ftype := http.DetectContentType(buffer[:n])
	if(path_from_root == ""){
		if(ftype == "image/gif"){
			return 	fmt.Sprintf("g%s\t/%s\r\n",fp.Name(),fp.Name())
		}else if(strings.Contains(ftype,"image")){
			return fmt.Sprintf("I%s\t/%s\r\n",fp.Name(),fp.Name())
		}else if(strings.Contains(ftype,"text")){
			return fmt.Sprintf("0%s\t/%s\r\n",fp.Name(),fp.Name())
		}else if(strings.Contains(ftype,"application") || strings.Contains(ftype,"audio")){
			return   fmt.Sprintf("9%s\t/%s\r\n",fp.Name(),fp.Name())
		}
	}else{
		if(ftype == "image/gif"){
			return 	fmt.Sprintf("g%s\t/%s/%s\r\n",fp.Name(),path_from_root,fp.Name())
		}else if(strings.Contains(ftype,"image")){
			return fmt.Sprintf("I%s\t/%s/%s\r\n",fp.Name(),path_from_root,fp.Name())
		}else if(strings.Contains(ftype,"text")){
			return fmt.Sprintf("0%s\t/%s/%s\r\n",fp.Name(),path_from_root,fp.Name())
		}else if(strings.Contains(ftype,"application") || strings.Contains(ftype,"audio")){
			return fmt.Sprintf("9%s\t/%s/%s\r\n",fp.Name(),path_from_root,fp.Name())
		}		
	}
	panic("Unreachable code: create_gopher_entry 39")
	return ""
}

func get_depth(starting_path string)string{
	path,err := os.Getwd()
	if err != nil{
		panic(err)
	}
	split_path := strings.Split(path,starting_path);
	relative_path := split_path[1]
	relative_path = relative_path[1:]
	return relative_path
}

func map_directory_contents(start_path string) {
	path_queue := make([]string,100)
	qp := 0
	ql := 1
	path_queue[0] = start_path
	builder := strings.Builder{}
	builder.Grow(512)
	for(qp < ql){
		fmt.Println(os.Getwd())
		err := os.Chdir(path_queue[qp])
		if err != nil{
			panic(err)
		}
		entries, err := os.ReadDir(".")
		if err != nil{
			panic(err)
		}
		for _, entry := range entries{
			if (entry.Name() ==  "gophermap"){
				continue 
			}
			if entry.IsDir(){
				cwd,err := os.Getwd()
				if err != nil{
					panic(err)
				}
				dir_path := (cwd + "/" + entry.Name())
				append_queue(path_queue,&ql,dir_path)
				if(path_queue[qp] == start_path){
 					formated_gopher := fmt.Sprintf("1%s\t/%s\r\n",entry.Name(),entry.Name())
					builder.WriteString(formated_gopher)
				}else{
					path_from_root := get_depth(start_path)
					formated_gopher := fmt.Sprintf("1%s\t/%s/%s\r\n",entry.Name(),path_from_root,entry.Name())
					builder.WriteString(formated_gopher)
				}
			}else{
				gopher_entry := ""
				if(path_queue[qp] == start_path){
					gopher_entry = create_gopher_entry(entry.Name(),"")
				}else{
					path_from_root := get_depth(start_path)
					gopher_entry = create_gopher_entry(entry.Name(),path_from_root)
				}
				builder.WriteString(gopher_entry)
			}	
		}
		qp++;
		gopher_map := builder.String()
		gopherfile,err := os.Create("gophermap")
		if err != nil{
			panic(err)
		}
		gopherfile.WriteString(gopher_map)
		gopherfile.Close()
		builder.Reset()
	}
}

func main(){
	path := ""
	if(len(os.Args) < 2){
		cwd,err := os.Getwd()
		if err != nil{
			panic(err)
		}
		path = cwd
	}else{
		full_path, err := os.Getwd()
		if err != nil{
			panic(err)
		}
		path = full_path + "/" + os.Args[1]
	}
	fmt.Println(os.Getwd())
	map_directory_contents(path)
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {

	// Check Arguments
	arglen := len(os.Args)
	if arglen < 2 {
		fmt.Printf("Usage: fs input.txt\n")
		return
	}

	// Open Input file with commands list
	cmdFile := os.Args[1]
	f, err := os.Open(cmdFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	// Initialize scanner for reading Input file
	newsc := bufio.NewScanner(f)

	// Initialize root directory
	fs := InitFileSystem()

	// while there is data, keep reading and parsing the lines
	for newsc.Scan() {
		cmdtxt := newsc.Text()
		cmdopts := strings.Split(cmdtxt, "  ")
		switch cmdopts[0] {
		case "dir":
			fmt.Printf("Command: %s\n", cmdtxt)
			fs.dir()
		case "mkdir":
			// for now ignore
			if len(cmdopts) < 2 {
				continue
			}
			fmt.Printf("Command: %s\n", cmdtxt)
			mkdircmdopts := strings.Split(cmdtxt, "   ")
			fs.mkdir(mkdircmdopts[1])
		case "cd":
			fmt.Printf("Command: %s\n", cmdtxt)
			cdcmdopts := strings.Split(cmdtxt, "      ")
			fs.cd(cdcmdopts[1])
		case "up":
			fmt.Printf("Command: %s\n", cmdtxt)
			fs.up()
		case "mv":
			fmt.Printf("Command: %s\n", cmdtxt)
			dir1 := ""
			dir2 := ""
			firstdone := false
			for i := 2; i < len(cmdtxt); i++ {
				if cmdtxt[i] == ' ' {
					if len(dir1) > 0 {
						firstdone = true
					}
					continue
				}
				if !firstdone {
					dir1 += string(cmdtxt[i])
				} else {
					dir2 += string(cmdtxt[i])
				}
			}
			fs.mv(dir1, dir2)
		case "tree":
			fmt.Printf("Command: %s\n", cmdtxt)
			fs.tree()
		default:
			continue
		}

	}
}

// Init flag to check if everything is initialized
var fileSystemInit = false

// Each node in the filesystem
type DirNode struct {
	dirname string
	depth   int
	dirs    map[string]DirNode
	parent  *DirNode
}

// Root File System Interface
type FileSystem struct {
	rootDir  DirNode
	cdirName string // Where are we right now
	cdirNode *DirNode
}

func createNode(cdir *DirNode, depth int, name string) DirNode {
	dnode := DirNode{dirname: name,
		dirs:   map[string]DirNode{},
		depth:  depth,
		parent: cdir,
	}
	return dnode
}

func InitFileSystem() FileSystem {
	// Initialize root node
	rnode := DirNode{dirname: "root",
		dirs:   map[string]DirNode{},
		depth:  0,
		parent: nil,
	}
	Fs := FileSystem{rootDir: rnode,
		cdirName: "root",
		cdirNode: &rnode,
	}
	fileSystemInit = true
	return Fs
}

func (fs *FileSystem) mkdir(dirname string) {
	cdir := fs.cdirNode

	//fmt.Printf("mkdir: current directory is %s\n", fs.cdirName)
	// Check if the directory exists
	_, ok := cdir.dirs[dirname]
	if ok {
		fmt.Println("Subdirectory already exists")
		return
	}
	v := strings.Split(fs.cdirName, "\\")
	newnode := createNode(cdir, len(v), dirname)
	cdir.dirs[dirname] = newnode
	return
}

func (fs *FileSystem) dir() {
	cdir := fs.cdirNode
	dirlist := []string{}

	fmt.Printf("Directory of %s:\n", fs.cdirName)
	for k, _ := range cdir.dirs {
		dirlist = append(dirlist, k)
	}
	sort.Strings(dirlist)
	if len(dirlist) == 0 {
		fmt.Println("No subdirectories")
		return
	}
	dl := len(dirlist)
	nlines := dl / 10
	if dl%10 > 0 {
		nlines += 1
	}
	for i := 0; i < dl; i++ {
		if (i)%10 == 0 && i != 0 {
			fmt.Println()
		}
		if i != dl-1 {
			fmt.Printf("%-8v", dirlist[i])
		} else {
			fmt.Printf("%s", dirlist[i])
		}
	}
	fmt.Println()
	return
}

func (fs *FileSystem) up() {
	cdirNode := fs.cdirNode
	//fmt.Printf("up: current directory is %s\n", fs.cdirName)
	if fs.cdirName == "root" {
		fmt.Println("Cannot move up from root directory")
		return
	} else {
		// parse the pathname
		tstr := strings.Split(fs.cdirName, "\\")
		newstr := ""
		for i := 0; i < len(tstr)-1; i++ {
			newstr += tstr[i]
			if i != len(tstr)-2 {
				newstr += "\\"
			}
		}
		fs.cdirName = newstr
		parent := cdirNode.parent
		fs.cdirNode = parent
	}
}

func (fs *FileSystem) cd(dirname string) {
	//fmt.Printf("cd: current directory is %s\n", fs.cdirName)
	cdirNode := fs.cdirNode
	tdir, ok := cdirNode.dirs[dirname]
	if !ok {
		fmt.Println("Subdirectory does not exist")
		return
	}
	fs.cdirNode = &tdir
	fs.cdirName += "\\"
	fs.cdirName += dirname
}

var startmarker string = "├──"
var endmarker string = "└──"
var nomarker string = "│   "

func treeHelper(node *DirNode, buf string) {
	dirnames := []string{}
	for k := range node.dirs {
		v, _ := node.dirs[k]
		dirnames = append(dirnames, v.dirname)
	}
	sort.Strings(dirnames)
	for index, dirent := range dirnames {
		dirv, _ := node.dirs[dirent]
		if index == len(dirnames)-1 {
			fmt.Println(buf+endmarker, dirent)
			treeHelper(&dirv, buf+"    ")
		} else {
			fmt.Println(buf+startmarker, dirent)
			treeHelper(&dirv, buf+nomarker)
		}
	}
}

func (fs *FileSystem) tree() {
	fmt.Printf("Tree of %s:\n", fs.cdirName)
	fmt.Println(".")
	treeHelper(fs.cdirNode, "")
}

// returns dirNode, exists or not, local or not
func getdestNode(fs *FileSystem, dest string) (*DirNode, bool, bool) {
	cdir := fs.cdirNode
	dstarr := strings.Split(dest, "\\")
	if len(dstarr) == 1 || (len(dstarr) == 2 && dstarr[0] == ".") {
		ddirNode, dok := cdir.dirs[dstarr[1]]
		if dok {
			return &ddirNode, true, true
		} else {
			return nil, false, true
		}
	} else {
		cnode := fs.cdirNode
		for i := 0; i < len(dstarr); i++ {
			if dstarr[i] == "." {
				continue
			}
			if dstarr[i] == ".." {
				cnode = cnode.parent
			} else {
				d, ok := cnode.dirs[dstarr[i]]
				if ok {
					cnode = &d
					return cnode, true, false
				} else {
					return cnode, false, false
				}
			}
		}
	}
	return nil, false, false
}

// Recursively go through the node and update
// depths of all nodes underneath this node.
func updateDepth(dnode *DirNode) {
	if len(dnode.dirs) > 0 {
		for k := range dnode.dirs {
			v, ok := dnode.dirs[k]
			if ok {
				v.depth = dnode.depth + 1
				updateDepth(&v)
			}
		}
	}
}

func (fs *FileSystem) mv(tsrc, dest string) {
	cdir := fs.cdirNode
	src := ""
	parsesrc := strings.Split(tsrc, "\\")
	if len(parsesrc) == 1 {
		src = parsesrc[0]
	} else {
		src = parsesrc[len(parsesrc)-1]
	}
	sdirNode, sok := cdir.dirs[src]
	if !sok {
		fmt.Printf("%s directory doesnt exist.\n", src)
		return
	}
	// check if this is local or not
	// if its not-local and dest doesnt exist, fetch the remote dir parent
	// if dest exists, get the dirNode
	ddirNode, exists, local := getdestNode(fs, dest)
	if !exists && local {
		// Directory not present
		// Rename the src name to destination name
		sdirNode.dirname = dest
		delete(cdir.dirs, src)
		cdir.dirs[dest] = sdirNode
	}
	// These conditions can be collapsed as they do similar work.
	// But for easy reading, leaving them separate.
	if exists && local {
		_, ok := ddirNode.dirs[src]
		if ok {
			fmt.Println("Subdirectory already exists")
			return
		}
		// destination dir present
		// Move src to destination updating parent
		// update depth in all the src path objects
		sdirNode.parent = ddirNode
		sdirNode.depth = ddirNode.depth + 1
		// add new node to the destination directory
		ddirNode.dirs[src] = sdirNode
		//Remove sdir from dirs map of curdir
		delete(cdir.dirs, src)
		// update depth of all nodes under this
		updateDepth(&sdirNode)
	}
	if exists && !local {
		_, ok := ddirNode.dirs[src]
		if ok {
			fmt.Println("Subdirectory already exists")
			return
		}
		sdirNode.parent = ddirNode
		sdirNode.depth = ddirNode.depth + 1
		delete(cdir.dirs, src)
		ddirNode.dirs[src] = sdirNode
		updateDepth(&sdirNode)
	}

	if !exists && !local {
		sdirNode.dirname = dest
		delete(cdir.dirs, src)
		// add this node to dest directory
		sdirNode.parent = ddirNode
		sdirNode.depth = ddirNode.depth + 1
		ddirNode.dirs[src] = sdirNode
		updateDepth(&sdirNode)
	}
}

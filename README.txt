   Virtual FileSystem Implementation

1. Goals:
    Build an in-memory File-System with support for these commands:
    1. dir  (Display Objects in cur directory)
    2. mkdir (Create a directory in cur directory)
    3. cd (Change directory)
    4. up (Go One level up in the current directory)
    5. tree (Display tree of the current virtual file system)
    6. mv (Move a directory)


2. Input
    Input/Tests are mentioned in the attached document input1.txt 
    and input2.txt with respective output in oend1.txt and oend2.txt.


3. Design
    Will separate this into three parts.
    1. Input Processing
    2. Actual Implementation.
    3. Display   


    Input Processing:
    Input is read from a file line by line using Golang bufio Scanner and
    parsed based on the requirement document. There is different spacing for
    each command so its done in different ways.
    Invalid input is skipped as there is no requirement to handle that.


    Implementation:
    There are two structures I Used to get this work done. Organized
    directories into a tree like structure.

    A Directory Node representation
    type DirNode struct {
        dirname string
        depth   int
        dirs    map[string]DirNode
        parent  *DirNode // Parent of the current node.
    }

    A File System Struct to implement the main Interface to handle API's.
    type FileSystem struct {
        rootDir  DirNode
        cdirName string     // Path Where are we right now
        cdirNode *DirNode   // DirNode object for the same
    }

    API's:
    1. mkdir(name string)
        Creates an object in the curDirNode from FS struct.
    2. dir()
        Lists all files in the curDir from the FS struct.
    3. cd(name string)
        Changes the directory to the specified name if it exists.
        Gets cur dir Obj handle and reads the children to find it the name
        exists. 
    3. up()
        Changes the directory to one level up.
        Every Object has a pointer to the parent object.
    3. tree()
        Prints of the tree of the current directory.
        Use DFS algorithm and using "depth" parameter in the DirNode object to 
        print the levels into a strings Builder interface in Golang.
    4.  mv(src, dest string)
        Moved src to dest depending on if the destination folder exists.
        Handle depth parameter for the new tree path recursively.


    Implementation:
    ----------------
    Golang
    
    Files:
    -----
    README.md
    filesystem.go (Source Code)
    input1.txt, input2.txt (Input Data for the Binary)
    output1.txt, output2.txt (Output generated from Binary)
    expout1.txt, expout2.txt (Expected Output for the tests)
    
    Usage:
    ------
    compile
    go build filesystem.go

    ./filesystem input1.txt > output1.txt
    ./filesystem input2.txt > output2.txt
    
        


